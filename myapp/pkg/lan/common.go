package lan

import (
	// "bufio"
	// "context"
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	// "net"
	"net/http"
	// "os"
	// "sync"
	// "syscall"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"myapp/pkg/konst"
	"time"
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

// inviteSocket 发起 WebSocket 邀请并聊天
func InviteSocket(guestID string) (*websocket.Conn, error) {
	if peer, ok := PEERS[guestID]; ok {
		peerIP, peerWSPort := peer.IP, peer.WSPort
		url := fmt.Sprintf("ws://%s:%d/ws", peerIP, peerWSPort)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Printf("Dial %s 失败: %v", url, err)
			return nil, err
		}
		// 发送 invite
		conn.WriteJSON(map[string]string{"type": "invite", "from": id})

		// 读取响应
		var resp struct{ Type string }
		if err := conn.ReadJSON(&resp); err != nil {
			log.Printf("Failed to read response: %v", err)
			return nil, err
		}
		if resp.Type != "accept" {
			log.Printf("The invitation is rejected")
			return nil, err
		}
		// chatLoop(conn)
		return conn, nil
	}
	return nil, errors.New("不存在此Peer")

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
	runtime.EventsEmit(konst.Ctx, "lan:receive_invite", msg.From)

	result := <-konst.BoolChan

	if result {
		conn.WriteJSON(map[string]string{"type": "accept"})
		ChatLoop(conn)
	} else {
		conn.WriteJSON(map[string]string{"type": "reject"})
		return
	}


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
		
		// line := "hi"

		// if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
		// 	log.Println("WS 发送错误:", err)
		// 	return
		// }
		time.Sleep(time.Second)
	}
}
