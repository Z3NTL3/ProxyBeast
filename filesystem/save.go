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
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func WriteToSaveFile(data string) error {
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, "\r", "")

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	basePath, err := filepath.Abs(cwd)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(basePath, "saves", "goods.txt"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()
	return nil
}
