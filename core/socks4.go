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

	socks "github.com/z3ntl3/socks/client"
)

func (c *CheckerCtx) SOCKS4(proxy Proxy) (anonimity string, err error) {
	if !c.Multi || c.Scheme == "socks4" {
		proxy = Proxy(fmt.Sprintf("socks4://%s", proxy))
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
		Port:     443,
	}

	proxyCtx := socks.Context{
		Resolver: net.ParseIP(uri.Hostname()),
		Port:     port,
	}

	client, err := socks.New(&socks.Socks4Client{}, targetCtx, proxyCtx)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	if err = socks.Connect(client, ctx); err != nil {
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