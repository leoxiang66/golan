//go:build !windows

package lan

import (
	// "bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	// "os"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const (
	udpPort      = 9999
	broadcastIP  = "255.255.255.255"
	broadcastInt = time.Second
)

// DiscoveryMsg 用于 UDP 广播/接收
type DiscoveryMsg struct {
	ID     string `json:"id"`
	WSPort int    `json:"wsPort"`
}

type Peer struct {
	IP     string
	WSPort int
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var PEERS map[string]Peer
var id string

func Launch() {
	// 1) 生成基于当前时间的唯一 ID
	id = time.Now().Format("2006-01-02T15:04:05")
	log.Printf("本节点 ID = %s\n", id)

	// 2) 随机端口启动 WebSocket 服务
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wsHandler)

	wsListener, err := net.Listen("tcp", ":0") // 随机可用端口
	if err != nil {
		log.Fatalf("无法监听 WS 端口: %v", err)
	}
	wsPort := wsListener.Addr().(*net.TCPAddr).Port
	log.Printf("WebSocket 服务监听 TCP :%d/ws\n", wsPort)
	go http.Serve(wsListener, mux)

	// 3) UDP ListenConfig：打开 SO_BROADCAST/SO_REUSEADDR/SO_REUSEPORT
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var serr error
			if err := c.Control(func(fd uintptr) {
				// 允许 UDP 广播
				serr = syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
				if serr != nil {
					return
				}
				// 允许地址/端口重用
				serr = syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if serr != nil {
					return
				}
				serr = syscall.SetsockoptInt(int(fd),
					syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
			}); err != nil {
				return err
			}
			return serr
		},
	}

	udpConn, err := lc.ListenPacket(context.Background(), "udp4", fmt.Sprintf(":%d", udpPort))
	if err != nil {
		log.Fatalf("无法监听 UDP %d: %v", udpPort, err)
	}
	defer udpConn.Close()

	// 维护已见 peers
	PEERS = make(map[string]Peer)
	var mu sync.Mutex

	// 4) UDP 接收，发现新节点时提示邀请
	go func() {
		buf := make([]byte, 512)
		for {
			n, addr, err := udpConn.ReadFrom(buf)
			if err != nil {
				log.Printf("UDP 读取错误: %v", err)
				continue
			}
			var msg DiscoveryMsg
			if err := json.Unmarshal(buf[:n], &msg); err != nil {
				continue // 忽略非本协议报文
			}
			if msg.ID == id {
				continue // 排除自己
			}

			mu.Lock()
			_, seen := PEERS[msg.ID]
			if !seen {
				ip := addr.(*net.UDPAddr).IP.String()
				PEERS[msg.ID] = Peer{
					IP:     ip,
					WSPort: msg.WSPort,
				}
			}
			mu.Unlock()

			// if !seen {

			//     fmt.Printf("\n发现新节点：ID=%q, 地址=%s, WS 端口=%d\n", msg.ID, ip, msg.WSPort)
			//     promptInvite(ip, msg.WSPort, id)
			// }
		}
	}()

	// 5) UDP 持续广播自己的 ID + WS 端口
	go func() {
		dst := &net.UDPAddr{IP: net.ParseIP(broadcastIP), Port: udpPort}
		msg := DiscoveryMsg{ID: id, WSPort: wsPort}
		data, _ := json.Marshal(msg)
		for {
			if _, err := udpConn.WriteTo(data, dst); err != nil {
				log.Printf("UDP 广播失败: %v", err)
			} else {
				// log.Printf("已广播: %s", string(data))
			}
			time.Sleep(broadcastInt)
		}
	}()

	// 阻塞主线程
	select {}
}

// promptInvite 询问用户是否邀请对端
// func promptInvite(peerIP string, peerWSPort int, id string) {
//     fmt.Printf("是否邀请 %s:%d 聊天? [y/N]: ", peerIP, peerWSPort)
//     var ans string
//     fmt.Scanln(&ans)
//     if ans != "y" && ans != "Y" {
//         return
//     }
//     inviteSocket(peerIP, peerWSPort, id)
// }

// inviteSocket 发起 WebSocket 邀请并聊天
func InviteSocket(guestID string) (*websocket.Conn,error){
	if peer, ok := PEERS[guestID]; ok {
		peerIP, peerWSPort := peer.IP, peer.WSPort
		url := fmt.Sprintf("ws://%s:%d/ws", peerIP, peerWSPort)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Dial %s 失败: %v", url, err)
			return nil,err
		}
		// 发送 invite
		conn.WriteJSON(map[string]string{"type": "invite", "from": id})

		// 读取响应
		var resp struct{ Type string }
		if err := conn.ReadJSON(&resp); err != nil {
			log.Printf("Failed to read response: %v", err)
			return nil,err
		}
		if resp.Type != "accept" {
			log.Printf("The invitation is rejected")
			return nil,err
		}
		// chatLoop(conn)
        return conn,nil
	}
    return nil,errors.New("不存在此Peer")

}

// wsHandler 处理传入的 invite
func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS Upgrade 失败:", err)
		return
	}
	defer conn.Close()

	// 读取 invite 消息
	var msg struct {
		Type string `json:"type"`
		From string `json:"from"`
	}
	if err := conn.ReadJSON(&msg); err != nil {
		log.Println("读取 invite 失败:", err)
		return
	}
	if msg.Type != "invite" {
		return
	}

	// 提示是否接受
	fmt.Printf("\n收到 %s 的聊天邀请，接受? [y/N]: ", msg.From)
	var ans string
	fmt.Scanln(&ans)
	if ans != "y" && ans != "Y" {
		conn.WriteJSON(map[string]string{"type": "reject"})
		return
	}

	// 同意并进入聊天
	conn.WriteJSON(map[string]string{"type": "accept"})
	ChatLoop(conn)
}

// chatLoop 终端与 WebSocket 双向聊天
func ChatLoop(conn *websocket.Conn) {
	// 从 WS -> 终端
	go func() {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("WS 读取错误:", err)
                break
				// os.Exit(0)
			}
			fmt.Printf("\n<< %s\n>>> ", string(data)) //todo
		}
	}()

	// todo: send
	for {
		line := "hi"

		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			log.Println("WS 发送错误:", err)
			return
		}
        time.Sleep(time.Second)
	}
}
