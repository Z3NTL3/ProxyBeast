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
	timeout   = flag.Int("timeout", 5, "Set the timeout in seconds using a custom value")
	protocol  = flag.String("protocol", "", "Specify a required protocol, options include: http, https, socks4, and socks5")
	proxyFile = flag.String("file", "proxies.txt", "Determine the name of the file containing proxies, must be in .txt format")
	retries   = flag.Int("retry", 1, "Specify the number of attempts to reconnect to a failed proxy")
	threads   = flag.Int("threads", runtime.NumCPU(), "Choose the number of threads to use for checking proxies, default is the number of CPU cores available")
)

func areArgsValid(timeout *int, protocol, proxyfile *string, threads *int) (areValid bool) {
	fileRegex := regexp.MustCompile(`^.*\.txt$`)
	numberRegex := regexp.MustCompile(`^[0-9]+$`)
	areValid = true

	if !fileRegex.MatchString((*proxyfile)) {
		areValid = false
	}
	if !numberRegex.MatchString(strconv.Itoa(*timeout)) {
		areValid = false
	}
	if !numberRegex.MatchString(strconv.Itoa(*threads)) {
		areValid = false
	}
	if *protocol != "" {
		if !(strings.EqualFold(*protocol, "http") ||
			strings.EqualFold(*protocol, "https") ||
			strings.EqualFold(*protocol, "socks4") ||
			strings.EqualFold(*protocol, "socks5")) {
			areValid = false
		}
	}

	return
}

func main() {
	flag.Parse()

	builder.Logo()

	group := new(errgroup.Group)

	runtime.GOMAXPROCS((*threads))
	group.SetLimit(-1)

	if !areArgsValid(timeout, protocol, proxyFile, threads) {
		handlers.Err("Invalid arguments. Get usage info by passing -h flag!")
		os.Exit(-1)
	}

	scanner, err := filesystem.LineByLine_Scanner(proxyFile)
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
		proxyToCheck := scanner.Text()
		group.Go(func() error {
			return proxy.CheckProxy(&proxyToCheck, timeout, protocol, retries)
		})
	}

	err = group.Wait()
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

}
