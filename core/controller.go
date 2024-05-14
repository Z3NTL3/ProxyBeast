/*
     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App core]
*/

package core

import (
	"context"
	"sync/atomic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// worker pool
// its size can be customized on the settings
// page
type Workers struct {
	proxy Proxy
}

// file descriptor pool
// which is always a size of 20 concurrent
// goroutines
type FD_Pool struct {
	Proxy     Proxy  `json:"proxy"`
	Latency   string `json:"latency"`
	Anonimity string `json:"anonimity"`
}


// controller for the cooperative pools
type Controller struct {
	current uint32
	load uint32

	cthread int32
	

	capacity int32
	worker_pool chan Workers // worker pool
	fd_pool     chan FD_Pool // fd pool

	done   context.Context     // signal context to know when to exit/cancel
	cancel *context.CancelFunc // cancel helper for the context, intended for exit, gets invoked explicitly or when task completes
}

func (c *Controller) ThreadCompletion() {
	atomic.AddInt32(&c.cthread, 1)
}

func (c *Controller) CurrentThread() int32 {
	return atomic.LoadInt32(&c.cthread)
}

func (c *Controller) Done() {
	atomic.AddUint32(&c.current, 1)
	go runtime.EventsEmit(APP.ctx, Fire_CurrentThread, c.Current())
}

func (c *Controller) Current() uint32 {
	return atomic.LoadUint32(&c.current)
}

func (c *Controller) SetLoad(n uint32) {
	atomic.StoreUint32(&c.load, n)
}

func (c *Controller) GetLoad() uint32 {
	return atomic.LoadUint32(&c.load)
}

func (c *Controller) ShouldStop() <- chan struct{} {
	return c.done.Done()
}

func (c *Controller) CanExit() {
	done := make(chan int)
	go func(){
		for {
			if c.CurrentThread() == atomic.LoadInt32(&c.capacity) {
				done <- 1
				return
			}
		}
	}()

	select {
		case <-done:
			return
	}
}

func (c *Controller) Register(ctx context.Context, cancel context.CancelFunc) {
	c.done = ctx
	c.cancel = &cancel
}

// reset controller state
func (c *Controller) Reset() {
	close(c.fd_pool)
	close(c.worker_pool)

	c.fd_pool = make(chan FD_Pool, 20)
	c.worker_pool = make(chan Workers, DefaultPoolSize)
	c.current = 0
	c.cthread = 0

	// clear fd
	// FD[InputFile] = nil
	// FD[SaveFile] = nil

	// register new context and cancel func
	c.Register(context.WithCancel(context.Background()))
}
