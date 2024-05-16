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

	socks "github.com/z3ntl3/socks/client"
)

func (c *CheckerCtx) SOCKS5(proxy Proxy) (anonimity string, err error) {
	if !c.Multi || c.Scheme == SOCKS5 {
		proxy = Proxy(fmt.Sprintf("%s://%s", SOCKS5, proxy))
	}

	uri, err := url.Parse(string(proxy))
	if err != nil {
		return
	}

	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return
	}

	addr, err := socks.LookupHost(JUDGE)
	if err != nil {
		return
	}

	targetCtx := socks.Context{
		Resolver: net.ParseIP(addr[0]),
		Port:     80,
	}

	proxyCtx := socks.Context{
		Resolver: net.ParseIP(uri.Hostname()),
		Port:     port,
	}

	client, err := socks.New(&socks.Socks5Client{}, targetCtx, proxyCtx)
	if err != nil {
		return
	}

	if uri.User != nil {
		client.Auth.Username = uri.User.Username()
		if passwd, canUse := uri.User.Password(); canUse {
			client.Auth.Password = passwd
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), AppSettings.Store.Timeout)
	defer cancel()

	if err = socks.Connect(client, ctx); err != nil {
		return
	}

	defer client.Close()
	client.SetLinger(0)

	tlsConn := tls.Client(client, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err = tlsConn.Write(
		[]byte(
			fmt.Sprintf("GET /judge.php HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n", JUDGE),
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