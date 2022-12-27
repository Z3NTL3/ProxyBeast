package builder

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import "fmt"

func Log(log_type,color, info,delim string) {
	fmt.Printf("\033[1m\033[38;5;218m[\033[0m\033[1m\033[38;5;146m%s\033[0m\033[38;5;218m\033[1m]\033[0m \033[1m%s%s\033[0m%s", log_type, color,info,delim)
}