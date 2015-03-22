package kwiscale

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

// Enable debug logs.
var debug = false

// Change debug mode
func SetDebug(mode bool) {
	debug = mode
}

// IBaseHandler is the main handler interface that every handler sould
// implement.
type IBaseHandler interface {
	setVars(map[string]string, http.ResponseWriter, *http.Request)
	setApp(*App)
	getRequest() *http.Request
	getResponse() http.ResponseWriter
	GetSession(interface{}) (interface{}, error)
	SetSession(interface{}, interface{})

	setSessionStore(ISessionStore)
}

// BaseHandler is the parent struct of every Handler.
type BaseHandler struct {
	Response     http.ResponseWriter
	Request      *http.Request
	Vars         map[string]string
	sessionStore ISessionStore

	app *App
}

// setVars initialize vars from url
func (r *BaseHandler) setVars(v map[string]string, w http.ResponseWriter, req *http.Request) {
	r.Vars = v
	r.Response = w
	r.Request = req
}

// setApp assign App to the handler
func (r *BaseHandler) setApp(a *App) {
	r.app = a
}

// getReponse returns the current response.
func (r *BaseHandler) getResponse() http.ResponseWriter {
	return r.Response
}

// getRequest returns the current request.
func (b *BaseHandler) getRequest() *http.Request {
	return b.Request
}

// SetSessionStore defines the session store to use.
func (b *BaseHandler) setSessionStore(store ISessionStore) {
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

// GetPayload returns the Body content.
func (b *BaseHandler) GetPayload() []byte {
	content, err := ioutil.ReadAll(b.Request.Body)
	if err != nil {
		return nil
	}
	return content
}

// GetJSONPayload unmarshal body to the "v" interface.
func (b *BaseHandler) GetJSONPayload(v interface{}) error {
	return json.Unmarshal(b.GetPayload(), v)
}

// GetPos return the post data for the given "name" argument.
func (b *BaseHandler) GetPost(name string) string {
	return b.Request.PostFormValue(name)
}

// GetPostValues returns the entire posted values.
func (b *BaseHandler) GetPostValues() url.Values {
	b.Request.ParseForm()
	return b.Request.PostForm
}

// GetPostFile returns the "name" file pointer and information from the post data.
func (b *BaseHandler) GetPostFile(name string) (multipart.File, *multipart.FileHeader, error) {
	b.Request.ParseForm()
	return b.Request.FormFile(name)
}

// SavePostFile save the given "name" file to the "to" path.
func (b *BaseHandler) SavePostFile(name, to string) error {
	b.Request.ParseForm()
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
