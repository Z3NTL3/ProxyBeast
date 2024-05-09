package core

import "os"

func OpenFileRDO(loc string) (*os.File, error) {
	return os.Open(loc)
}