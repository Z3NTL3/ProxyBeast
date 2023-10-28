package config

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"h12.io/socks"
)

func Configure(protocol, proxy *string) (*http.Transport, error) {
	proxyUrl, err := url.Parse(fmt.Sprintf("%s://%s", *protocol, strings.TrimSpace(*proxy)))
	if err != nil {
		return nil, err
	}
	transport := new(http.Transport)

	switch true {
	case strings.Contains(proxyUrl.Scheme, "socks4"):
		transport.DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
			f := socks.Dial(proxyUrl.String())
			return f(network, addr)
		}
	default:
		// http.ProxyURL func can handle http, https and socks5
		transport = &http.Transport{
			Proxy:             http.ProxyURL(proxyUrl),
			ForceAttemptHTTP2: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	}

	return transport, nil
}
