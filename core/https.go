package core

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Z3NTL3/proxifier"
)

func (c *CheckerCtx) HTTPS(proxy Proxy) (anonimity string, err error) {
	if !proxy.IsHTTP() {
		proxy = Proxy(fmt.Sprintf("%s://%s", HTTP, proxy))
	}

	uri, err := url.Parse(string(proxy))
	if err != nil {
		return
	}

	port, err := strconv.Atoi(uri.Port())
	if err != nil {
		return
	}

	auth := proxifier.Auth{}
	if uri.User != nil {
		if username := uri.User.Username(); username != "" {
			auth.Username = username
		}
		if passwd, canUse := uri.User.Password(); canUse {
			auth.Password = passwd
		}
	}

	httpClient := proxifier.HTTPClient{
		Auth: auth,
		TLS: true,
	}

	conn, err := httpClient.PROXY(AppSettings.Store.Judge.String(), proxifier.Context{
		Resolver: net.ParseIP(uri.Hostname()),
		Port:     port,
	}, DefaultTimeout)
	if err != nil || conn == nil {
		return
	}
	
	if err = conn.SetDeadline(time.Now().Add(DefaultTimeout)); err != nil {return}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
	})

	if _, err = tlsConn.Write([]byte(
		fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
			AppSettings.Store.Judge.Path,
			AppSettings.Store.Judge.Hostname(),
		))); err != nil {
		return
	}

	data, err := io.ReadAll(tlsConn)
	if err != nil {
		return
	}

	fmt.Println(string(data))

	anon := Anonimity(string(data))
	anonimity = (&anon).GetAnonimity()

	return
}