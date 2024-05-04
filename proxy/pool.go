package proxy

/*
	About the concept, particularly how it operates with
	several pools, in an event driven manner.


	Take the following as example:

	Before starting any worker, the fd_pool workers gets spawned.
	They can process N amount of file writes at the same time, while only
	one file descriptor stays open, until the end (just one no dup). Waiting
	until they can consume from the channel which is buffered.

	The proxies in the provided proxy list get's read line by line
	in the main thread. Everytime a line is read it will push new worker into the worker pool, which
	just checks one given proxy. Plus incrementing ``current``. The
	worker pool is spawned after the FD_pool, they consume from the buffered
	workers channel.

	This means all your proxies are checked in isolated environments and all in parallel in a co-operative way
	with each other. Once one finalizes it i'll decrement one of ``current``

	After those workers get spawned, the pool watcher is waiting until
	all goroutines have finished, if in the meantime the user interrupts the operation
	for example with CTRL + C it will kill and stop all workers and finally exit.
*/

type Controller struct {
	Current uint64 // current
	Total uint64 // total work
}

type Workers chan struct {
	proxy Proxy
}

type FD_Pool chan struct {
	proxy     Proxy
	latency   uint32
	anonimity Anonimity
}
