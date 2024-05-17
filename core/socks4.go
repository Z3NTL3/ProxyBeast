/*
	ProxyBeast GUI

The ultimate proxy checker

	   by @z3ntl3

	[proxy.pix4.dev]

License: GNU
Note: If you've liked ProxyBeast, please consider starring our Github repository.

[App core]
*/
package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func (c *CheckerCtx) SOCKS4(proxy Proxy) (anonimity string, err error) {
	if !c.Multi || c.Scheme == SOCKS4 {
		proxy = Proxy(fmt.Sprintf("%s://%s", SOCKS4, proxy))
	}

	uri, err := url.Parse(string(proxy))
	if err != nil {
		return
	}

	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return
	}

	addr, err := proxifier.LookupHost(AppSettings.Store.Judge.Hostname())
	if err != nil {
		return
	}

	targetCtx := proxifier.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     443,
	}

	proxyCtx := proxifier.Context{
		Resolver: net.ParseIP(uri.Hostname()),
		Port:     port,
	}

	client, err := proxifier.New(&proxifier.Socks4Client{}, targetCtx, proxyCtx)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), AppSettings.Store.Timeout)
	defer cancel()

	if err = proxifier.Connect(client, ctx); err != nil {
		return
	}

	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	tlsConn.SetDeadline(time.Now().Add(DefaultTimeout))

	if _, err = tlsConn.Write(
		[]byte(
			fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", 
			AppSettings.Store.Judge.Path,
			AppSettings.Store.Judge.Hostname()),
		),
	); err != nil {
		return
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		return
	}

	anon := Anonimity(string(data))
	anonimity = (&anon).GetAnonimity()

	return
}