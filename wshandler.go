package kwiscale

import (
	"log"

	"github.com/gorilla/websocket"
)

type IWSHandler interface {
	// Serve is the method to implement inside the project
	Serve()
	upgrade()
}

// WebsockerHandler type
type WebSocketHandler struct {
	BaseHandler
	conn *websocket.Conn
}

// upgrade protocol to use websocket communication
func (ws *WebSocketHandler) upgrade() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws.conn, _ = upgrader.Upgrade(ws.Response, ws.Request, nil)
}

// GetConnection returns the websocket client connection.
func (ws *WebSocketHandler) GetConnection() *websocket.Conn {
	return ws.conn
}

// Serve is the method to implement to serve websocket.
func (ws *WebSocketHandler) Serve() {
	if DEBUG {
		log.Println("Serve method not implemented")
	}
}
