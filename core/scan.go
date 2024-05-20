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
	"bufio"
	"context"
	"encoding/json"
	"io"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (c *Controller) StartScan(ctx context.Context, proto string) {
	var err error

	// if err fire to js runtime
	defer func(err_ *error) {
		if *err_ != nil {
			runtime.EventsEmit(APP.ctx, Fire_ErrEvent, (*err_).Error())
		}
	}(&err)

	// determine if a scan is ongoing
	if c.CurrentThread() != 0 {
		err = ErrOngoingCheck
		return
	}

	c.Reset()

	// validite input file
	err = FD.Validate()
	if err != nil {
		return
	}

	runtime.EventsEmit(APP.ctx, Fire_CheckerTotalLoad, c.load)
	c.threads = int32(cap(c.fd_pool) + cap(c.worker_pool))

	defer c.Reset() // reset controller
	
	checker := &CheckerCtx{
		Scheme: proto,
	}

	// Spawn goroutines for FD pool
	for range cap(c.fd_pool){
		go func() {
			defer c.ThreadCompletion()

			for {
				select {
					// Working proxy waiting to be saved
					case proxy := <-c.fd_pool:
						raw, err := json.Marshal(&proxy)
						if err != nil || len(raw) == 0{
							c.Abort(err) // fatal
							return
						}
						
						if _, err = FD[SaveFile].WriteString(string(raw) + "\n"); err != nil {
							c.Abort(err) // fatal
							return
						}

						c.Done()
					// Stop signal
					case <-c.ShouldStop():
						return // kill goroutine
					}
			}
		}()
	}

	// Spawn worker pool
	for range cap(c.worker_pool) {
		go func() {
			defer c.ThreadCompletion()

			for {
				select {
					case proxy := <-c.worker_pool:
						if proto == Multi {
							protocols := []string{SOCKS4, SOCKS5, HTTP, HTTPS}
							done := make(chan int, len(protocols))

							for _, proc := range protocols {
								go checker.WRAP_COMPLETION(proc, proxy.proxy, done)
							}

							for range cap(done) {
								<-done
							}

						} else {
							checker.WRAP(proto, proxy.proxy)
						}

						MX.Done()
					case <-c.ShouldStop():
						return // kill goroutine
				}
			}
		}()
	}

	checker.fd_pool = c.fd_pool

	go func(){
		defer c.Cancel()
		for {
			if c.Current() == c.GetLoad() {
				runtime.EventsEmit(APP.ctx, Fire_MsgEvent, "Cleaning up all goroutines in the pool, please wait.")
				return
			}
		}
	}()

	// Push proxies that need checking
	buff := bufio.NewScanner(FD[InputFile])
	defer FD[InputFile].Seek(0, io.SeekStart) // go back to offset 0 

	for buff.Scan() {
		c.worker_pool <- struct{ proxy Proxy }{
			proxy: Proxy(buff.Text()),
		}
	}

	runtime.EventsEmit(APP.ctx, Fire_CheckerStart)

	<-c.ShouldStop()
	c.CanExit()

	runtime.EventsEmit(APP.ctx, Fire_CheckerEnd)
}