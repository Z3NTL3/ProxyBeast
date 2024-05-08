package core

type Controller struct {
	started     bool
	current     uint64  // current
	total       uint64  // total work
	worker_pool Workers // worker pool
	fd_pool     FD_Pool // fd pool
	done        chan int
}

type Workers chan struct {
	proxy Proxy
}

type FD_Pool chan struct {
	proxy     Proxy
	latency   uint32
	anonimity Anonimity
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
