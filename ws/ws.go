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
	readFn      func([]byte) error
	option      *WsOption
	lock        sync.Mutex
	isReConning bool
	writeMsg    chan []byte
	readMsg     chan []byte
	exitSign    chan string
	heartTime   int64
}

func NewWs(endpoint string) (*Ws, error) {
	conn, err := newWebsocketConn(endpoint)
	if err != nil {
		return nil, err
	}
	ws := &Ws{
		conn:     conn,
		endpoint: endpoint,

		writeMsg:    make(chan []byte, 1000),
		readMsg:     make(chan []byte, 1000),
		exitSign:    make(chan string, 1),
		isReConning: false,
		option:      DefaultOption(),
		heartTime:   60,
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

func (ws *Ws) Exit() {
	close(ws.exitSign)
}

func (ws *Ws) closeConn() {
	err := ws.conn.Close()
	if err != nil {
		//log.Printf("ws conn error: %v", err)
	}
}

func (ws *Ws) write() {
	defer func() {
		//log.Printf("close websocket conn")
		ws.closeConn()
	}()

	ticker := time.NewTicker(time.Duration(ws.heartTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-ws.writeMsg:
			if ok {
				err := ws.conn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					//log.Printf("write message error: %v \n", err)
				}
			}

		case t := <-ticker.C:
			err := ws.conn.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				reErr := ws.ReConnect()
				if reErr != nil {
					//log.Printf("ws re connect error: %v", reErr)
				}
				//log.Printf("write ping message error: %v", err)
			}
			//log.Printf("websocket ping message")
		case <-ws.exitSign:
			err := ws.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				//log.Printf("ws exit  write close message error: %v", err)
			}
			//log.Printf("rev ws write exit sign")

			return
		}
	}
}

func (ws *Ws) read() {

	for {

		mType, message, err := ws.conn.ReadMessage()
		if err != nil {
			time.Sleep(time.Duration(ws.heartTime) * time.Second)
			//log.Printf("ws read message error: %v", err)
		}
		switch mType {
		case websocket.TextMessage:
			ws.readMsg <- message
		case websocket.PongMessage:
			//log.Printf("websocket rec pong message")
		case websocket.CloseMessage:
			//log.Printf("websocket close message")
			ws.Exit()
			return
		}
		select {
		case <-ws.exitSign:
			//log.Printf("rev ws read exit sign")
			return
		default:

		}

	}

}

func (ws *Ws) IsReConning() bool {
	return ws.isReConning
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
