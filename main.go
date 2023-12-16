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
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"Z3NTL3/proxy-checker/builder"
	"Z3NTL3/proxy-checker/cmd"
	"Z3NTL3/proxy-checker/filesystem"
	"Z3NTL3/proxy-checker/globals"
	"Z3NTL3/proxy-checker/handlers"
	"Z3NTL3/proxy-checker/proxy"

	"golang.org/x/sync/errgroup"
)

func checkArgs(timeout int, protocol, proxyfile string) (validity bool) {
	file_regex := regexp.MustCompile(`^.*\.txt$`)
	validity = true

	if !file_regex.MatchString(proxyfile) {
		validity = false
	}

	switch strings.ToLower(protocol) {
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

	if timeout <= 0 {
		log.Fatal("timeout cannot be 0 or a negative value")
		validity = false
	}

	return
}

/*
ONLY USED WHEN rotating option is set;

used for backbone proxy checking
*/
func getip() error {
	res, err := http.Get("http://ip-api.com/json")
	if err != nil {
		return err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	data := struct {
		Query string `json:"query"`
	}{}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if len(data.Query) <= 1 {
		return errors.New("Could not obtain local IP")
	}

	*(&globals.LocalIP) = data.Query
	return nil
}

func main() {
	builder.Logo()
	cmd.Init()

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if *&globals.Rotating {
		if err := getip(); err != nil {
			log.Fatalf("[COULD NOT OBTAIN LOCAL IP]: %s", err)
		}
	}

	if slices.Contains(os.Args, "--help") ||
		slices.Contains(os.Args, "-h") {
		os.Exit(0)
	}

	group := new(errgroup.Group)
	max_worker_count := runtime.NumCPU()
	free_cores := 3

	runtime.GOMAXPROCS((max_worker_count - free_cores))
	group.SetLimit(-1)

	validity := checkArgs(globals.Timeout, globals.Protocol, globals.ProxyFile)
	if !validity {
		handlers.Err("Invalid command line arguments. Get usage info by passing -h flag!")
		os.Exit(-1)
	}
	path := *&globals.ProxyFile

	scanner, err := filesystem.LineByLine_Scanner(&path)
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

	err = filesystem.RecursiveInit()
	if err != nil {
		handlers.Err(err.Error())
		os.Exit(-1)
	}

	for scanner.Scan() {
		text := scanner.Text()
		time.Sleep(50 * time.Millisecond)
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
