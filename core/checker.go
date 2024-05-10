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

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//todo
func (c *Controller) StartChecking(ctx context.Context, proto string) {
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

	if MX.DidStart() || MX.Current() != 0  {
		fmt.Println(MX.started, MX.current)
		err = ErrOngoingCheck
		return
	}

	// Determine if multi protocol scan
	// or just one.
	//
	// If multi, then input file requires to satisfy PROXY URI: protocol:// etc
	scheme := Scheme
	if proto != "" { scheme = NoScheme }

	// All proxy files are valid
	err = FD.Validate(scheme)
	if err != nil {return}

	defer MX.Reset()

	// Spawn goroutines for FD pool
	for range cap(MX.fd_pool) {
		go func(){
			for {
				select {
					// Working proxy waiting to be saved
					case wProxy := <-MX.fd_pool:
						fmt.Printf("%+v", wProxy)
						MX.Done(1)

					// Stop signal
					case <-MX.ShouldStop():
						return // kill goroutine
				}
			}
		}()
	}

	// Spawn worker pool 
	for range cap(MX.worker_pool) {
		go func ()  {
			init := false
			for {
				select {
					case msg := <-MX.worker_pool:
						if !init {
							MX.Start()
							init = true
						}

						MX.fd_pool <- struct{proxy Proxy; latency uint32; anonimity Anonimity}{
							proxy: msg.proxy,
						}
						// fmt.Println(proxy)

						// proxy checking needs to happen here
					case <-MX.ShouldStop():
						return // kill goroutine
				}
			}	
		}()
	}
	
	// Push proxies that needs checking
	buff := bufio.NewScanner(FD[InputFile])
	defer FD[InputFile].Seek(0, io.SeekStart)

	go func(){
		for {
			if MX.Current() == 0 && MX.DidStart() {
				MX.Cancel()
				return
			}
		}
	}()


	for buff.Scan() {
		MX.Add(1)

		MX.worker_pool <- struct{proxy Proxy}{
			proxy: Proxy(buff.Text()),
		}
	}

	<-MX.ShouldStop()
	runtime.LogDebug(APP.ctx, "SHOULD STOP")
	
}