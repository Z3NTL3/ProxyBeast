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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// todo nothing final here
func (c *Controller) StartScan(ctx context.Context, proto string) {
	var err error

	// if err fire to js runtime
	defer func(err_ *error) {
		if *err_ != nil {
			runtime.EventsEmit(APP.ctx, Fire_ErrEvent, (*err_).Error())
		}
	}(&err)

	selection := Proxy(proto)

	// check if given proto is known
	if selection != "" && !(&selection).IsValid(Scheme) {
		runtime.EventsEmit(APP.ctx, Fire_ProtoUnknown)
		return
	}

	// determine if a scan is ongoing
	if c.CurrentThread() != 0 {
		err = ErrOngoingCheck
		return
	}

	// determine if multi protocol scan
	isMulti := Scheme
	if proto != "" {
		isMulti = NoScheme
	}

	// validity check
	err = FD.Validate(isMulti)
	if err != nil {
		return
	}

	runtime.EventsEmit(APP.ctx, Fire_CheckerTotalLoad, c.load)
	c.threads = int32(cap(c.fd_pool) + cap(c.worker_pool))

	defer c.Reset() // reset controller

	checker := &Checker{}

	// Spawn goroutines for FD pool
	for range cap(c.fd_pool){
		go func() {
			defer c.ThreadCompletion()

			for {
				select {
					// Working proxy waiting to be saved
					case proxy := <-c.fd_pool:
						raw, err := json.Marshal(&proxy)
						if err != nil {
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
						start := time.Now().UnixMilli()

						//todo
						level, err := checker.SOCKS4(proxy.proxy)
						if err != nil { 
							c.Done()
							continue
						}

						latency := fmt.Sprintf("%dms", time.Now().UnixMilli() - start)
						go func(){
							defer recover()
							c.fd_pool <- FD_Pool{
								Proxy: Proxy("test"),
								Latency: latency,
								Anonimity: level,
							}
						}()
						
					case <-c.ShouldStop():
						return // kill goroutine
				}
			}
		}()
	}

	go func(){
		defer (*c.cancel)()
		for {
			if c.Current() == c.GetLoad() {
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