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

const (
	Fire_MsgEvent      Operation = "msg"
	Fire_ErrEvent      Operation = "error"
	Fire_ErrSvdirEvent Operation = "svdir_failure"
	Fire_ProtoUnknown  Operation = "proto_unknown"

	// go runtime event listeners
	OnStartScan Operation = "scan"
	OnDialog    Operation = "dialog"
)