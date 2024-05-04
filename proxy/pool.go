package proxy

var Current uint32

type Workers chan struct {
	proxy Proxy
}

type FD_Pool chan struct {
	proxy     Proxy
	latency   uint32
	anonimity Anonimity
}