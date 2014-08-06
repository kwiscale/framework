package kwiscale

import (
	"log"
	"net/http"
)

// Type that handles URL params, this is the type to declare as http verb method
// in handlers. Eg. MyHandler.Get => func(m *MyHandler) Get(p ContextParams).
type ContextParams map[interface{}]interface{}

// Response to use to respond to client.
type Response struct {
	ResponseWriter http.ResponseWriter
	templateEngine ITemplateEngine
	Server         *Server
}

// Alias to http.ResponseWriter.WriteHeader.
func (r *Response) WriteHeader(status int) { r.ResponseWriter.WriteHeader(status) }

// Alias to http.ReponseWriter.Write but use string instead of []byte.
func (r *Response) Write(b []byte) { r.ResponseWriter.Write(b) }

// WriteString calls Write() after string to []byte conversion.
func (r *Response) WriteString(m string) { r.Write([]byte(m)) }

// Render tplname with context parameters.
func (r *Response) Render(tplname string, ctx ContextParams) {
	log.Println(r.templateEngine)
	r.templateEngine.Render(r, tplname, ctx)
}
