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

	if MX.CurrentThread() != 0 {
		fmt.Println("current", MX.cthread)
		err = ErrOngoingCheck
		return
	}

	// Determine if multi protocol scan
	isMulti := Scheme
	if proto != "" {
		isMulti = NoScheme
	}

	// All proxy files are valid
	err = FD.Validate(isMulti)
	if err != nil {
		return
	}

	runtime.EventsEmit(APP.ctx, Fire_CheckerTotalLoad, MX.load)

	MX.Reset()
	MX.capacity = int32(cap(MX.fd_pool) + cap(MX.worker_pool))

	defer MX.Reset()


	// Spawn goroutines for FD pool
	for range cap(MX.fd_pool){
		go func() {
			defer MX.ThreadCompletion()

			for {
				select {
				// Working proxy waiting to be saved
					case <-MX.fd_pool:
						time.Sleep(time.Second * 1)
						MX.Done()
					// Stop signal
					case <-MX.ShouldStop():
						return // kill goroutine
					}
			}
		}()
	}

	// Spawn worker pool
	for range cap(MX.worker_pool) {
		go func() {
			defer MX.ThreadCompletion()

			for {
				select {
					case <-MX.worker_pool:
						time.Sleep(time.Second * 1)
						MX.fd_pool <- FD_Pool{
							Proxy: Proxy("test"),
						}

				case <-MX.ShouldStop():
					return // kill goroutine
				}
			}
		}()
	}


	// Push proxies that need checking
	buff := bufio.NewScanner(FD[InputFile])
	defer FD[InputFile].Seek(0, io.SeekStart)

	for buff.Scan() {
		MX.worker_pool <- struct{ proxy Proxy }{
			proxy: Proxy(buff.Text()),
		}
	}
	


	time.Sleep(time.Millisecond * 100)
	runtime.EventsEmit(APP.ctx, Fire_CheckerStart)

	<-MX.ShouldStop()
	fmt.Println("SHOULD STOP")
	MX.CanExit()
	fmt.Println("CAN EXIT")
	runtime.EventsEmit(APP.ctx, Fire_CheckerEnd)
}