package kwiscale

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// handlerManager is used to manage handler production.
type handlerManager struct {

	// the handler type to produce
	handler reflect.Type

	// record closers
	closer chan int

	// record handlers (as interface)
	producer chan interface{}
}

// produceHandlers continuously generates new handlers.
// The number of handlers to generate in cache is set by
// Config.NbHandlerCache.
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

// App handles router and handlers.
type App struct {

	// configuration
	Config *Config

	// Global context shared to handlers
	Context map[string]interface{}

	// session store
	sessionstore SessionStore

	// Template engine instance.
	templateEngine Template

	// The router that will be used
	router *mux.Router

	// List of handler "names" mapped to route (will be create by a factory)
	handlers map[*mux.Route]string

	database DB

	errorHandler string
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
		Context:  make(map[string]interface{}),

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

// NewAppFromConfigFile import config file and returns *App.
func NewAppFromConfigFile(filename ...string) *App {
	if len(filename) > 1 {
		panic(errors.New("You should give only one file in NewAppFromConfigFile"))
	}
	file := "config.yml"
	if len(filename) > 0 {
		file = filename[0]
	}
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	cfg := yamlConf{}
	yaml.Unmarshal(content, &cfg)
	return NewApp(cfg.parse())
}

// ListenAndServe calls http.ListenAndServe method
func (a *App) ListenAndServe(port ...string) {
	p := a.Config.Port
	if len(port) > 0 {
		p = port[0]
	}
	log.Println("Listening", p)
	log.Fatal(http.ListenAndServe(p, a))
}

// SetStatic set the route "prefix" to serve files configured in Config.StaticDir
func (a *App) SetStatic(prefix string) {
	path, _ := filepath.Abs(prefix)
	prefix = filepath.Base(path)
	a.AddNamedRoute("/"+prefix+"/{file:.*}", staticHandler{}, "statics")
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
			if code <= 0 {
				Log("Init method returns no error but a status <= 0")
				return
			}
			// Init() stops the request with error and with a status code to use
			// REM HandleError(code, req.getResponse(), req.getRequest(), err)
			app.Error(code, w, err)
			return
		}
		// No stop, so we can
		// prepare defered destroy
		defer req.Destroy()
	} else {
		// Handler is not an handler...
		app.Error(http.StatusNotFound, w, ErrNotFound, r.Method+" "+r.URL.String())
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
			app.Error(http.StatusNotImplemented, w, ErrNotImplemented)
		}
		return
	}

	// Standard Request
	if req, ok := req.(HTTPRequestHandler); ok {
		// RequestHandler case
		w.Header().Add("Connection", "close")
		Log("Respond to RequestHandler", r.Method, req)
		if req == nil {
			app.Error(http.StatusNotFound, w, ErrNotFound, r.Method, req)
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
			app.Error(http.StatusNotImplemented, w, ErrNotImplemented)
		}
		return
	}

	details := "" +
		fmt.Sprintf("Registry: %+v\n", handlerRegistry) +
		fmt.Sprintf("RequestWriter: %+v\n", w) +
		fmt.Sprintf("Reponse: %+v", r) +
		fmt.Sprintf("KwiscaleHandler: %+v\n", req)
	Log(details)
	app.Error(http.StatusInternalServerError, w, ErrInternalError, details)
}

// Append handler in handlerRegistry and start producing.
// return the name of the handler.
func (app *App) handle(h interface{}, name string) string {
	handlerType := reflect.TypeOf(h)
	if name == "" {
		name = handlerType.String()
	}
	Log("Register ", name)

	if _, ok := handlerRegistry[name]; ok {
		// do not create registry manager if it exists
		Log("Registry manager for", name, "already exists")
		return name
	}

	// Append a new handler manager in registry
	handlerRegistry[name] = handlerManager{
		handler:  handlerType,
		closer:   make(chan int, 0),
		producer: make(chan interface{}, app.Config.NbHandlerCache),
	}

	// start to produce handlers
	go handlerRegistry[name].produceHandlers()

	// return the handler name
	return name
}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler or WebSocketHandler).
func (app *App) AddRoute(route string, handler interface{}) {
	app.addRoute(route, handler, "")
}

// AddNamedRoute does the same as AddRoute but set the route name instead of
// using the handler name. If the given name already exists or is empty, the method
// panics.
func (app *App) AddNamedRoute(route string, handler interface{}, name string) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		panic(errors.New("The given name is empty"))
	}

	app.addRoute(route, handler, name)
}

// Add route to the stack.
func (app *App) addRoute(route string, handler interface{}, routename string) {
	var name string
	handlerType := reflect.TypeOf(handler)
	if len(routename) == 0 {
		name = handlerType.String()
	} else {
		name = routename
	}

	// Try to find handler with the given name. If
	// it already exists, there may be a problem !
	//exists := false
	//for n, _ := range app.handlers {
	//	if n.GetName() == name {
	//		exists = true
	//	}
	//}

	// record a route
	//if !exists {
	r := app.router.NewRoute()
	r.Path(route)
	r.Name(name)
	app.handlers[r] = name
	//}

	app.handle(handler, name)
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

// SetErrorHandler set error handler to replace the default ErrorHandler.
func (app *App) SetErrorHandler(h HTTPErrorHandler) {
	app.errorHandler = app.handle(h, "")
}

// Error displays an error page with details if any.
func (app *App) Error(status int, w http.ResponseWriter, err error, details ...interface{}) {
	var handler WebHandler
	if app.errorHandler == "" {
		handler = &ErrorHandler{}
	} else {
		handler = (<-handlerRegistry[app.errorHandler].producer).(WebHandler)
	}
	handler.setApp(app)
	handler.setVars(nil, w, nil)
	handler.(HTTPErrorHandler).setStatus(status)
	handler.(HTTPErrorHandler).setError(err)
	handler.(HTTPErrorHandler).setDetails(details)
	log.Printf("%T, %+v\n", handler, handler)
	handler.(HTTPRequestHandler).Get()
}
