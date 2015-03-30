package kwiscale

import (
	"log"
	"net/http"
	"reflect"

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

type HandlerFactory func() IBaseHandler

// handlerFactory continuously generates new handlers in registry.
// It launches a goroutine to produce those handlers. The number of
// handlers to generate in cache is set by Config.NbHandlerCache.
// Return a chanel to write in to close handler production
func (manager handlerManager) produceHandlers() {
	// forever produce handlers until closer is called
	for {
		select {
		case manager.producer <- reflect.New(manager.handler).Interface():
			if debug {
				log.Println("Appended handler ", manager.handler.Name())
			}
		case <-manager.closer:
			// Someone closed the factory
			return
		}
	}
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
	SessionsEngine string
	// SessionName is the name of session, eg. Cookie name, default is "kwiscale-session"
	SessionName string
	// A secret string to encrypt cookie
	SessionSecret []byte
	// Static directory (to put css, images, and so on...)
	StaticDir string
	// Activate static in memory cache
	StaticCacheEnabled bool

	// DBDriver should be the name of a
	// registered DB Driver (sqlite3, postgresql, mysql/mariadb...)
	//DBDriver string

	// DBURL is the connection path/url to the database
	//DBURL string
}

// App handles router and handlers.
type App struct {

	// configuration
	Config *Config

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

	// DB connection
	//DB IORM
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

	if config.SessionsEngine == "" {
		config.SessionsEngine = "default"
	}
	if config.SessionName == "" {
		config.SessionName = "kwiscale-session"
	}
	if config.SessionSecret == nil {
		config.SessionSecret = []byte("A very long secret string you should change")
	}

	/* no default ! use may not want to have a static handler
	//
	if config.StaticDir == "" {
		config.StaticDir = "static"
	}
	/**/
	return config
}

// NewApp Create new *App - App constructor.
func NewApp(config *Config) *App {

	// fill up config for non-set values
	config = initConfig(config)

	if debug {
		log.Printf("%+v\n", config)
	}

	// generate app, assign config, router and handlers map
	a := &App{
		nbHandlerCache: config.NbHandlerCache,
		router:         mux.NewRouter(),
		handlers:       make(map[*mux.Route]string),

		// Get template engine from config
		templateEngine: templateEngine[config.TemplateEngine],
	}

	a.templateEngine.SetTemplateDir(config.TemplateDir)
	a.templateEngine.SetTemplateOptions(&config.TemplateEngineOptions)

	// set sessstion store
	a.sessionstore = sessionEngine[config.SessionsEngine]
	a.sessionstore.Name(config.SessionName)
	a.sessionstore.SetSecret(config.SessionSecret)
	a.sessionstore.Init()

	if config.StaticDir != "" {
		a.SetStatic(config.StaticDir)
	}

	/*if config.DBDriver != "" {
		db := ormDriverRegistry[config.DBDriver]
		if db == nil {
			panic(errors.New("Unable to find driver " + config.DBDriver))
		}
		db.ConnectionString(config.DBURL)
		a.DB = db
	}*/

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
	a.AddRoute("/"+prefix+"/{file:.*}", staticHandler{})
}

// Implement http.Handler ServeHTTP method.
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
			req = <-handlerRegistry[handler].producer

			if debug {
				log.Print("Handler found ", req)
			}

			//assign some vars
			req.(IBaseHandler).setVars(match.Vars, w, r)
			req.(IBaseHandler).setApp(app)
			req.(IBaseHandler).setSessionStore(app.sessionstore)
			break //that's ok, we can continue
		}
		// code hasn't breaked, so we didn't found handler
	}

	if req, ok := req.(IBaseHandler); ok {
		// Call Init before starting response
		if err, code := req.Init(); err != nil {
			// Init stops the request with error
			HandleError(code, req.getResponse(), req.getRequest(), err)
			return
		}
		// Prepare defered destroy
		defer req.Destroy()
	} else {
		HandleError(http.StatusNotFound, w, r, nil)
		return
	}

	// Websocket case
	if req, ok := req.(IWSHandler); ok {
		req.upgrade()
		req.Serve()
		return
	}
	if req, ok := req.(IRequestHandler); ok {
		// RequestHandler case
		w.Header().Add("Connection", "close")
		if debug {
			log.Println("Respond to IRequestHandler", r.Method, req)
		}
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
	} else {
		HandleError(http.StatusInternalServerError, w, r, nil)
		if debug {
			log.Printf("RequestWriter: %+v\n", w)
			log.Printf("Reponse: %+v", r)
			log.Printf("KwiscaleHandler: %+v\n", req)
		}
	}
}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler).
func (app *App) AddRoute(route string, handler interface{}) {
	r := app.router.NewRoute()
	r.Path(route)
	r.Name(route)

	handlerType := reflect.TypeOf(handler)
	// keep in mind that "route" is an pointer
	app.handlers[r] = handlerType.String()
	if debug {
		log.Print("Register ", handlerType.String())
	}

	// register factory channel
	manager := handlerManager{
		handler:  handlerType,
		closer:   make(chan int, 0),
		producer: make(chan interface{}, app.nbHandlerCache),
	}
	// produce handlers
	handlerRegistry[handlerType.String()] = manager
	go manager.produceHandlers()
}

// HangOut stops each handler manager goroutine (useful for testing).
func (app *App) SoftStop() {
	for name, closer := range handlerRegistry {
		if debug {
			log.Println("Closing ", name)
		}
		closer.closer <- 1
	}
}
