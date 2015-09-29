package kwiscale

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
)

// handlerManager is used to manage handler production and close
type handlerManager struct {

	// the handler type to produce
	handler reflect.Type

	// record closers
	closer chan int

	// record handlers (as interface)
	producer chan interface{}
}

// produceHandlers continuously generates new handlers in registry.
// It launches a goroutine to produce those handlers. The number of
// handlers to generate in cache is set by Config.NbHandlerCache.
// Return a chanel to write in to close handler production
func (manager handlerManager) produceHandlers() {
	// forever produce handlers until closer is called
	for {
		select {
		case manager.producer <- reflect.New(manager.handler).Interface():
			Log("Appended handler ", manager.handler.Name())
		case <-manager.closer:
			// Someone closed the factory
			break
		}
	}
	Log("Quitting ", manager.handler.Name, "producer")
}

// the full registry
var handlerRegistry = make(map[string]handlerManager)

// Config structure that holds configuration
type Config struct {
	// Root directory where TemplateEngine will get files
	TemplateDir string
	// Port to listen
	Port string
	// Number of handler to prepare
	NbHandlerCache int
	// TemplateEngine to use (default, pango2...)
	TemplateEngine string
	// Template engine options (some addons need options)
	TemplateEngineOptions TplOptions

	// SessionEngine (default is a file storage)
	SessionEngine string
	// SessionName is the name of session, eg. Cookie name, default is "kwiscale-session"
	SessionName string
	// A secret string to encrypt cookie
	SessionSecret []byte
	// Configuration for SessionEngine
	SessionEngineOptions SessionEngineOptions

	// Static directory (to put css, images, and so on...)
	StaticDir string
	// Activate static in memory cache
	StaticCacheEnabled bool

	// StrictSlash allows to match route that have trailing slashes
	StrictSlash bool

	// Datastrore
	DB        string
	DBOptions DBOptions
}

// App handles router and handlers.
type App struct {

	// configuration
	Config *Config

	// session store
	sessionstore SessionStore

	// Template engine instance.
	templateEngine Template

	// The router that will be used
	router *mux.Router

	// List of handler "names" mapped to route (will be create by a factory)
	handlers map[*mux.Route]string

	database DB
}

// Initialize config default values if some are not defined
func initConfig(config *Config) *Config {
	if config == nil {
		config = new(Config)
	}

	if config.Port == "" {
		config.Port = ":8000"
	}

	if config.NbHandlerCache == 0 {
		config.NbHandlerCache = 5
	}

	if config.TemplateEngine == "" {
		config.TemplateEngine = "basic"
	}

	if config.TemplateEngineOptions == nil {
		config.TemplateEngineOptions = make(TplOptions)
	}

	if config.SessionEngine == "" {
		config.SessionEngine = "default"
	}
	if config.SessionName == "" {
		config.SessionName = "kwiscale-session"
	}
	if config.SessionSecret == nil {
		config.SessionSecret = []byte("A very long secret string you should change")
	}
	if config.SessionEngineOptions == nil {
		config.SessionEngineOptions = make(SessionEngineOptions)
	}

	return config
}

// NewApp Create new *App - App constructor.
func NewApp(config *Config) *App {

	// fill up config for non-set values
	config = initConfig(config)

	Log(fmt.Sprintf("%+v\n", config))

	// generate app, assign config, router and handlers map
	a := &App{
		Config:   config,
		router:   mux.NewRouter(),
		handlers: make(map[*mux.Route]string),

		// Get template engine from config
		templateEngine: templateEngine[config.TemplateEngine],
	}

	a.templateEngine.SetTemplateDir(config.TemplateDir)
	a.templateEngine.SetTemplateOptions(config.TemplateEngineOptions)

	// set sessstion store
	a.sessionstore = sessionEngine[config.SessionEngine]
	a.sessionstore.Name(config.SessionName)
	a.sessionstore.SetSecret(config.SessionSecret)
	a.sessionstore.SetOptions(config.SessionEngineOptions)
	a.sessionstore.Init()

	// set Datastore
	if config.DB != "" {
		a.database = dbdrivers[config.DB]
		a.database.SetOptions(config.DBOptions)
		a.database.Init()

	}

	if config.StaticDir != "" {
		a.SetStatic(config.StaticDir)
	}

	a.router.StrictSlash(config.StrictSlash)

	// keep config
	a.Config = config

	return a
}

// ListenAndServe calls http.ListenAndServe method
func (a *App) ListenAndServe(port ...string) {
	p := a.Config.Port
	if len(port) > 0 {
		p = port[0]
	}
	log.Println("Listening", p)
	http.ListenAndServe(p, a)
}

// SetStatic set the route "prefix" to serve files configured in Config.StaticDir
func (a *App) SetStatic(prefix string) {
	path, _ := filepath.Abs(prefix)
	prefix = filepath.Base(path)
	a.AddRoute("/"+prefix+"/{file:.*}", staticHandler{})
}

