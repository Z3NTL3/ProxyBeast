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
	"log"
	"os"

	"Z3NTL3/proxy-checker/globals"
)

func LineByLine_Scanner(path *string) (*bufio.Scanner, error) {
	file, err := os.Open(*path)
	if err != nil {
		return nil, err
	}
	s, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if s.Size() == int64(0) {
		log.Fatalf("zero proxies in your proxy file: '%s'!", *&globals.ProxyFile)
	}
	// dont close file, we give it back as pointer to use it
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	return scanner, nil
}
