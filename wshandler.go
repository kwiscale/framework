package kwiscale

import (
	"log"

	"github.com/gorilla/websocket"
)

var (
	rooms = make(map[string]*wsroom, 0)
)

type wsroom struct {
	conns map[*websocket.Conn]bool
}

func getRoom(path string) *wsroom {
	if r, ok := rooms[path]; !ok {
		r = new(wsroom)
		r.conns = make(map[*websocket.Conn]bool)
		rooms[path] = r
	}

	return rooms[path]
}

func (room *wsroom) add(c *websocket.Conn) {
	room.conns[c] = true
}

func (room *wsroom) remove(c *websocket.Conn) {
	if _, ok := room.conns[c]; ok {
		if debug {
			log.Println("Remove websocket connection", c)
		}
		delete(room.conns, c)
	}
}

// IWSHander is the base template to implement to be able to use
// Websocket
type IWSHandler interface {
	// Serve is the method to implement inside the project
	upgrade() error
	OnConnect() error
	GetConnection() *websocket.Conn
	Close()
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
	if err == nil {

		// record room and append connection
		path := ws.Request.URL.Path
		room := getRoom(path)
		room.add(ws.conn)
	}
	return err
}

// GetConnection returns the websocket client connection.
func (ws *WebSocketHandler) GetConnection() *websocket.Conn {
	return ws.conn
}

func (ws *WebSocketHandler) OnConnect() error {
	return nil
}

func (ws *WebSocketHandler) Write(b []byte) error {
	return ws.conn.WriteMessage(websocket.TextMessage, b)
}

func (ws *WebSocketHandler) WriteString(m string) error {
	return ws.Write([]byte(m))
}

func (ws *WebSocketHandler) WriteJSON(i interface{}) error {
	return ws.conn.WriteJSON(i)
}

func (ws *WebSocketHandler) Close() {
	defer ws.conn.Close()
	if room, ok := rooms[ws.Request.URL.Path]; ok {
		room.remove(ws.conn)
		if len(room.conns) == 0 {
			log.Println("Room", ws.Request.URL.Path, "is empty, deleting")
			delete(rooms, ws.Request.URL.Path)
		}
	}
}

// If WebSocketHandler implements WSServerHandler,
type WSServerHandler interface {
	Serve()
}

// If WebSocketHandler implements WSJsonHandler, framework will
// read socket and call OnJSON each time a json message is received.
type WSJsonHandler interface {
	OnJSON(map[string]interface{}, error)
}

// if WebSocketHandler implements WSStringHandler, framework
// will read socket and call OnMessage() each time a string is received.
type WSStringHandler interface {
	OnMessage(int, string, error)
}

func serveWS(w IWSHandler) {
	defer w.Close()
	w.(WSServerHandler).Serve()
}

func serveJSON(w IWSHandler) {
	c := w.GetConnection()
	defer w.Close()
	for {
		var i map[string]interface{}
		err := c.ReadJSON(i)
		w.(WSJsonHandler).OnJSON(i, err)
	}
}

func serveString(w IWSHandler) {
	c := w.GetConnection()
	defer w.Close()
	for {
		i, p, err := c.ReadMessage()
		w.(WSStringHandler).OnMessage(i, string(p), err)
		if err != nil {
			return
		}
	}

}
