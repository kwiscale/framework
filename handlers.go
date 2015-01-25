package kwiscale

import "net/http"

// IRequestHandler interface which declare HTTP verbs.
type IRequestHandler interface {
	Get()
	Post()
	Put()
	Head()
	Patch()
	Delete()
}

// RequestHandler that should be composed by users
type RequestHandler struct {
	BaseHandler
	app *App
}

// Get implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Get() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Put implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Put() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Post implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Post() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Delete implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Delete() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Head implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Head() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Patch implements IRequestHandler Method - default "not found"
func (r *RequestHandler) Patch() {
	r.Response.WriteHeader(http.StatusNotFound)
}
