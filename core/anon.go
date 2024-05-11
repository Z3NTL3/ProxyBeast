package core

import "strings"

type Anonimity string

const (
	XForwardedFor   string = "x-forwarded-for"
	Via             string = "via"
	ProxyConnection string = "proxy-connection"

	Elite       string = "elite"
	Anonymous   string = "anonymous"
	Transparent string = "transparent"
)

func (a *Anonimity) GetAnonimity() string {
	*a = Anonimity(strings.ToLower(string(*a)))
	
	if  a.IsElite() {return Elite}
	if a.IsAnonymous() {return Anonymous}

	return Transparent
}

func (a *Anonimity) IsElite() bool {
	return !a.Contains(XForwardedFor, Via, ProxyConnection)
}

func (a *Anonimity) IsAnonymous() bool {
	via := a.Contains(Via, ProxyConnection)
	forwardFor := a.Contains(XForwardedFor)

	if via && !forwardFor {
		return true
	}

	return false
}

func (a *Anonimity) IsTransparent() bool {
	return a.Contains(XForwardedFor, Via, ProxyConnection)
}


func (a *Anonimity) Contains(header ...string) bool {
	for _, h := range header {
		if strings.Contains(string(*a), h) {
			return true
		}
	}
	return false
} 