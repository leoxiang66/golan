//go:build windows
package lan

import (
    "bufio"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
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

type DiscoveryMsg struct {
    ID     string `json:"id"`
    WSPort int    `json:"wsPort"`
}

type Peer struct {
    IP     string
    WSPort int
}



var PEERS map[string]Peer
var id string
var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}



func Launch() {
    id = time.Now().Format("2006-01-02T15:04:05")
    log.Printf("本节点 ID = %s\n", id)

    // 随机端口 WebSocket
    mux := http.NewServeMux()
    mux.HandleFunc("/ws", wsHandler)
    wsListener, err := net.Listen("tcp", ":0")
    if err != nil {
        log.Fatalf("无法监听 WS 端口: %v", err)
    }
    wsPort := wsListener.Addr().(*net.TCPAddr).Port
    log.Printf("WebSocket 服务监听 TCP :%d/ws\n", wsPort)
    go http.Serve(wsListener, mux)

    // UDP：打开广播 + 重用地址
    lc := net.ListenConfig{
        Control: func(network, address string, c syscall.RawConn) error {
            var serr error
            if err := c.Control(func(fd uintptr) {
                h := syscall.Handle(fd)
                if network == "udp4" {
                    serr = syscall.SetsockoptInt(h, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
                    if serr != nil {
                        return
                    }
                }
                serr = syscall.SetsockoptInt(h, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
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

    PEERS = make(map[string]Peer)
    var mu sync.Mutex

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
                continue
            }
            if msg.ID == id {
                continue
            }
            mu.Lock()
            _, seen := PEERS[msg.ID]
            if !seen {
                ip := addr.(*net.UDPAddr).IP.String()
                PEERS[msg.ID] = Peer{
                    IP: ip,
                    WSPort: msg.WSPort,
                }
            }
            mu.Unlock()
            // if !seen {
            //     peerIP := addr.(*net.UDPAddr).IP.String()
            //     fmt.Printf("\n发现新节点：ID=%q, 地址=%s, WS端口=%d\n", msg.ID, peerIP, msg.WSPort)
            //     promptInvite(peerIP, msg.WSPort, id)
            // }
        }
    }()

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

    select {}
}



func InviteSocket(guestID string) bool {
	if peer, ok := PEERS[guestID]; ok {
        peerIP, peerWSPort := peer.IP,peer.WSPort 
		url := fmt.Sprintf("ws://%s:%d/ws", peerIP, peerWSPort)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Dial %s 失败: %v", url, err)
			return false
		}
		defer conn.Close()

		// 发送 invite
		conn.WriteJSON(map[string]string{"type": "invite", "from": id})

		// 读取响应
		var resp struct{ Type string }
		if err := conn.ReadJSON(&resp); err != nil {
			log.Printf("读取响应失败: %v", err)
			return false
		}
		if resp.Type != "accept" {
			log.Printf("对方拒绝邀请")
			return false
		}
		chatLoop(conn)
		return true
	} else {
		return false
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("WS Upgrade 失败:", err)
        return
    }
    defer conn.Close()
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
    fmt.Printf("\n收到 %s 的聊天邀请，接受? [y/N]: ", msg.From)
    var ans string
    fmt.Scanln(&ans)
    if ans != "y" && ans != "Y" {
        conn.WriteJSON(map[string]string{"type": "reject"})
        return
    }
    conn.WriteJSON(map[string]string{"type": "accept"})
    chatLoop(conn)
}

func chatLoop(conn *websocket.Conn) {
    go func() {
        for {
            _, data, err := conn.ReadMessage()
            if err != nil {
                log.Println("WS 读取错误:", err)
                os.Exit(0)
            }
            fmt.Printf("\n<< %s\n>>> ", string(data))
        }
    }()
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(">>> ")
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Println("stdin 读取错误:", err)
            return
        }
        if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
            log.Println("WS 发送错误:", err)
            return
        }
        fmt.Print(">>> ")
    }
}
