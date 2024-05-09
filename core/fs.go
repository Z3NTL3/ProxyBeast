package core

import "os"

func ReadFromFile(loc string) (*os.File, error) {
	return os.Open(loc)
}