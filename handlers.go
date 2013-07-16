package kwiscale

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// thanks Russ Cox - https://groups.google.com/forum/#!topic/golang-nuts/OEdSDgEC7js
// eq reports whether the first argument is equal to
// any of the remaining arguments.
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case string, int, int64, byte, float32, float64:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

// request handler represents an handler
type IRequestHandler interface {
	Get()
	Post()
	Head()
	Delete()
	Put()
	Render(template string, context interface{})
	GetSession(name string) interface{}
	SetSession(name string, value interface{})

	setParams(w http.ResponseWriter, r *http.Request, u []string)
}

// object that will handler routes and request information
type RequestHandler struct {
	IRequestHandler
	Routes    []string
	Response  http.ResponseWriter
	Request   *http.Request
	UrlParams []string
	SessionId string
}

// method that set request and response object
func (this *RequestHandler) setParams(w http.ResponseWriter, r *http.Request, urlparams []string) {
	this.Response = w
	this.Request = r
	this.UrlParams = urlparams
}

// alias to http.Response.Write + type conversion from string to []byte
func (this *RequestHandler) Write(s string) {
	this.Response.Write([]byte(s))
}

// Render a template using override directive if any
func (this *RequestHandler) Render(tpl string, context interface{}) {

	content := getCachedTemplate(tpl)

	fm := template.FuncMap{
		"title": strings.Title,
	}

	re := regexp.MustCompile(`\{\{\s*override\s+"(.*)"\s*\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	for len(matches) > 0 {
		for _, m := range matches {
			// The capture is in m[1]
			// The whole override line is m[0]
			over := getCachedTemplate(m[1])
			content = strings.Replace(content, m[0], string(over), -1)
		}
		matches = re.FindAllStringSubmatch(content, -1)
	}

	t, _ := template.New("main").Funcs(fm).Parse(content)
	t.Execute(this.Response, context)

}

// generate a session uuid
func (this *RequestHandler) GenSessionID() {

	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	this.SessionId = uuid

}

func (this *RequestHandler) CheckSessid() {
	if id, err := this.Request.Cookie("SESSID"); err == nil {
		this.SessionId = id.Value
	} else {
		this.GenSessionID()
	}

	c := http.Cookie{
		Name:  "SESSID",
		Value: this.SessionId,
		Path:  "/",
	}

	http.SetCookie(this.Response, &c)
}

// get sessions
func (this *RequestHandler) GetSession(name string) interface{} {
	this.CheckSessid()
	if sessions == nil {
		sessions = make(map[string]map[string]interface{})
	}

	if s := sessions[this.SessionId]; s != nil {
		if r := s[name]; r != nil {
			return r
		}
	}

	return nil
}

func (this *RequestHandler) SetSession(name string, val interface{}) {

	this.CheckSessid()

	if sessions == nil {
		sessions = make(map[string]map[string]interface{})
	}

	s := sessions[this.SessionId]
	if s == nil {
		sessions[this.SessionId] = make(map[string]interface{})
		s = sessions[this.SessionId]
	}

	s[name] = val
}

// Redirect to given url
func (this *RequestHandler) Redirect(url string) {
	http.Redirect(this.Response, this.Request, url, http.StatusSeeOther)
}

// unset the entire session for the current request
func (this *RequestHandler) EmptySession() {
	delete(sessions, this.SessionId)
}
