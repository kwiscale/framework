package main

import (
	"fmt"
	"net/http"

	"github.com/metal3d/kwiscale/v3"
)

type HomeHandler struct {
	kwiscale.RequestHandler
}

func (h *HomeHandler) Get() {
	println("ok")
}

type OtherHandler struct {
	kwiscale.RequestHandler
}

type WSHandler struct {
	kwiscale.WSHandler
}

func (w *WSHandler) Serve() {
	conn := w.GetConnection()
	defer conn.Close()
	for {
		i, m, err := conn.ReadMessage()
		if err == nil {
			conn.WriteMessage(i, m)
		}
	}
}

func (h *OtherHandler) Get() {
	fmt.Println(h.Vars)
	println("other")
	h.Response.Write([]byte("coucou\n"))
}

func main() {
	kwiscale.DEBUG = true
	app := kwiscale.NewApp()
	app.AddRoute("/ws", WSHandler{})
	app.AddRoute("/toto", HomeHandler{})
	app.AddRoute("/other/{userid:[0-9]+}", OtherHandler{})
	http.ListenAndServe(":8000", app)
}
