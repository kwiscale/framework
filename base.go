package kwiscale

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

type PayloadType int

const (
	JSON PayloadType = iota
	BYTES
	STRING
)

// Enable debug logs.
var debug = false

// Change debug mode
func SetDebug(mode bool) {
	debug = mode
}

// WebHandler is the main handler interface that every handler sould
// implement.
type WebHandler interface {
	setVars(map[string]string, http.ResponseWriter, *http.Request)
	setApp(*App)
	App() *App
	setRoute(*mux.Route)

	getRequest() *http.Request
	getResponse() http.ResponseWriter

	Request() *http.Request
	Response() http.ResponseWriter

	GetSession(interface{}) (interface{}, error)
	SetSession(interface{}, interface{})

	setSessionStore(SessionStore)
	Init() (status int, message error)
	Destroy()
	URL(...string) (*url.URL, error)
	// GetApp() *App                       // deprecated
	// GetURL(...string) (*url.URL, error) // deprecated
	// GetRequest() *http.Request          // deprecated
	// GetResponse() http.ResponseWriter   // deprecated
}

// BaseHandler is the parent struct of every Handler.
// Implement WebHandler.
type BaseHandler struct {
	response     http.ResponseWriter
	request      *http.Request
	Vars         map[string]string
	sessionStore SessionStore

	route     *mux.Route
	app       *App
	routepath string
}

// Init is called before the begin of response (before Get, Post, and so on).
// If error is not nil, framework will write response with the second argument as http status.
func (r *BaseHandler) Init() (int, error) {
	return -1, nil
}

// Destroy is called as defered function after response.
func (r *BaseHandler) Destroy() {}

// setVars initialize vars from url
func (r *BaseHandler) setVars(v map[string]string, w http.ResponseWriter, req *http.Request) {
	r.Vars = v
	r.response = w
	r.request = req
}

// setApp assign App to the handler
func (r *BaseHandler) setApp(a *App) {
	r.app = a
}

// setRoute register mux.Route in the handler.
func (b *BaseHandler) setRoute(r *mux.Route) {
	b.route = r
}

// getReponse returns the current response.
func (r *BaseHandler) getResponse() http.ResponseWriter {
	return r.response
}

// getRequest returns the current request.
func (b *BaseHandler) getRequest() *http.Request {
	return b.request
}

// Reponse returns the current response.
func (r *BaseHandler) Response() http.ResponseWriter {
	return r.getResponse()
}

// Request returns the current request.
func (b *BaseHandler) Request() *http.Request {
	return b.getRequest()
}

// SetSessionStore defines the session store to use.
func (b *BaseHandler) setSessionStore(store SessionStore) {
	b.sessionStore = store
}

// GetSession return the session value of "key".
func (b *BaseHandler) GetSession(key interface{}) (interface{}, error) {
	return b.sessionStore.Get(b, key)
}

// SetSession set the "key" session to "value".
func (b *BaseHandler) SetSession(key interface{}, value interface{}) {
	b.sessionStore.Set(b, key, value)
}

// CleanSession remove every key/value of the current session.
func (b *BaseHandler) CleanSession() {
	b.sessionStore.Clean(b)
}

// Payload returns the Body content.
func (b *BaseHandler) Payload() []byte {
	content, err := ioutil.ReadAll(b.request.Body)
	if err != nil {
		return nil
	}
	return content
}

// GetJSONPayload unmarshal body to the "v" interface.
func (b *BaseHandler) JSONPayload(v interface{}) error {
	return json.Unmarshal(b.Payload(), v)
}

// PostValue returns the post data for the given "name" argument.
// If POST value is empty, return "def" instead. If no "def" is provided, return an empty string by default.
func (b *BaseHandler) PostValue(name string, def ...string) string {
	if res := b.request.PostFormValue(name); res != "" {
		return res
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// PostValues returns the entire posted values.
func (b *BaseHandler) PostValues() url.Values {
	b.request.ParseForm()
	return b.request.PostForm
}

// GetPostFile returns the "name" file pointer and information from the post data.
func (b *BaseHandler) GetPostFile(name string) (multipart.File, *multipart.FileHeader, error) {
	b.request.ParseForm()
	return b.request.FormFile(name)
}

// SavePostFile save the given "name" file to the "to" path.
func (b *BaseHandler) SavePostFile(name, to string) error {
	b.request.ParseForm()
	file, _, err := b.GetPostFile(name)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := os.Create(to)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)

	return err
}

// URL return an url based on the declared route and given string pair.
func (b *BaseHandler) URL(s ...string) (*url.URL, error) {
	return b.route.URL(s...)
}

// App() returns the current application.
func (b *BaseHandler) App() *App {
	return b.app
}
