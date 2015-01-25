package main

import (
	"fmt"

	"github.com/metal3d/kwiscale"
)

type HomeHandler struct {
	kwiscale.RequestHandler
}

func (h *HomeHandler) Get() {
	println("ok")
}

type WSHandler struct {
	kwiscale.WebSocketHandler
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

type OtherHandler struct {
	kwiscale.RequestHandler
}

func (h *OtherHandler) Get() {
	fmt.Println(h.Vars)
	h.Render("content/user.html", h.Vars)
}

func main() {
	kwiscale.SetDebug(true)
	app := kwiscale.NewApp(kwiscale.Config{
		TemplateDir: "template",
		Port:        ":8000",
	})
	app.AddRoute("/", HomeHandler{})
	app.AddRoute("/ws", WSHandler{})
	app.AddRoute("/other/{userid:[0-9]+}", OtherHandler{})
	app.ListenAndServe()
}
