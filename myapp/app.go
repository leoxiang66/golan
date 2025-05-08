package main

import (
	"context"
	"fmt"
	"myapp/lan"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	go lan.Launch()
	go func() {
		var mu sync.Mutex
		for {

			// publish
			publish_peers(ctx)

			// clear
			mu.Lock()
			lan.PEERS = make(map[string]lan.Peer)
			mu.Unlock()

			// interval
			time.Sleep(10 * time.Second)
		}
	}()
}

func publish_peers(ctx context.Context) {
	peers := []string{}
	for k := range lan.PEERS {
		peers = append(peers, k)
	}
	runtime.EventsEmit(ctx, "lan:peers", peers)
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) {
	fmt.Println(name)
}

func (a *App) InviteSocket(id string, timeout_s int) (bool, error) {
	// 1) 定义一个结果类型，既包含连接也包含错误
	type inviteResult struct {
		Conn *websocket.Conn
		Err  error
	}

	result := make(chan inviteResult, 1)

	// 2) 后台 goroutine 去执行实际的 InviteSocket
	go func() {
		conn, err := lan.InviteSocket(id)
		result <- inviteResult{Conn: conn, Err: err}
	}()

	// 3) 等待结果或超时
	select {
	case res := <-result:
		if res.Conn == nil {
			fmt.Println("connection失败,返回false")
			return false, res.Err
		} else {
			go func() {
				// 成功拿到 conn，进入聊天
				// 聊天结束后才关闭连接并返回
				lan.ChatLoop(res.Conn)
				res.Conn.Close()
				runtime.EventsEmit(a.ctx,"lan:conn_closed")
				fmt.Println("Connection closed")
			}()
			return true, nil
		}

	case <-time.After(time.Duration(timeout_s) * time.Second):
		// 超时返回自定义错误
		return false, fmt.Errorf("invite to %s timed out after %d seconds", id, timeout_s)
	}
}
