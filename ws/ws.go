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
	reTime      time.Time

	option      *WsOption
	lock        sync.Mutex
	isReConning bool
	writeMsg    chan []byte
	readMsg     chan []byte
	exitSign    chan string
}

func NewWs(endpoint string) (*Ws, error) {
	conn, err := newWebsocketConn(endpoint)
	if err != nil {
		return nil, err
	}
	ws := &Ws{
		conn:        conn,
		endpoint:    endpoint,
		writeMsg:    make(chan []byte, 1),
		readMsg:     make(chan []byte, 1),
		exitSign:    make(chan string, 1),
		isReConning: false,
		option:      DefaultOption(),
		reTime:      time.Now(),
	}
	return ws, nil
}

func (ws *Ws) Run() {
	go ws.write()
	go ws.read()
	if ws.option.KeepAlive {
		go ws.KeepAlive()
	}
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

	for {
		select {
		case msg, ok := <-ws.writeMsg:
			if ok {
				err := ws.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {

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
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if ws.option.AutoReConn {
				ws.reConnect()
			}
			time.Sleep(300 * time.Millisecond)
			continue
		}
		ws.readMsg <- message
		select {
		case <-ws.exitSign:
			return
		default:

		}
	}

}

// KeepAlive send ping-pong
func (ws *Ws) KeepAlive() {
	ticker := time.NewTicker(ws.option.HeartTime)
	lastResponse := time.Now()
	ws.conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})
	defer ticker.Stop()
	for {
		deadline := time.Now().Add(10 * time.Second)
		err := ws.conn.WriteControl(websocket.PingMessage, []byte{}, deadline)
		if err != nil {

		}
		<-ticker.C
		if time.Since(lastResponse) > ws.option.HeartTime {
			ws.conn.Close()
			if ws.option.AutoReConn {
				ws.reConnect()
			}
		}
	}
}

func (ws *Ws) SetReConnFn(fn func() error) {
	ws.reConnectFn = fn
}

func (ws *Ws) reConnect() error {
	fmt.Printf("start re connect now %v \n", time.Now())
	if !ws.option.AutoReConn || time.Since(ws.reTime) < ws.option.HeartTime {
		return nil
	}
	if ws.isReConning {
		return nil
	}
	ws.lock.Lock()
	defer ws.lock.Unlock()

	if ws.isReConning {
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
	ws.isReConning = false
	ws.reTime = time.Now()
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
	conn.SetReadLimit(655350)
	return conn, nil
}
