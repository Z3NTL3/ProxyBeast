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
	"fmt"
	"regexp"
	"strings"
)

type (
	Proxy     string
)

const (
	SOCKS4 string = "socks4"
	SOCKS5 string = "socks5"
	HTTP   string = "http"
	HTTPS  string = "https"
	Multi string = "multi"
)

var (
	PROXY_URI = regexp.MustCompile(`^(?:(\w+):\/\/)?(?:([\w%]+)(?::([\w%]+))?@)?((?:\d{1,3}\.){3}\d{1,3}:\d+)$`)
)

func(p *Proxy) IsValid() bool {
	return PROXY_URI.MatchString(string(*p))
}


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