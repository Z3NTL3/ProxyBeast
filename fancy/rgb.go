package fancy

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import "fmt"

type rgb []string

var (
	Palettes rgb = []string{
		colorize(125),
		colorize(126),
		colorize(127),
		colorize(128),
		colorize(129),
		colorize(130),
		colorize(134),
		colorize(135),
		colorize(146),
		colorize(147),
		colorize(206),
		colorize(207),
		colorize(219),
		colorize(218),
		colorize(168),
		colorize(169),
	}
	Bold string = "\033[1m"
	Reset string = "\033[0m"
	Endline string = "\r\n"
)

func colorize(color int) (string){
	return fmt.Sprintf("\033[38;5;%dm", color)
}