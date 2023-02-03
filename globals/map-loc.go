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
	Locations = map[string]string{
		"http":   "goods-http.txt",
		"https":  "goods-https.txt",
		"socks4": "goods-socks4.txt",
		"socks5": "goods-socks5.txt",
	}
)
