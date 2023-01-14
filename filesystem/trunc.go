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

func TruncateAtStart() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath, err := filepath.Abs(cwd)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(basePath, "saves", "goods.txt"), os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
