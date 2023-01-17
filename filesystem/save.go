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
	file, err := os.OpenFile(path.Join(basePath, "saves", "goods.txt"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("%s\r\n", data))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
