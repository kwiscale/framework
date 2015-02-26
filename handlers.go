package kwiscale

// There, we define HTTP handlers - they are based on BaseHandler

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
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Put implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Put() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Post implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Post() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Delete implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Delete() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Head implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Head() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Patch implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Patch() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
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
