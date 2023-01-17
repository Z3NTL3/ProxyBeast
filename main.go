package main

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
	"Z3NTL3/proxy-checker/filesystem"
	"Z3NTL3/proxy-checker/globals"
	"Z3NTL3/proxy-checker/handlers"
	"Z3NTL3/proxy-checker/proxy"
	"flag"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	Timeout   = flag.String("timeout", "5", "Set custom timeout in seconds")
	Protocol  = flag.String("protocol", "", "Required flag, can be one of http, https, socks4 or socks5")
	ProxyFile = flag.String("file", "proxies.txt", "Determines your proxy file name requires to be *.txt matching")
	Retry     = flag.Int("retry", 1, "The amount of tries to retry to connect to a failure proxy")
)

func checkArgs(timeout, protocol, proxyfile *string) (validity bool) {
	file_regex := regexp.MustCompile(`^.*\.txt$`)
	number_regex := regexp.MustCompile(`^[0-9]+$`)
	validity = true

	if !file_regex.MatchString((*proxyfile)) {
		validity = false
	}

	switch strings.ToLower(*protocol) {
	case "http":
		globals.Protocol = "http"
	case "https":
		globals.Protocol = "https"
	case "socks4":
		globals.Protocol = "socks4"
	case "socks5":
		globals.Protocol = "socks5"
	default:
		validity = false
	}

	if !number_regex.MatchString(*timeout) {
		delay, err := strconv.Atoi(*timeout)
		if err != nil {
			handlers.Err("Cannot process -timeout delay option")
			validity = false
			return
		}
		globals.Timeout = delay
		validity = false
	}

	return
}

func main() {
	builder.Logo()
	flag.Parse()
	globals.Retries = *Retry

	group := new(errgroup.Group)
	max_worker_count := runtime.NumCPU()
	free_cores := 3

	runtime.GOMAXPROCS((max_worker_count - free_cores))
	group.SetLimit(-1)

	validity := checkArgs(Timeout, Protocol, ProxyFile)
	if !validity {
		handlers.Err("Invalid command line arguments. Get usage info by passing -h flag!")
		os.Exit(-1)
	}
	path := *ProxyFile

	scanner, err := filesystem.LineByLine_Scanner(&path)
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

	err = filesystem.TruncateAtStart()
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

	for scanner.Scan() {
		text := scanner.Text()
		group.Go(func() error {
			return proxy.CheckProxy(text)
		})
	}
	err = group.Wait()
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

}
