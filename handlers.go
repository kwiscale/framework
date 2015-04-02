package kwiscale

// There, we define HTTP handlers - they are based on BaseHandler

import (
	"encoding/json"
	"net/http"
)

// IRequestHandler interface which declare HTTP verbs.
type IRequestHandler interface {
	Get()
	Post()
	Put()
	Head()
	Patch()
	Delete()
	Options()
	Trace()
	Redirect(url string)
	RedirectWithStatus(url string, httpStatus int)
	GlobalCtx() map[string]interface{}
}

// RequestHandler that should be composed by users.
type RequestHandler struct {
	BaseHandler
	tplGlobals map[string]interface{}
}

// GlobalCtx Returns global template context.
func (r *RequestHandler) GlobalCtx() map[string]interface{} {
	if r.tplGlobals == nil {
		r.tplGlobals = make(map[string]interface{})
	}
	return r.tplGlobals
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

// Options implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Options() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Trace implements IRequestHandler Method - default "not found".
func (r *RequestHandler) Trace() {
	HandleError(http.StatusNotFound, r.getResponse(), r.getRequest(), nil)
}

// Write is an alias to RequestHandler.Request.Write. That implements io.Writer.
func (r *RequestHandler) Write(data []byte) (int, error) {
	return r.Response.Write(data)
}

// WriteString is converts param to []byte then use Write method.
func (r *RequestHandler) WriteString(data string) (int, error) {
	return r.Write([]byte(data))
}

// WriteJSON converts data to json then send bytes.
// This methods set content-type to application/json (RFC 4627)
func (r *RequestHandler) WriteJSON(data interface{}) (int, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return -1, err
	}
	r.Response.Header().Add("Content-Type", "application/json")
	return r.Write(b)
}

// Stauts write int status to header (use htt.StatusXXX as status).
func (r *RequestHandler) Status(status int) {
	r.Response.WriteHeader(status)
}

// Render calls assigned template engine Render method.
// This method copies globalCtx and write ctx inside. So, contexts are not overriden, it
// only merge 2 context in a new one that is passed to template.
func (r *RequestHandler) Render(file string, ctx map[string]interface{}) error {
	// merge global context with the given
	// ctx should override gobal context
	newctx := make(map[string]interface{})
	for k, v := range r.GlobalCtx() {
		newctx[k] = v
	}
	for k, v := range ctx {
		newctx[k] = v
	}
	return r.app.templateEngine.Render(r, file, newctx)
}

// redirect client with http status.
func (r *RequestHandler) redirect(uri string, status int) {
	r.Response.Header().Add("Location", uri)
	if status < 0 {
		// by default, we use SeeOther status. This status
		// should change HTTP verb to GET
		status = http.StatusSeeOther
	}
	r.Status(status)
}

// Redirect will redirect client to uri using http.StatusSeeOther.
func (r *RequestHandler) Redirect(uri string) {
	r.redirect(uri, -1)
}

// RedirectWithStatus will redirect client to uri using given status.
func (r *RequestHandler) RedirectWithStatus(uri string, status int) {
	r.redirect(uri, status)
}
