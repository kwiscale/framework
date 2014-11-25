package kwiscale

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type IRequestHandler interface {
	Get()
	Post()
	Put()
	Head()
	Patch()
	Delete()

	setVars(map[string]string)
}

// this is the main Handler
type App struct {
	Router   *mux.Router
	Handlers map[*mux.Route]IRequestHandler
}

func NewApp() *App {
	a := new(App)
	a.Router = mux.NewRouter()
	a.Handlers = make(map[*mux.Route]IRequestHandler)

	return a
}

type RequestHandler struct {
	app      *App
	Response http.ResponseWriter
	Request  *http.Request
	Vars     map[string]string
}

// Handle IRequestHandeler
func (r *RequestHandler) Get()    {}
func (r *RequestHandler) Put()    {}
func (r *RequestHandler) Post()   {}
func (r *RequestHandler) Delete() {}
func (r *RequestHandler) Head()   {}
func (r *RequestHandler) Patch()  {}

// setVars initialize vars from url
func (r *RequestHandler) setVars(v map[string]string) {
	r.Vars = v
}

// Implement http.Handler ServeHTTP method
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var req IRequestHandler
	for route, handler := range app.Handlers {
		var match mux.RouteMatch
		if route.Match(r, &match) {
			req = handler
			req.setVars(match.Vars)
			break
		}
	}

	if req == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch strings.ToUpper(r.Method) {
	case "GET":
		req.Get()
	case "PUT":
		req.Put()
	case "POST":
		req.Post()
	case "DELETE":
		req.Delete()
	case "HEAD":
		req.Post()
	case "PATCH":
		req.Patch()
	}
}

func (app *App) AddRoute(route string, rh IRequestHandler) {
	r := app.Router.NewRoute()
	r.Path(route)
	r.Name(route)
	app.Handlers[r] = rh
}
