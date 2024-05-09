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