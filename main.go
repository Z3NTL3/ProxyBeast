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
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	
	"golang.org/x/sync/errgroup"
)

func checkArgs(cliArgs *[]string) (validity bool) {
	file_regex := regexp.MustCompile(`^.*\.txt$`)
	number_regex := regexp.MustCompile(`^[0-9]+$`)

	validity = true

	if len(*(cliArgs)) != 3 {
		validity = false
		return
	}
	// protocol arg
	switch strings.ToLower((*cliArgs)[0]) {
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

	if !file_regex.MatchString((*cliArgs)[1]) {
		validity = false
	}

	if !number_regex.MatchString((*cliArgs)[2]) {
		delay,_ := strconv.Atoi((*cliArgs)[0])
		globals.Timeout = delay
		validity = false
	}

	return
}

func main() {
	builder.Logo()
	args := os.Args[1:]

	group := new(errgroup.Group)
	max_worker_count := runtime.NumCPU()
	free_cores := 3
	
	runtime.GOMAXPROCS((max_worker_count-free_cores))
	group.SetLimit(10000)
	
	
	validity := checkArgs(&args); if !validity {
		handlers.Err("Invalid command line arguments. Example usage: ./proxy-checker.exe <protocol> <proxyFile.txt> <timeout>")
		os.Exit(-1)
	}
	path := args[1]
	
	scanner, err := filesystem.LineByLine_Scanner(&path); if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}
	
	err = filesystem.TruncateAtStart(); if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

	for scanner.Scan() {
		text := scanner.Text()
		group.Go(func()(error){
			return proxy.CheckProxy(text)
		})
	}
	err = group.Wait(); if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}
	
}
