//go:build windows

package lan

import (
	// "bufio"
	"context"
	"encoding/json"
	// "errors"
	"fmt"
	"log"
	"net"
	"net/http"
	// "os"
	"sync"
	"syscall"
	"time"

	// "github.com/gorilla/websocket"
)





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




