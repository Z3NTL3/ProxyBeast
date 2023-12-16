package proxy

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"Z3NTL3/proxy-checker/builder"
	"Z3NTL3/proxy-checker/config"
	"Z3NTL3/proxy-checker/filesystem"
	"Z3NTL3/proxy-checker/globals"
	"Z3NTL3/proxy-checker/typedefs"
)

type host struct {
	Origin string `json:"origin"`
}

func CheckProxy(
	proxy string,
) error {
	proc := (&globals.Protocol)
	retries := globals.Retries
	var retryTimes int
	workingProxy := false

	transport, err := config.Configure(proc, &proxy)
	if err != nil {
		builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}

	seconds := time.Duration(globals.Timeout) * time.Second
	client := http.Client{
		Transport: transport,
		Timeout:   seconds,
	}

	var reader io.ReadCloser // auto close
	Req := new(http.Request)
	Resp := new(http.Response)

	all := []string{*proc}
	if globals.Multi {
		for k := range globals.Locations {
			if k == "http" {
				continue
			}
			all = append(all, k)
		}
	}

	for _, v := range all {
		transport, err := config.Configure(&v, &proxy)
		if err != nil {
			builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
			continue
		}
		client.Transport = transport

		for i := 0; i < retries; i++ {
			if retryTimes > retries || workingProxy {
				break
			}
			req, err := http.NewRequest("GET", "https://httpbin.org/ip", reader)
			if err != nil {
				retryTimes++
				continue
			}
			req.Header.Add("user-agent", typedefs.UA[rand.Intn(len(typedefs.UA))])
			req.Header.Add("content-type", "application/json")
			req.Header.Add("referer", typedefs.Referers[rand.Intn(len(typedefs.Referers))])

			resp, err := client.Do(req)
			if err != nil {
				builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mDIDNT RESPOND RETRYING! \033[1m\033[38;5;127m[\033[0m%d\033[1m\033[38;5;196m/\033[0m%d\033[1m\033[38;5;127m]\033[0m", proxy, retryTimes, retries), "\n")

				retryTimes++
				continue
			}
			defer resp.Body.Close()

			workingProxy = true
			Req = req
			Resp = resp
			break
		}

		if retryTimes > retries || !workingProxy {
			builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID \033[1m\033[38;5;127m[\033[0m%d\033[1m\033[38;5;196m/\033[0m%d\033[1m\033[38;5;127m]\033[0m", proxy, retryTimes, retries), "\n")
			continue
		}

		var jsonData host
		err = json.NewDecoder(Resp.Body).Decode(&jsonData)
		if err != nil && err != io.EOF {
			builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
			continue
		}
		if Req == nil {
			return nil
		}

		if *&globals.Rotating && jsonData.Origin != *&globals.LocalIP {
			filesystem.WriteToSaveFile(proxy, globals.Locations[v])
			builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s [%s] \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;118mVALID [ROTATING]\033[0m ", proxy, jsonData.Origin), "\n")
			return nil
		} else if jsonData.Origin == strings.Split(proxy, ":")[0] {
			filesystem.WriteToSaveFile(proxy, globals.Locations[v])
			builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;118mVALID\033[0m ", proxy), "\n")
			return nil
		}

		builder.Log("INFO", "\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")

	}

	return nil
}
