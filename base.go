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
	getRequest() *http.Request
	getResponse() http.ResponseWriter
	GetSession(interface{}) interface{}
	SetSession(interface{}, interface{})

	setSessionStore(ISessionStore)
}

type BaseHandler struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Vars         map[string]string
	sessionStore ISessionStore

	app *App
}

// setVars initialize vars from url
func (r *BaseHandler) setVars(v map[string]string, w http.ResponseWriter, req *http.Request) {
	r.Vars = v
	r.Response = w
	r.Request = req
}

// setApp assign App to the handler
func (r *BaseHandler) setApp(a *App) {
	r.app = a
}

// getReponse returns the current response
func (r *BaseHandler) getResponse() http.ResponseWriter {
	return r.Response
}

// getRequest returns the current request
func (b *BaseHandler) getRequest() *http.Request {
	return b.Request
}

// SetSessionStore defines the session store to use
func (b *BaseHandler) setSessionStore(store ISessionStore) {
	b.sessionStore = store
}

func (b *BaseHandler) GetSession(key interface{}) interface{} {
	return b.sessionStore.Get(b, key)
}

func (b *BaseHandler) SetSession(key interface{}, value interface{}) {
	b.sessionStore.Set(b, key, value)
}
