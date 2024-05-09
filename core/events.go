package core

const (
	Fire_MsgEvent        Operation = "msg"
	Fire_ErrEvent        Operation = "error"
	Fire_ErrSvdirEvent   Operation = "svdir_failure"
	Fire_ProxyListLoaded Operation = "list_loaded"
	Fire_SaveFileLoaded  Operation = "sf_loaded"

	// go runtime event listeners
	OnStartScan Operation = "scan"
	OnDialog    Operation = "dialog"
)