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
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Controller struct {
	started     bool
	current     uint64  // smart counter, 0 means work has finished

	cntr_current 	uint32  // progress counter, points to current 
	cntr_load		uint32 // total amount of load

	worker_pool chan Workers // worker pool
	fd_pool     chan FD_Pool // fd pool

	done        context.Context
	cancel 		*context.CancelFunc

	mx sync.RWMutex
}

type Workers  struct {
	proxy Proxy
}

type FD_Pool struct {
	Proxy     Proxy `json:"proxy"`
	Latency   string `json:"latency"`
	Anonimity string `json:"anonimity"`
}

func(c *Controller) SetLoad(n uint32) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.cntr_load = n
}

func(c *Controller) Start() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.started = true
}

func (c *Controller) Cancel()  {
	(*c.cancel)()
}

func(c *Controller) ShouldStop() <-chan struct{}{
	return c.done.Done()
}

func(c *Controller) Register(ctx context.Context, cancel context.CancelFunc){
	c.done = ctx
	c.cancel = &cancel
}

func (c *Controller) Add(n uint64) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.current += n
}

func (c *Controller) Done(n uint64) {
	defer c.mx.Unlock()
	c.mx.Lock()

	if c.current - n < 0 || c.cntr_current +1 > c.cntr_load { return }

	c.current -= n
	c.cntr_current += 1

	runtime.EventsEmit(APP.ctx, Fire_CurrentThread, c.cntr_current)
}

func (c *Controller) Current() uint64 {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.current
}

func (c *Controller) DidStart() bool {
	c.mx.RLock()
	defer c.mx.RUnlock()
	
	return c.started
}

func(c *Controller) Reset() {
	c.fd_pool = make(chan FD_Pool, 20)
	c.worker_pool = make(chan Workers, DefaultPoolSize)
	c.started = false
	c.current = 0
	c.cntr_load = 0
	c.cntr_current = 0

	c.Register(context.WithCancel(context.Background()))
}

func (c *Controller) CanExit() {
	exit := make(chan int)

	go func(){
		for {
			if c.Current() == 0 {
				exit <- 1
			}
		}
	}()

	select {
		case <-exit:
			return
	}
}