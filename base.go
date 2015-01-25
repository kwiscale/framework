package kwiscale

import "net/http"

// Enable debug logs.
var DEBUG = false

func SetDebug(mode bool) {
	DEBUG = mode
}

// Register canals that handle of RequestHandlers.
var handlerRegistry = make(map[string]chan interface{})

type IBaseHandler interface {
	setVars(map[string]string, http.ResponseWriter, *http.Request)
}

type BaseHandler struct {
	Response http.ResponseWriter
	Request  *http.Request
	Vars     map[string]string
}

// setVars initialize vars from url
func (r *BaseHandler) setVars(v map[string]string, w http.ResponseWriter, req *http.Request) {
	r.Vars = v
	r.Response = w
	r.Request = req
}
