/*
	ProxyBeast GUI

The ultimate proxy checker

	   by @z3ntl3

	[proxy.pix4.dev]

License: GNU
Note: Please do give us a star on Github, if you like ProxyBeast

[App core]
*/
package core

import (
	"bufio"
	"io"
	"os"
)

type Filesystem map[string]*os.File

var FD Filesystem = map[string]*os.File{
	InputFile: nil,
	SaveFile: nil,
}

func OpenFileRDO(loc string) (*os.File, error) {
	return os.OpenFile(loc, os.O_APPEND, os.ModeAppend)
}

func (fs *Filesystem) Validate(scheme bool) (err error) {
	defer func(){
		if err != ErrFDNil {
			(*fs)[InputFile].Seek(0, io.SeekStart)
		}
	}()
	if (*fs)[InputFile] == nil || (*fs)[SaveFile] == nil {
		err = ErrFDNil
		return
	}

	buff := bufio.NewScanner((*fs)[InputFile])

	load := 0
	for buff.Scan() {
		proxy := Proxy(buff.Text())
		if !(&proxy).IsValid(scheme) {
			err = ErrInvalidProxyURI
			break
		}
		load++
	}

	if load < 1 {
		err = ErrNoProxiesFound
		return
	}

	MX.SetLoad(uint32(load))
	return
}