package proxy

import (
	"reflect"
	"time"
)

type (
	timeout time.Duration
	filePath string
)
var (
	Timeout timeout
	ProxyFilePath filePath
)

func(t *timeout) String() string {
	return time.Duration(*t).String()
}

func (t *timeout) Set(input string) error {
	t_, err := time.ParseDuration(input)
	if err != nil { return err }

	*t = timeout(t_)
	
	return nil
}

func(t *timeout) Type() string {
	return reflect.TypeOf(*t).String()
}

func(t *filePath) String() string {
	return string(*t)
}

func (t *filePath) Set(input string) error {
	*t = filePath(input)
	return nil
}

func(t *filePath) Type() string {
	return reflect.TypeOf(*t).String()
}