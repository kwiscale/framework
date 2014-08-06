package kwiscale

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Factory is a type to function that is used to serve Route.
// It should return *Handler and the return type must be IHandler.
type Factory func() IHandler

// IHandler is use to interface HTTP verb. This is the type that should return each Factory
// declared in project.
type IHandler interface {
	Get(map[string]string)
	Post(map[string]string)
	Head(map[string]string)
	Delete(map[string]string)
	Patch(map[string]string)
	Put(map[string]string)
	bind(http.ResponseWriter, *http.Request)
	getResponse() *Response
	getRequest() *http.Request
	setTemplateEngine(ITemplateEngine)
	getTemplateEngine() ITemplateEngine
	setServer(*Server)
	getServer() *Server
}

// Handler is the base class that should be composed by projects
// Each route should return a Factory that return *Handler.
type Handler struct {
	Response Response
	Request  *http.Request
	Server   *Server
}

func (h *Handler) getResponse() *Response {
	return &h.Response
}

func (h *Handler) getRequest() *http.Request {
	return h.Request
}

func (h *Handler) getTemplateEngine() ITemplateEngine {
	return h.Response.templateEngine
}

func (h *Handler) setTemplateEngine(t ITemplateEngine) {
	h.Response.templateEngine = t
}

func (h *Handler) setServer(s *Server) {
	h.Server = s
}

func (h *Handler) getServer() *Server {
	return h.Server
}

/* Default response is 501 Unimplemented*/

// Default GET response 501 Unimplemented.
func (h *Handler) Get(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// Default POST response 501 Unimplemented.
func (h *Handler) Post(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// Default HEAD response 501 Unimplemented.
func (h *Handler) Head(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// Default DELETE response 501 Unimplemented.
func (h *Handler) Delete(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// Default PATCH response 501 Unimplemented.
func (h *Handler) Patch(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// Default PUT response 501 Unimplemented.
func (h *Handler) Put(map[string]string) { h.Response.WriteHeader(http.StatusNotImplemented) }

// bind response and request to the handler.
func (h *Handler) bind(w http.ResponseWriter, r *http.Request) {
	h.Response = Response{w, nil, h.getServer()}
	h.Request = r
}

// make the handler respond to the correect HTTP verb.
func dispatch(h IHandler, w http.ResponseWriter, r *http.Request) {
	h.bind(w, r)

	if h.getTemplateEngine() == nil {
		h.setTemplateEngine(&BaseTemplateEngine{})
	}

	switch r.Method {
	case "GET":
		h.Get(mux.Vars(r))
	case "POST":
		h.Post(mux.Vars(r))
	case "DELETE":
		h.Delete(mux.Vars(r))
	case "HEAD":
		h.Head(mux.Vars(r))
	case "PATCH":
		h.Patch(mux.Vars(r))
	case "PUT":
		h.Put(mux.Vars(r))
	default:
		h.getResponse().WriteHeader(http.StatusBadRequest)
	}
}