// Implement http.Handler ServeHTTP method.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var req interface{}
	handler, route, match := getBestRoute(app, r)

	if handler != "" {
		// wait for a built handler from registry
		req = <-handlerRegistry[handler].producer

		Log("Handler found ", req)

		//assign some vars
		req.(WebHandler).setRoute(route)
		req.(WebHandler).setVars(match.Vars, w, r)
		req.(WebHandler).setApp(app)
		req.(WebHandler).setSessionStore(app.sessionstore)
	}

	if req, ok := req.(WebHandler); ok {
		// Call Init before starting response
		if code, err := req.Init(); err != nil {
			Log(err)
			// if returned status is <0, let Init() method do the work
			if code < 0 {
				Log("Init method returns no error but a status < 0")
				return
			}
			// Else
			// Init() stops the request with error with a status code to use
			HandleError(code, req.getResponse(), req.getRequest(), err)
			return
		}
		// No stop, so we can
		// prepare defered destroy
		defer req.Destroy()
	} else {
		// Handler is not an handler...
		HandleError(http.StatusNotFound, w, r, nil)
		return
	}

	// Websocket case
	if req, ok := req.(WSHandler); ok {
		if err := req.upgrade(); err != nil {
			log.Println("Error upgrading Websocket protocol", err)
			return
		}

		defer req.OnClose()
		req.OnConnect()

		if _, ok := req.(WSServerHandler); ok {
			serveWS(req)
		} else if _, ok := req.(WSJsonHandler); ok {
			serveJSON(req)
		} else if _, ok := req.(WSStringHandler); ok {
			serveString(req)
		} else {
			log.Println("ws handler has not implemented one of the method: OnJSON, OnString or Serve")
			HandleError(http.StatusNotImplemented, w, r, nil)
		}
		return
	}

	// Standard Request
	if req, ok := req.(HTTPRequestHandler); ok {
		// RequestHandler case
		w.Header().Add("Connection", "close")
		Log("Respond to IRequestHandler", r.Method, req)
		if req == nil {
			HandleError(http.StatusNotFound, w, r, nil)
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
		case "OPTIONS":
			req.Options()
		case "TRACE":
			req.Trace()
		default:
			HandleError(http.StatusNotImplemented, w, r, nil)
		}
		return
	}

	HandleError(http.StatusInternalServerError, w, r, nil)
	Log(fmt.Sprintf("Registry: %+v\n", handlerRegistry))
	Log(fmt.Sprintf("RequestWriter: %+v\n", w))
	Log(fmt.Sprintf("Reponse: %+v", r))
	Log(fmt.Sprintf("KwiscaleHandler: %+v\n", req))
}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler or WebSocketHandler).
func (app *App) AddRoute(route string, handler interface{}) {
	app.addRoute(route, handler)
}

// AddNamedRoute does the same as AddRoute but set the route name instead of
// using the handler name. If the given name already exists or is empty, the method
// panics.
func (app *App) AddNamedRoute(route string, handler interface{}, name string) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		panic(errors.New("The given name is empty"))
	}

	for n, _ := range app.handlers {
		if n.GetName() == name {
			panic(errors.New("The given name already exists:" + name))
		}
	}

	app.addRoute(route, handler, name)
}

// Add route to the stack
func (app *App) addRoute(route string, handler interface{}, routename ...string) {
	var name string
	handlerType := reflect.TypeOf(handler)
	if len(routename) == 0 {
		name = handlerType.String()
	} else {
		name = routename[0]
	}

	// record a route
	r := app.router.NewRoute()
	r.Path(route)
	r.Name(name)

	app.handlers[r] = name
	Log("Register ", name)

	if _, ok := handlerRegistry[name]; ok {
		// do not create registry manager if it exists
		Log("Registry manager for", name, "already exists")
		return
	}
	// register factory channel
	manager := handlerManager{
		handler:  handlerType,
		closer:   make(chan int, 0),
		producer: make(chan interface{}, app.Config.NbHandlerCache),
	}
	// produce handlers
	handlerRegistry[name] = manager
	go manager.produceHandlers()
}

// SoftStop stops each handler manager goroutine (useful for testing).
func (app *App) SoftStop() chan int {
	c := make(chan int, 0)
	go func() {
		for name, closer := range handlerRegistry {
			Log("Closing ", name)
			closer.closer <- 1
			Log("Closed ", name)
		}
		c <- 1
	}()
	return c
}

// GetRoute return the *mux.Route that have the given name.
func (a *App) GetRoute(name string) *mux.Route {

	for route, _ := range a.handlers {
		if route.GetName() == name {
			return route
		}
	}
	return nil
}

// DB returns the App.database configured from Config.
func (app *App) DB() DB {
	return app.database
}
