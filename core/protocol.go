/*
     ProxyBeast GUI

The ultimate proxy checker
       by @z3ntl3

    [proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App core]
*/

package core

import (
	"fmt"
	"strings"
)

const (
	SOCKS4 string = "socks4"
	SOCKS5 string = "socks5"
	HTTP   string = "http"
	HTTPS  string = "https"
)

func (p *Proxy) IsSOCKS4() bool {
	return p.protocol(SOCKS4)
}

func (p *Proxy) IsSOCKS5() bool {
	return p.protocol(SOCKS5)
}

func (p *Proxy) IsHTTP() bool {
	return p.protocol(HTTP)
}

func (p *Proxy) IsHTTPS() bool {
	return p.protocol(HTTPS)
}

func (p *Proxy) protocol(scheme string) bool {
	return strings.Contains(string(*p), fmt.Sprintf("%s://", scheme))
}