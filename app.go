package kwiscale

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// Config structure that holds configuration
type Config struct {
	TemplateDir    string
	Port           string
	NbHandlerCache int
	TemplateEngine string
	SessionsEngine string
	SessionName    string
	SessionSecret  []byte
}

// App handles router and handlers.
type App struct {

	// configuration
	Config Config

	// session store
	sessionstore ISessionStore

	// Template engine instance.
	templateEngine ITemplate

	// The router that will be used
	router *mux.Router

	// List of handler "names" mapped to route (will be create by a factory)
	handlers map[*mux.Route]string

	// number of handler to keep in a channel
	nbHandlerCache int
}

// Initialize config default values if some are not defined
func initConfig(config *Config) {
	if config.Port == "" {
		config.Port = ":8000"
	}

	if config.NbHandlerCache == 0 {
		config.NbHandlerCache = 5
	}

	if config.TemplateEngine == "" {
		config.TemplateEngine = "basic"
	}

	if config.SessionsEngine == "" {
		config.SessionsEngine = "default"
	}
	if config.SessionName == "" {
		config.SessionName = "kwiscale-session"
	}
	if config.SessionSecret == nil {
		config.SessionSecret = []byte("A very long secret string you should change")
	}
}

// NewApp Create new *App - App constructor.
func NewApp(config Config) *App {

	// fill up config for non-set values
	initConfig(&config)

	// generate app, assigne config, router and handlers map
	a := new(App)
	a.nbHandlerCache = config.NbHandlerCache
	a.router = mux.NewRouter()
	a.handlers = make(map[*mux.Route]string)

	// Get template engine from config
	a.templateEngine = templateEngine[config.TemplateEngine]
	a.templateEngine.SetTemplateDir(config.TemplateDir)

	// set sessstion store
	a.sessionstore = sessionEngine[config.SessionsEngine]
	a.sessionstore.Init()
	a.sessionstore.Name(config.SessionName)
	a.sessionstore.SetSecret(config.SessionSecret)

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
	for route, handler := range app.handlers {
		var match mux.RouteMatch
		if route.Match(r, &match) {
			// construct handler from its name
			if debug {
				log.Printf("Route matches %#v\n", route)
				log.Println("Handler to fetch", handler)
			}

			// wait for a built handler from registry
			req = <-handlerRegistry[handler]

			//assign some vars
			req.(IBaseHandler).setVars(match.Vars, w, r)
			req.(IBaseHandler).setApp(app)
			req.(IBaseHandler).setSessionStore(app.sessionstore)
			break //that's ok, we can continue
		}
	}

	// Websocket case
	if req, ok := req.(IWSHandler); ok {
		req.upgrade()
		req.Serve()
		return
	}

	// RequestHandler case
	if req, ok := req.(IRequestHandler); ok {
		if debug {
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
	if debug {
		log.Printf("RequestWriter: %+v\n", w)
		log.Printf("Reponse: %+v", r)
		log.Printf("KwiscaleHandler: %+v\n", req)
	}
}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler)
func (app *App) AddRoute(route string, rh interface{}) {
	r := app.router.NewRoute()
	r.Path(route)
	r.Name(route)
	app.registerHandler(r, route, rh)
}

// registerHandler add the reflect.Type of handler in registry.
func (app *App) registerHandler(route *mux.Route, name string, h interface{}) {
	handlerType := reflect.TypeOf(h)
	// keep in mind that "route" is an pointer
	app.handlers[route] = handlerType.String()
	// produce handlers
	go app.handlerFactory(handlerType)
}

// handlerFactory continuously generates new handlers in registry.
// It launches a goroutine to produce those handlers. The number of
// handlers to generate in cache is set by Config.NbHandlerCache.
func (app *App) handlerFactory(h reflect.Type) {
	// register factory channel
	// bufferize handler creation in channels
	c := make(chan interface{}, app.nbHandlerCache)
	handlerRegistry[h.String()] = c

	// forever produce handlers
	for {
		if debug {
			log.Println("Append handler in channel: ", reflect.TypeOf(h))
		}
		c <- reflect.New(h).Interface()
	}
}
