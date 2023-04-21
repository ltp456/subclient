package wsclient

import "time"

type WsOption struct {
	AutoReConn   bool
	RetryConnNum int
	KeepAlive    bool
	HeartTime    time.Duration
}

func DefaultOption() *WsOption {
	return &WsOption{
		AutoReConn:   true,
		RetryConnNum: 3,
		KeepAlive:    true,
		HeartTime:    60 * time.Second,
	}
}
