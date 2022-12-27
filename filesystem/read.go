package filesystem

/*
*
*   Author: Z3NTL3 (aka Efdal)
*   License: GNU
*   Telegram: @z3ntl3
*   Description: Super-duper fast and accurate proxy checker amplified with Goroutines
*
 */

import (
	"bufio"
	"os"
)

func LineByLine_Scanner(path *string) (*bufio.Scanner,error) {
	file, err := os.Open(*path); if err != nil {
		return nil, err
	}
	// dont close file, we give it back as pointer to use it
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return scanner, nil
}