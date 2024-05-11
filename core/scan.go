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

	if MX.DidStart() || MX.Current() != 0  {
		err = ErrOngoingCheck
		return
	}

	// Determine if multi protocol scan
	isMulti := Scheme
	if proto != "" { isMulti = NoScheme }

	// All proxy files are valid
	err = FD.Validate(isMulti)
	if err != nil {return}

	runtime.EventsEmit(APP.ctx, Fire_CheckerTotalLoad, MX.cntr_load)
	defer MX.Reset()

	checker := &Checker{
		Scheme: proto,
		Multi: isMulti,
	}

	// Spawn goroutines for FD pool
	for range cap(MX.fd_pool) {
		go func(){
			for {
				select {
					// Working proxy waiting to be saved
					case data := <-MX.fd_pool:
						MX.Done(1)
						raw, err := json.Marshal(&data)
						if err != nil {
							runtime.LogError(APP.ctx, err.Error())
							return
						}

						_, err = FD[SaveFile].WriteString(string(raw)+"\n")
						if err != nil {
							runtime.LogError(APP.ctx, err.Error())
						}
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
							runtime.EventsEmit(APP.ctx, Fire_CheckerStart)
						}

						data := FD_Pool{
							Proxy: msg.proxy,
						}
						start := time.Now().UnixMilli()
						if checker.Scheme == SOCKS4{
							anon, err := checker.SOCKS4(msg.proxy)
							if err != nil { 
								MX.Done(1)
							} else {
								data.Anonimity = anon
								end := time.Now().UnixMilli()
								latency := fmt.Sprintf("%dms", end - start)
	
								data.Latency = latency
								MX.fd_pool <- data
							}
						}

						
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
	runtime.EventsEmit(APP.ctx, Fire_CheckerEnd)
}