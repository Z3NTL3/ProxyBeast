package proxy

import "net/url"

const (
	SOCKS4 string = "socks4"
	SOCKS5 string = "socks5"
	HTTP   string = "http"
	HTTPS  string = "https"
)

func (p *Proxy) IsSOCKS4() bool {
	match, err := p.protocol(SOCKS4)
	if err != nil {
		panic("proxy uri could not be parsed")
	}

	return match
}

func (p *Proxy) IsSOCKS5() bool {
	match, err := p.protocol(SOCKS5)
	if err != nil {
		panic("proxy uri could not be parsed")
	}

	return match
}

func (p *Proxy) IsHTTP() bool {
	match, err := p.protocol(HTTP)
	if err != nil {
		panic("proxy uri could not be parsed")
	}

	return match
}

func (p *Proxy) IsHTTPS() bool {
	match, err := p.protocol(HTTPS)
	if err != nil {
		panic("proxy uri could not be parsed")
	}

	return match
}

func (p *Proxy) protocol(scheme string) (match bool, err error) {
	uri, err := url.Parse(string(*p))
	if err != nil {
		return
	}

	match = (uri.Scheme == scheme)
	return
}