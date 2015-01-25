package kwiscale

import (
	"log"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
)

type Config struct {
	TemplateDir    string
	Port           string
	NbHandlerCache int
}

// App handles router and handlers.
type App struct {
	NbHandlerCache int

	// The router that will be used
	Router *mux.Router

	// List of handler "names" mapped to route (will be create by a factory)
	Handlers map[*mux.Route]string

	// configuration
	Config Config

	// Template engine instance.
	templateEngine ITemplate
}

// setTempate set template engine instance to the current App
func (a *App) setTemplate(t ITemplate) {
	a.templateEngine = t
}

func initConfig(config *Config) {
	if config.Port == "" {
		config.Port = ":8000"
	}

	if config.NbHandlerCache == 0 {
		config.NbHandlerCache = 5
	}

}

// NewApp Create new *App - App constructor.
func NewApp(config Config) *App {
	initConfig(&config)

	a := new(App)
	a.NbHandlerCache = config.NbHandlerCache
	a.Router = mux.NewRouter()
	a.Handlers = make(map[*mux.Route]string)
	SetTemplateDir(config.TemplateDir)
	// use default template engine
	a.setTemplate(new(Template))

	// keep config
	a.Config = config

	return a
}

// ListenAndServe calls http.ListenAndServe method
func (a *App) ListenAndServe() {
	log.Println("Listening", a.Config.Port)
	http.ListenAndServe(a.Config.Port, a)
}

// Implement http.Handler ServeHTTP method
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var req interface{}
	for route, handler := range app.Handlers {
		var match mux.RouteMatch
		if route.Match(r, &match) {
			// construct handler from its name
			if DEBUG {
				log.Printf("Route matches %#v\n", route)
				log.Println("Handler to fetch", handler)
			}

			// wait for a built handler from registry
			req = <-handlerRegistry[handler]
			req.(IBaseHandler).setVars(match.Vars, w, r)
			req.(IBaseHandler).setApp(app)
			break
		}
	}

	if req, ok := req.(IWSHandler); ok {
		req.upgrade()
		req.Serve()
		return
	}

	if req, ok := req.(IRequestHandler); ok {
		if DEBUG {
			log.Println("Respond to IRequestHandler", r.Method, req)
		}
		if req == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		switch r.Method {
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
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	log.Println("The request cannot be served, type of handler is not correct")
	if DEBUG {
		log.Printf("%+v\n", req)
	}

}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler)
func (app *App) AddRoute(route string, rh interface{}) {
	r := app.Router.NewRoute()
	r.Path(route)
	r.Name(route)
	app.registerHandler(r, route, rh)
}

// register type
// handlerFactory will be able to create copies...
func (app *App) registerHandler(route *mux.Route, name string, h interface{}) {
	handlerType := reflect.TypeOf(h)
	// keep in mind that "route" is an pointer
	app.Handlers[route] = handlerType.String()
	go app.handlerFactory(handlerType)
}

// Fill chan c with created handler. No type assertion there.
func (app *App) handlerFactory(h reflect.Type) {
	// register factory channel
	// bufferize handler creation in channels
	c := make(chan interface{}, app.NbHandlerCache)
	handlerRegistry[h.String()] = c

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case *runtime.TypeAssertionError:
				log.Fatal(r)
			default:
				log.Fatal("WTF ???")
			}
		}
	}()

	for {

		if DEBUG {
			log.Println("Append handler in channel: ", reflect.TypeOf(h))
		}
		v := reflect.New(h) // should make a copy
		c <- v.Interface()
	}
}
