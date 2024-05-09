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

import "context"

var (
	JUDGE string = "https://pool.proxyspace.pro/"
)

type Controller struct {
	started     bool
	current     uint64  // current

	worker_pool Workers // worker pool
	fd_pool     FD_Pool // fd pool

	done        context.Context
	cancel 		*context.CancelFunc
}

type Workers chan struct {
	proxy Proxy
}

type FD_Pool chan struct {
	proxy     Proxy
	latency   uint32
	anonimity Anonimity
}

func(c *Controller) Register(ctx context.Context, cancel context.CancelFunc){
	c.done = ctx
	c.cancel = &cancel
}

func (c *Controller) Add(n uint64) {
	c.current += n
}

func (c *Controller) Done(n uint64) {
	c.current -= n
}

func (c *Controller) Current() uint64 {
	return c.current
}

func (c *Controller) DidStart() bool {
	return c.started
}

func (c *Controller) Cancel()  {
	(*c.cancel)()
}

func(c *Controller) ShouldCancel() <-chan struct{}{
	return c.done.Done()
}