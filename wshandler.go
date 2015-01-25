package kwiscale

import (
	"log"

	"github.com/gorilla/websocket"
)

type IWSHandler interface {
	Serve()
	upgrade()
}

// WebsockerHandler type
type WSHandler struct {
	BaseHandler
	conn *websocket.Conn
}

func (ws *WSHandler) upgrade() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws.conn, _ = upgrader.Upgrade(ws.Response, ws.Request, nil)
}

func (ws *WSHandler) GetConnection() *websocket.Conn {
	return ws.conn
}

// Serve is the method to implement to serve websocket.
func (ws *WSHandler) Serve() {
	log.Println("Serve method not implemented")
}
