package proxy

type (
	Anonimity string
	Proxy     string

	ProxyClient struct {
		pool chan struct {
			latency   uint32
			proxy     Proxy
			anonimity Anonimity
		} // todo props
		quit chan int
	}

	ClientImpl interface {
		*ProxyClient
		Check() // todo
	}
)
