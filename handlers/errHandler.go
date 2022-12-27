package handlers

import (
	"fmt"
)

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

func Err(details string) {
	fmt.Printf("\033[1m\033[38;5;218m[\033[38;5;196m%s\033[38;5;218m\033[1m]\033[0m \033[1m\033[38;5;168m%s\033[0m\r\n", "Error", details)
}