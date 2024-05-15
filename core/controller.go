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
	"fmt"
	"sync/atomic"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Worker pool
// Distributed workers that consume their duties from a channel.
// They check proxies concurrently and when a proxy is valid, it is further sent
// to the FD_Pool, which is a pool that basically saves working proxies in the associated
// FD[SaveFile].
type Workers struct {
	proxy Proxy
}

// File descriptor pool
// Please note that this is just a reference to output file.
// 
// These workers listen on a channel, which only keep track of working proxies.
// They get saved in the FD[SaveFile]	
type FD_Pool struct {
	Proxy     Proxy  `json:"proxy"`
	Latency   string `json:"latency"`
	Anonimity string `json:"anonimity"`
}


// Controller for threads to work co-operatively
// with no blocks
type Controller struct {
	current uint32 // progress
	load uint32 // total work	
	threads int32 // amount of threads active
	abort int32 
	
	worker_pool chan Workers // worker pool aka the pool that consumes proxy sent to it from a distributed channel and starts checking
	fd_pool     chan FD_Pool // fd pool aka file descriptor pool for associated files, it consumes good proxies from its channel and saves them in the FD[SaveFile]

	done   context.Context     // signal context to know when to exit/cancel
	cancel *context.CancelFunc // cancel helper for the context, intended for exit, gets invoked explicitly or when task completes
}

// Marks thread completion
func (c *Controller) ThreadCompletion() {
	atomic.AddInt32(&c.threads, -1)
}

// The current thread, 0 means all threads were shutdown
// and it is safe to exit
func (c *Controller) CurrentThread() int32 {
	return atomic.LoadInt32(&c.threads)
}

// Signals completion when one operation in a thread completes
func (c *Controller) Done() {
	atomic.AddUint32(&c.current, 1)
	go runtime.EventsEmit(APP.ctx, Fire_CurrentThread, c.Current())
}

// Current position {current/total_load}
func (c *Controller) Current() uint32 {
	return atomic.LoadUint32(&c.current)
}

// Sets load, the load is the amount of proxies to be checked
func (c *Controller) SetLoad(n uint32) {
	atomic.StoreUint32(&c.load, n)
}

// Gets load (total amount of proxies being load)
func (c *Controller) GetLoad() uint32 {
	return atomic.LoadUint32(&c.load)
}

// Signal for the goroutines to know when to stop and kill itself
func (c *Controller) ShouldStop() <- chan struct{} {
	return c.done.Done()
}


// Signals abort due to fatal error
func (c *Controller) Abort(err error) {
	(*c.cancel)()

	if atomic.LoadInt32(&c.abort) != 1 {
		go runtime.EventsEmit(
			APP.ctx, 
			Fire_ErrEvent, 
			fmt.Sprintf("[ERROR] Aborting due: %s", err.Error()),
		)

		atomic.StoreInt32(&c.abort, 1)
	}
}

// Signals exit, when all goroutines are shutdown
func (c *Controller) CanExit() {
	done := make(chan int)
	go func(){
		for {
			if c.CurrentThread() == 0 {
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

// Initiazes new context and associated cancel
func (c *Controller) Register(ctx context.Context, cancel context.CancelFunc) {
	c.done = ctx
	c.cancel = &cancel
}

// Resets controller state
func (c *Controller) Reset() {
	c.fd_pool = make(chan FD_Pool, 20)
	c.worker_pool = make(chan Workers, DefaultPoolSize)
	c.current = 0
	c.threads = 0
	c.abort = 0

	// clear fd
	//FD[InputFile] = nil
	//FD[SaveFile] = nil

	// register new context and cancel func
	c.Register(context.WithCancel(context.Background()))
}
