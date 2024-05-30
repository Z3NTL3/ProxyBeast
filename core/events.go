/*
	ProxyBeast GUI

The ultimate proxy checker

	   by @z3ntl3

	[proxy.pix4.dev]

License: GNU
Note: If you've liked ProxyBeast, please consider starring our Github repository.

[App core]
*/
package core

const (
	Fire_MsgEvent      Operation = "msg"
	Fire_ErrEvent      Operation = "error"
	Fire_ErrSvdirEvent Operation = "svdir_failure"
	Fire_ProtoUnknown  Operation = "proto_unknown"
	Fire_CheckerTotalLoad Operation = "checker_load"
	Fire_CheckerStart Operation = "checker_start"
	Fire_CheckerEnd Operation = "checker_end"
	Fire_CurrentThread Operation = "current_thread"
	Fire_FatalError Operation = "fatal_error"
	Fire_ProxyData Operation = "proxy_data"

	// go runtime event listeners
	OnStartScan Operation = "scan"
	OnDialog    Operation = "dialog"
	OnSettingsModifyTimeout Operation = "modify_timeout" 
	OnCancelScan Operation = "cancel_scan" 
)
