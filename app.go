package kwiscale

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// handlerManagerRegistry is a map of [name]handlerManager.
var handlerManagerRegistry = make(map[string]*handlerManager)

// handlerRegistry keep the entire handlers - map[name]type.
var handlerRegistry = make(map[string]reflect.Type)

// regexp to find url ordered params from gorilla form.
var urlParamRegexp = regexp.MustCompile(`\{(.+?):`)

// Register takes webhandler and keep type in handlerRegistry.
// It can be called directly (to set handler accessible
// by configuration file), or implicitally by "AddRoute" and "AddNamedRoute()".
func Register(h WebHandler) {
	elem := reflect.ValueOf(h).Elem().Type()
	name := elem.String()
	Log("Registering", name)
	if _, exists := handlerRegistry[name]; !exists {
		handlerRegistry[name] = elem
	}
}

type handlerRouteMap struct {
	handlername string
	route       string
}

// App handles router and handlers.
type App struct {

	// configuration
	Config *Config

	// Global context shared to handlers
	Context map[string]interface{}

	// session store
	sessionstore SessionStore

	// Template engine instance.
	//templateEngine Template

	// The router that will be used
	router *mux.Router

	// List of handler "names" mapped to route (will be create by a factory)
	handlers map[*mux.Route]handlerRouteMap

	// Handler name for error handler.
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
		handlers: make(map[*mux.Route]handlerRouteMap),
		Context:  make(map[string]interface{}),
	}

	// set sessstion store
	a.sessionstore = sessionEngine[config.SessionEngine]
	a.sessionstore.Name(config.SessionName)
	a.sessionstore.SetSecret(config.SessionSecret)
	a.sessionstore.SetOptions(config.SessionEngineOptions)
	a.sessionstore.Init()

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

	app := NewApp(cfg.parse())

	for route, v := range cfg.Routes {
		if handler, ok := handlerRegistry[v.Handler]; ok {
			h := reflect.New(handler).Interface().(WebHandler)
			log.Println(route, h, v.Alias)
			app.addRoute(route, h, v.Alias)
		} else {
			panic("Handler not found: " + v.Handler)
		}
	}

	return app
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
	s := &staticHandler{}
	s.prefix = prefix
	a.AddNamedRoute("/"+prefix+"/{file:.*}", s, "statics")
}

// Implement http.Handler ServeHTTP method.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// try to recover panic if possible, and display
	// an Error page
	defer func() {
		if err := recover(); err != nil {
			app.Error(http.StatusInternalServerError,
				w,
				errors.New("An unexpected error occured"),
				err,
			)
		}
	}()

	var handler WebHandler
	handlerName, route, match := getBestRoute(app, r)

	// if non match
	if _, ok := handlerManagerRegistry[handlerName]; !ok {
		app.Error(http.StatusNotFound, w, ErrNotFound, r.URL)
		return
	}

	// wait for a built handler from registry
	handler = <-handlerManagerRegistry[handlerName].produce()
	Log("Handler found ", handler)
	//assign some vars
	handler.setRoute(route)
	handler.setVars(match.Vars, w, r)
	handler.setApp(app)
	handler.setSessionStore(app.sessionstore)

	// Call Init before starting response
	if code, err := handler.Init(); err != nil {
		Log(err)
		// if returned status is <0, let Init() method do the work
		if code <= 0 {
			Log("Init method returns no error but a status <= 0")
			return
		}
		// Init() stops the request with error and with a status code to use
		app.Error(code, w, err)
		return
	}

	// Nothing stops the process calling Init(), so we can
	// prepare defered destroy
	defer handler.Destroy()

	// Websocket case
	if h, ok := handler.(WSHandler); ok {
		if err := h.upgrade(); err != nil {
			log.Println("Error upgrading Websocket protocol", err)
			return
		}

		defer h.OnClose()
		h.OnConnect()

		switch h.(type) {
		case WSServerHandler:
			serveWS(h)
		case WSJsonHandler:
			serveJSON(h)
		case WSStringHandler:
			serveString(h)
		default:
			app.Error(http.StatusNotImplemented, w, ErrNotImplemented)
		}

		return
	}

	// Standard Request
	if h, ok := handler.(HTTPRequestHandler); ok {
		// RequestHandler case
		w.Header().Add("Connection", "close")
		Log("Respond to RequestHandler", r.Method, h)
		if h == nil {
			app.Error(http.StatusNotFound, w, ErrNotFound, r.Method, h)
			return
		}

		switch r.Method {
		case "GET":
			h.Get()
		case "PUT":
			h.Put()
		case "POST":
			h.Post()
		case "DELETE":
			h.Delete()
		case "HEAD":
			h.Post()
		case "PATCH":
			h.Patch()
		case "OPTIONS":
			h.Options()
		case "TRACE":
			h.Trace()
		default:
			app.Error(http.StatusNotImplemented, w, ErrNotImplemented)
		}
		return
	}

	// if the method have parameters, we can try to call it.
	if app.callMethodWithParameters(r, handler, route, &match) {
		return
	}

	// we should NEVER go to this, but in case of...
	details := "" +
		fmt.Sprintf("Registry: %+v\n", handlerManagerRegistry) +
		fmt.Sprintf("RequestWriter: %+v\n", w) +
		fmt.Sprintf("Reponse: %+v", r) +
		fmt.Sprintf("KwiscaleHandler: %+v\n", handler)
	Log(details)
	app.Error(http.StatusInternalServerError, w, ErrInternalError, details)
}

