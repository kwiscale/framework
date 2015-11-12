package kwiscale

import "github.com/gorilla/websocket"

var (
	// keep connection by path.
	rooms = make(map[string]*wsroom, 0)
)

type wsroom struct {
	// connections for the room
	conns map[*WebSocketHandler]bool
}

// Returns the room named "path", or create one if not exists.
func getRoom(path string) *wsroom {
	if r, ok := rooms[path]; !ok {
		r = new(wsroom)
		r.conns = make(map[*WebSocketHandler]bool)
		rooms[path] = r
	}

	return rooms[path]
}

// Add a websocket handler to the room.
func (room *wsroom) add(c *WebSocketHandler) {
	room.conns[c] = true
}

// Remove a websocket handler from the room.
func (room *wsroom) remove(c *WebSocketHandler) {
	if _, ok := room.conns[c]; ok {
		Log("Remove websocket connection", c)
		delete(room.conns, c)
	}
}

// WSHandler is the base interface to implement to be able to use
// Websocket.
type WSHandler interface {
	// Serve is the method to implement inside the project
	upgrade() error
	OnConnect() error
	OnClose() error
	GetConnection() *websocket.Conn
	Close()
}

// WebSocketHandler type to compose a web socket handler.
// To use it, compose a handler with this type and implement one of
// OnJSON(), OnMessage() or Serve() method.
// Example:
//
//	type Example_WebSocketHandler struct{ WebSocketHandler }
//
//	func (m *Example_WebSocketHandler) OnJSON(i interface{}, err error) {
//		if err != nil {
//			m.SendJSON(map[string]string{
//				"error": err.Error(),
//			})
//			return
//		}
//
//		m.SendJSON(map[string]interface{}{
//			"greeting": "Hello !",
//			"data":     i,
//		})
//	}
//
// Previous example send back the message + a greeting message
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
	ws.conn, err = upgrader.Upgrade(ws.response, ws.request, nil)
	if err == nil {

		// record room and append connection
		path := ws.request.URL.Path
		room := getRoom(path)
		room.add(ws)
	}
	return err
}

// GetConnection returns the websocket client connection.
func (ws *WebSocketHandler) GetConnection() *websocket.Conn {
	return ws.conn
}

// OnConnect is called when a client connection is opened.
func (ws *WebSocketHandler) OnConnect() error {
	return nil
}

// OnClose is called when a client connection is closed.
func (ws *WebSocketHandler) OnClose() error {
	return nil
}

func (ws *WebSocketHandler) Write(b []byte) error {
	return ws.conn.WriteMessage(websocket.TextMessage, b)
}

// WriteString is an alias to SendText.
func (ws *WebSocketHandler) WriteString(m string) error {
	return ws.SendText(m)
}

// WriteJSON is an alias for SendJSON.
func (ws *WebSocketHandler) WriteJSON(i interface{}) error {
	return ws.SendJSON(i)
}

// SendJSON send interface "i" in json form to the current client.
func (ws *WebSocketHandler) SendJSON(i interface{}) error {
	return ws.conn.WriteJSON(i)
}

// SendText send string "s" to the current client.
func (ws *WebSocketHandler) SendText(s string) error {
	return ws.conn.WriteMessage(websocket.TextMessage, []byte(s))
}

// SendJSONToThisRoom send interface "i" in json form to the client connected
// to the same room of the current client connection.
func (ws *WebSocketHandler) SendJSONToThisRoom(i interface{}) {
	ws.SendJSONToRoom(ws.request.URL.Path, i)
}

// SendJSONToRoom send the interface "i" in json form to the client connected
// to the the room named "name".
func (ws *WebSocketHandler) SendJSONToRoom(room string, i interface{}) {
	for w := range rooms[room].conns {
		w.SendJSON(i)
	}
}

// SendJSONToAll send the interface "i" in json form to the entire
// client list.
func (ws *WebSocketHandler) SendJSONToAll(i interface{}) {
	for name := range rooms {
		ws.SendJSONToRoom(name, i)
	}
}

// SendTextToThisRoom send message s to the room of the
// current client connection.
func (ws *WebSocketHandler) SendTextToThisRoom(s string) {
	ws.SendTextToRoom(ws.request.URL.Path, s)
}

// SendTextToRoom send message "s" to the room named "name".
func (ws *WebSocketHandler) SendTextToRoom(name, s string) {
	for w := range rooms[name].conns {
		w.SendText(s)
	}
}

// SendTextToAll send message "s" to the entire list of connected clients.
func (ws *WebSocketHandler) SendTextToAll(s string) {
	for name := range rooms {
		ws.SendTextToRoom(name, s)
	}
}

// Close connection after having removed handler from the rooms stack.
func (ws *WebSocketHandler) Close() {
	defer ws.conn.Close()
	if room, ok := rooms[ws.request.URL.Path]; ok {
		room.remove(ws)
		if len(room.conns) == 0 {
			delete(rooms, ws.request.URL.Path)
		}
	}
}

// WSServerHandler interface to serve continuously.
type WSServerHandler interface {
	Serve()
}

// WSJsonHandler interface, framework will
// read socket and call OnJSON each time a json message is received.
type WSJsonHandler interface {
	OnJSON(interface{}, error)
}

// WSStringHandler interface, framework
// will read socket and call OnMessage() each time a string is received.
type WSStringHandler interface {
	OnMessage(int, string, error)
}

func serveWS(w WSHandler) {
	defer w.Close()
	w.(WSServerHandler).Serve()
}

// Serve JSON.
func serveJSON(w WSHandler) {
	c := w.GetConnection()
	defer w.Close()
	for {
		var i interface{}
		err := c.ReadJSON(&i)
		w.(WSJsonHandler).OnJSON(i, err)
		if err != nil {
			return
		}
	}
}

// Serve string messages
func serveString(w WSHandler) {
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
