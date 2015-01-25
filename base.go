package kwiscale

import "net/http"

// Enable debug logs.
var debug = false

// Change debug mode
func SetDebug(mode bool) {
	debug = mode
}

// Register canals that handle of RequestHandlers.
var handlerRegistry = make(map[string]chan interface{})

type IBaseHandler interface {
	setVars(map[string]string, http.ResponseWriter, *http.Request)
	setApp(*App)
}

type BaseHandler struct {
	Response http.ResponseWriter
	Request  *http.Request
	Vars     map[string]string
	app      *App
}

// setVars initialize vars from url
func (r *BaseHandler) setVars(v map[string]string, w http.ResponseWriter, req *http.Request) {
	r.Vars = v
	r.Response = w
	r.Request = req
}

func (r *BaseHandler) setApp(a *App) {
	r.app = a
}
