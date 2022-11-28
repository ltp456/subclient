package wsclient

type WsOption struct {
	CanReConn    bool
	RetryConnNum int
}

func DefaultOption() *WsOption {
	return &WsOption{
		CanReConn:    true,
		RetryConnNum: 3,
	}
}
