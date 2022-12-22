package wsclient

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var ws *Ws
var err error

func init() {
	endpoint := "ws://127.0.0.1:9944/"
	//endpoint := "wss://testnet.binance.vision/ws"
	//endpoint := "wss://rpc.polkadot.io/rpc"
	ws, err = NewWs(endpoint)
	if err != nil {
		panic(err)
	}
}

func TestUrl(t *testing.T) {
	res := "ws://127.0.0.1:9944"
	splits := strings.Split(res, "/")
	fmt.Println(len(splits), splits)
}

func TestWs(t *testing.T) {

	msg := ws.ReadMessage()
	go func() {
		for {
			select {
			case info := <-msg:
				fmt.Println("receiver msg: ", string(info))
			}
		}
	}()
	go func() {
		ticker := time.NewTicker(50 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sendMsg := `{"id":1661499051875090000,"jsonrpc":"2.0","method":"state_getStorage","params":["26aa394eea5630e07c48ae0c9558cef7b99d880ec681799c0cf30e8886371da9de1e86a9a8c739864cf3cc5ec2bea59fd43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d","0x72efed7544183274bcdcd6da8a1f551aa9151cc72677a9852857aa0bb74cd068"]}`
				//sendMsg := `{"method":"SUBSCRIBE","params":["btcusdt@aggTrade","btcusdt@depth"],"id":1}`
				ws.WriteMessage([]byte(sendMsg))

			}
		}
	}()
	ws.Run()
	time.Sleep(5 * time.Second)
	//ws.Exit()
	ws.Close()

	time.Sleep(10 * time.Minute)

}