// Append handler in handlerRegistry and start producing.
// return the name of the handler.
func (app *App) handle(h WebHandler, name string) string {
	handlerType := reflect.ValueOf(h).Elem().Type()
	handlerName := handlerType.String()
	if name == "" {
		name = handlerType.String()
	}
	Log("Register ", name)

	Register(h)
	if _, ok := handlerManagerRegistry[name]; ok {
		// do not create registry manager if it exists
		Log("Registry manager for", name, "already exists")
		return name
	}

	// Append a new handler manager in registry

	hm := &handlerManager{
		handler:  handlerName,
		closer:   make(chan int, 0),
		producer: make(chan WebHandler, app.Config.NbHandlerCache),
	}
	handlerManagerRegistry[name] = hm
	// to be able to fetch handler by real name, only if alias is not given
	if name != handlerName {
		handlerManagerRegistry[handlerName] = hm
	}

	// start to produce handlers
	go handlerManagerRegistry[name].produceHandlers()

	// return the handler name
	return name
}

// AddRoute appends route mapped to handler. Note that rh parameter should
// implement IRequestHandler (generally a struct composing RequestHandler or WebSocketHandler).
func (app *App) AddRoute(route string, handler WebHandler) {
	app.addRoute(route, handler, "")
}

// AddNamedRoute does the same as AddRoute but set the route name instead of
// using the handler name. If the given name already exists or is empty, the method
// panics.
func (app *App) AddNamedRoute(route string, handler WebHandler, name string) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		panic(errors.New("The given name is empty"))
	}

	app.addRoute(route, handler, name)
}

// Add route to the stack.
func (app *App) addRoute(route string, handler WebHandler, routename string) {
	var name string
	handlerType := reflect.ValueOf(handler).Elem().Type()
	if len(routename) == 0 {
		name = handlerType.String()
	} else {
		name = routename
	}

	// record a route
	r := app.router.NewRoute()
	r.Path(route)
	r.Name(name)
	app.handlers[r] = handlerRouteMap{name, route}

	app.handle(handler, name)
}

// SoftStop stops each handler manager goroutine (useful for testing).
func (app *App) SoftStop() chan int {
	c := make(chan int, 0)
	go func() {
		for name, closer := range handlerManagerRegistry {
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

// GetTemplate returns a new instance of Template.
func (a *App) GetTemplate() Template {
	engine := templateEngine[a.Config.TemplateEngine]
	//ttype := reflect.TypeOf(engine)
	t := reflect.New(engine).Interface().(Template)
	t.SetTemplateDir(a.Config.TemplateDir)
	t.SetTemplateOptions(a.Config.TemplateEngineOptions)
	return t
}

// GetRoutes get all routes for a handler
func (a *App) GetRoutes(name string) []*mux.Route {
	routes := make([]*mux.Route, 0)
	for route, _ := range a.handlers {
		if route.GetName() == name {
			routes = append(routes, route)
		}
	}
	return routes
}

// DB returns the App.database configured from Config.
/*func (app *App) DB() DB {
	if app.Config.DB != "" {
		dtype := dbdrivers[app.Config.DB]
		database := reflect.New(dtype).Interface().(DB)
		database.SetOptions(app.Config.DBOptions)
		database.Init()
		return database
	}
	Log("No db selected")
	return nil
}*/

// SetErrorHandler set error handler to replace the default ErrorHandler.
func (app *App) SetErrorHandler(h WebHandler) {
	app.errorHandler = app.handle(h, "")
}

// Error displays an error page with details if any.
func (app *App) Error(status int, w http.ResponseWriter, err error, details ...interface{}) {
	Log(err, details)
	var handler WebHandler
	if app.errorHandler == "" {
		handler = &ErrorHandler{}
	} else {
		handler = <-handlerManagerRegistry[app.errorHandler].produce()
	}
	handler.setApp(app)
	handler.setVars(nil, w, nil)
	handler.(HTTPErrorHandler).setStatus(status)
	handler.(HTTPErrorHandler).setError(err)
	handler.(HTTPErrorHandler).setDetails(details)
	handler.(HTTPRequestHandler).Get()
}

// Try to call method with parameters (if found)
func (app *App) callMethodWithParameters(r *http.Request, handler WebHandler, route *mux.Route, match *mux.RouteMatch) bool {

	h := reflect.ValueOf(handler)
	method := h.MethodByName(strings.Title(strings.ToLower(r.Method)))
	if method.Kind() == reflect.Invalid {
		return false
	}

	// if method is not a parametrized method, return false
	if method.Type().NumIn() == 0 {
		return false
	}

	maps := urlParamRegexp.FindAllStringSubmatch(app.handlers[route].route, -1)
	// build reflect.Value from match.Vars
	args := []reflect.Value{}
	for _, v := range maps {
		args = append(args, reflect.ValueOf(match.Vars[v[1]]))
	}

	// convert parameters, for now we can manage string (default), float, int and bool
	for i := 0; i < method.Type().NumIn(); i++ {
		typ := method.Type().In(i)      // function arg type
		arg := reflect.ValueOf(args[i]) // argument to map

		switch typ.Kind() {
		case reflect.Float32, reflect.Float64:
			v, err := strconv.ParseFloat(args[i].String(), 10)
			if err != nil {
				panic(err)
			}
			arg = reflect.ValueOf(v)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(args[i].String(), 10, 64)
			if err != nil {
				panic(err)
			}
			arg = reflect.ValueOf(v)
		case reflect.Bool:
			switch strings.ToLower(args[i].String()) {
			case "true", "1", "yes", "on":
				arg = reflect.ValueOf(true)
			case "false", "0", "no", "off":
				arg = reflect.ValueOf(false)
			default:
				panic(errors.New(fmt.Sprintf("Boolean URL '%s' value is not reconized", args[i])))
			}
		}

		// convertion
		if typ.Kind() != reflect.String {
			args[i] = arg.Convert(typ)
		}
	}

	method.Call(args)

	return true
}
