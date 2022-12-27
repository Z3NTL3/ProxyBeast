package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
)

func Configure(protocol, proxy *string) (*http.Transport, error) {
	proxyUrl, err := url.Parse(fmt.Sprintf("%s://%s",*protocol ,*proxy))
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy:             http.ProxyURL(proxyUrl),
		ForceAttemptHTTP2: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return transport, nil
}