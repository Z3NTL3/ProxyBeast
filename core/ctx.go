/*
	ProxyBeast GUI

The ultimate proxy checker

	   by @z3ntl3

	[proxy.pix4.dev]

License: GNU
Note: If you've liked ProxyBeast, please consider starring our Github repository.

[App core]
*/
package core

import (
	"fmt"
	"time"
)

type CheckerCtx struct {
	Scheme string
	fd_pool chan FD_Pool 
}

func (c *CheckerCtx) WRAP_COMPLETION(proto string, proxy Proxy, done chan<- int) error {
	defer func ()  {
		done <- 1	
	}()

	return c.WRAP(proto, proxy)
}

func (c *CheckerCtx) WRAP(proto string, proxy Proxy) (err error) {
	var level string

	start := time.Now().UnixMilli()
	switch {
		case proto == HTTP:
			// http
			level, err = c.HTTP(proxy)
			if err != nil {return}
		case proto == HTTPS:
			// https
			level, err = c.HTTPS(proxy)
			if err != nil {return}
		case proto == SOCKS4:
			// socks4
			level, err = c.SOCKS4(proxy)
			if err != nil {return}
		case proto == SOCKS5:
			// socks5
			level, err = c.SOCKS5(proxy)
			if err != nil {return}
		default:
			err = ErrUnknownProtocol
			return
	}

	anonimity := Anonimity(level)
	level = anonimity.GetAnonimity()
	latency := fmt.Sprintf("%vms", time.Now().UnixMilli() - start)

	go func(){
		// GODOC: If recover is called outside the deferred function it will not stop a panicking sequence.
		// so defer recover() - doesnt work has to be wrapped around a function
		defer func(){
			recover()

			// recover here is required, early cancellation signal
			// will close the channel, sends to closed channel panic
			// to perceive that state and continue from it, it is required
			// to handle in special way
		}()

		// on shut down, sents to channels will block.
		// on duty, sents to channel will block if the
		// channel's buffer is full, running in goroutine
		// ensures no blocks.
		c.fd_pool <- FD_Pool{
			Proxy: proxy,
			Latency: latency,
			Anonimity: level,
			Protocol: proto,
		}
	}()

	return
}