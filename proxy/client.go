package proxy

import "os"

type (
	Anonimity string
	Proxy     string

	/*
		This client is lightweight and very unique as it is working with
		several pools, in an event driven manner.


		Take the following as example:

		Before starting any worker, the fd_pool workers gets spawned.
		They can process 200 file writes at the same time, while only
		one file descriptor stays open, until the end.

		The proxies in the provided proxy list get's read line by line
		in the main thread. Everytime a line is read it will start a goroutine, which
		just checks one given proxy. Plus incrementing ``current``.

		This means all your proxies are checked in isolated environments and all in parallel
		with each other. Once one finalizes it i'll decrement one of ``current``

		After those workers get spawned, the pool watcher is waiting until
		all goroutines have finished, if in the meantime the user interrupts the operation
		for example with CTRL + C it will kill and stop all workers and finally exit.


	*/
	ProxyClient struct {
		*os.File // one open file socket to the 'saves' file until the end of the operation
		current uint64 // current amount of running goroutines; 0 means all tasks finished

		/*
			Transmits data to the channel
			in order for the good proxies to be written into
			the 'saves' file.

			FD pool has concurrent goroutines that listen for
			incoming transmits, in parallel
		*/
		fd_pool chan struct {
			proxy     Proxy
			latency   uint32
			anonimity Anonimity
		}

		/*
			Transmits proxies that need to be checked
			to the worker channel. If given proxy is working
			it will report its qualities and transmit further to the ``fd_pool``.
		*/
		workers chan struct {
			proxy Proxy
		}
	}
)
