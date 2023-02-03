package globals

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

var (
	Protocols = [...]string{
		"http",
		"https",
		"socks4",
		"socks5",
	}
	Multi bool = false // default, false, if set true through cli
	// then it will check for all protocols
)
