package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"myapp/lan"
	"sync"
	"time"
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
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) InviteSocket(id string) bool {

	go lan.InviteSocket(id)
	select{}
    // result channel, buffered so the goroutine can't block forever
    // result := make(chan bool, 1)

    // fire off the actual invite in its own goroutine
    // go func() {
    //     result <- lan.InviteSocket(id)
    // }()

    // // wait for either the InviteSocket result or the timeout
    // select {
    // case ok := <-result:
    //     return ok                // got a real answer before timeout
    // case <-time.After(time.Duration(timeout_s) * time.Second):
    //     return false             // timed out
    // }
}
