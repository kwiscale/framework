package kwiscale

import (
	"log"

	"github.com/gorilla/websocket"
)

// IWSHander is the base template to implement to be able to use
// Websocket
type IWSHandler interface {
	// Serve is the method to implement inside the project
	Serve()
	upgrade() error
}

// WebsockerHandler type
type WebSocketHandler struct {
	BaseHandler
	conn *websocket.Conn
}

// upgrade protocol to use websocket communication
func (ws *WebSocketHandler) upgrade() error {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	var err error
	ws.conn, err = upgrader.Upgrade(ws.Response, ws.Request, nil)
	return err
}

// GetConnection returns the websocket client connection.
func (ws *WebSocketHandler) GetConnection() *websocket.Conn {
	return ws.conn
}

// Serve is the method to implement to serve websocket.
func (ws *WebSocketHandler) Serve() {
	if debug {
		log.Println("Serve method not implemented")
	}
}
