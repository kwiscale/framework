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

// RequestHandler that should be composed by users.
type RequestHandler struct {
	BaseHandler
}

// Get implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Get() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Put implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Put() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Post implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Post() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Delete implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Delete() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Head implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Head() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Patch implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Patch() {
	r.Response.WriteHeader(http.StatusNotFound)
}

// Write is an alias to RequestHandler.Request.Write. That implements io.Writer.
func (r *RequestHandler) Write(data []byte) (int, error) {
	return r.Response.Write(data)
}

// WriteString is converts param to []byte then use Write method.
func (r *RequestHandler) WriteString(data string) {
	r.Response.Write([]byte(data))
}

// Stauts write int status to header (use htt.StatusXXX as status).
func (r *RequestHandler) Status(status int) {
	r.Response.WriteHeader(status)
}

// Render calls assigned template engine Render method.
func (r *RequestHandler) Render(file string, ctx interface{}) error {
	return r.app.templateEngine.Render(r, file, ctx)
}
