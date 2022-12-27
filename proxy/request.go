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
	"Z3NTL3/proxy-checker/builder"
	"Z3NTL3/proxy-checker/config"
	"Z3NTL3/proxy-checker/filesystem"
	"Z3NTL3/proxy-checker/globals"
	"Z3NTL3/proxy-checker/typedefs"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type host struct {
	Origin string `json:"origin"`
}

func CheckProxy(
	proxy string,
) error {
	proc := globals.Protocol
	transport, err := config.Configure(&proc, &proxy); if err != nil{
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}

	seconds :=  time.Duration(globals.Timeout)*time.Second
	client := http.Client{
		Transport: transport,
		Timeout: seconds,
	}

	var reader io.ReadCloser // auto close
	req, err := http.NewRequest("GET","https://httpbin.org/ip", reader); if err != nil {
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}
	req.Header.Add("user-agent", typedefs.UA[rand.Intn(len(typedefs.UA))])
	req.Header.Add("content-type", "application/json")
	req.Header.Add("referer", typedefs.Referers[rand.Intn(len(typedefs.Referers))])

	resp, err := client.Do(req); if err != nil {
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}
	defer resp.Body.Close()

	var jsonData host
	err = json.NewDecoder(resp.Body).Decode(&jsonData); if err != nil {
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}

	
	proxyUrl, err := transport.Proxy(req); if err != nil {
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxy), "\n")
		return nil
	}

	if fmt.Sprintf("%s:%s",jsonData.Origin, proxyUrl.Port()) == proxyUrl.Host {
		filesystem.WriteToSaveFile(proxyUrl.Host)
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;118mVALID\033[0m ", proxyUrl.Host), "\n")
	} else {
		builder.Log("INFO","\033[38;5;127m", fmt.Sprintf("\033[38;5;126m\033[1mProxy\033[0m\033[1m\033[38;5;127m[\033[38;5;147m %s \033[0m\033[1m\033[38;5;127m] \033[1m\033[38;5;196mINVALID\033[0m", proxyUrl.Host), "\n")
	
	}

	return nil
}