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
	"os"
	"path"
	"path/filepath"
)

var (
	loc = [...]string{
		"goods-http.txt",
		"goods-https.txt",
		"goods-socks4.txt",
		"goods-socks5.txt",
	}
)

func RecursiveInit() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath, err := filepath.Abs(cwd)
	if err != nil {
		return err
	}

	for _, v := range loc {
		file, err := os.OpenFile(path.Join(basePath, "saves", v), os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}
