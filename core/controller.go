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
)

type Controller struct {
	started     bool
	current     uint64  // current

	worker_pool chan Workers // worker pool
	fd_pool     chan FD_Pool // fd pool

	done        context.Context
	cancel 		*context.CancelFunc

	mx sync.Mutex
}

type Workers  struct {
	proxy Proxy
}

type FD_Pool struct {
	proxy     Proxy
	latency   uint32
	anonimity Anonimity
}

func(c *Controller) Start() {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.started = true
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
	c.mx.Lock()
	defer c.mx.Unlock()

	if c.current -n < 0 {
		return
	}
	c.current -= n
}

func (c *Controller) Current() uint64 {
	c.mx.Lock()
	defer c.mx.Unlock()
	
	return c.current
}

func (c *Controller) DidStart() bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	return c.started
}

func (c *Controller) Cancel()  {
	(*c.cancel)()
}

func(c *Controller) ShouldStop() <-chan struct{}{
	return c.done.Done()
}

func(c *Controller) Reset() {
	c.fd_pool = make(chan FD_Pool, 20)
	c.worker_pool = make(chan Workers, 2000)
	c.started = false
	c.current = 0

	c.Register(context.WithCancel(context.Background()))
}