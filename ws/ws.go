package wsclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Ws struct {
	endpoint    string
	conn        *websocket.Conn
	reConnectFn func() error
	option      *WsOption
	lock        sync.Mutex
	isReConning bool
	writeMsg    chan []byte
	readMsg     chan []byte
	exitSign    chan string
	heartTime   time.Duration
}

func NewWs(endpoint string) (*Ws, error) {
	conn, err := newWebsocketConn(endpoint)
	if err != nil {
		return nil, err
	}
	ws := &Ws{
		conn:        conn,
		endpoint:    endpoint,
		writeMsg:    make(chan []byte, 1000),
		readMsg:     make(chan []byte, 1000),
		exitSign:    make(chan string, 1),
		isReConning: false,
		option:      DefaultOption(),
		heartTime:   60 * time.Second,
	}
	return ws, nil
}

func (ws *Ws) Run() {
	go ws.write()
	go ws.read()

}

func (ws *Ws) WriteMessage(msg []byte) {
	ws.writeMsg <- msg

}

func (ws *Ws) WriteObj(obj interface{}) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	ws.writeMsg <- bytes
	return nil
}

func (ws *Ws) ReadMessage() chan []byte {
	return ws.readMsg
}

func (ws *Ws) Close() error {
	close(ws.exitSign)
	err := ws.conn.Close()
	if err != nil {
		return err
	}
	return nil

}

func (ws *Ws) write() {

	ticker := time.NewTicker(ws.heartTime)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ws.writeMsg:
			if ok {
				err := ws.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {

				}
			}

		case t := <-ticker.C:
			err := ws.conn.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				reErr := ws.ReConnect()
				if reErr != nil {
				}
			}

		case <-ws.exitSign:
			err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {

			}
			return
		}
	}
}

func (ws *Ws) read() {

	for {
		mType, message, err := ws.conn.ReadMessage()
		if err != nil {

		}
		switch mType {
		case websocket.TextMessage:
			ws.readMsg <- message
		case websocket.PongMessage:

		case websocket.CloseMessage:
			ws.Close()
			return
		}
		select {
		case <-ws.exitSign:
			return
		default:

		}

	}

}

func (ws *Ws) IsReConning() bool {
	return ws.isReConning
}

func (ws *Ws) SetReConnFn(fn func() error) {
	ws.reConnectFn = fn
}

func (ws *Ws) ReConnect() error {
	if !ws.option.CanReConn {
		return nil
	}
	if ws.IsReConning() {
		return nil
	}
	ws.lock.Lock()
	defer ws.lock.Unlock()
	// 需要再次判断，并发时防止多次重连
	if ws.IsReConning() {
		return nil
	}
	ws.isReConning = true
	conn, err := newWebsocketConn(ws.endpoint)
	if err != nil {
		return fmt.Errorf("new conn error: %v", err)
	}
	ws.conn = conn
	if ws.reConnectFn != nil {
		_ = ws.reConnectFn()
	}
	return nil

}

func newWebsocketConn(endpoint string) (*websocket.Conn, error) {
	urls := strings.Split(endpoint, "/")
	if len(urls) < 3 {
		return nil, fmt.Errorf("ws endpoint format error: example wss://127.0.0.1:9944/conn: %v", endpoint)
	}
	scheme := strings.TrimSuffix(urls[0], ":")
	path := "/"
	if len(urls) == 4 {
		path = fmt.Sprintf("/%s", urls[3])
	}
	u := url.URL{Scheme: scheme, Host: urls[2], Path: path}
	conn, _, err := websocket.DefaultDialer.DialContext(context.Background(), u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("ws dial error: %v", err)
	}
	return conn, nil
}
