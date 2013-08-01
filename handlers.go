package kwiscale

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

// request handler represents an handler
type IRequestHandler interface {
	New() IRequestHandler
	Get()
	Post()
	Head()
	Delete()
	Put()
	Render(template string, context interface{})
	GetSession(name string) interface{}
	SetSession(name string, value interface{})

	setParams(w http.ResponseWriter, r *http.Request, u []string)
	getHandler() IRequestHandler
}

// RequestHandler should be "override" by your handlers. Then your handlers can
// redefine Get(), Put(),... methods (see: IRequestHandler)
//
// Example:
//      type IndexHander struct {
//          kwiscale.RequestHandler `route:"/index"`
//      }
//
//      func (i *IndexHandler) Get(){
//          //...
//      }
type RequestHandler struct {
	IRequestHandler
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
		"_":     t_i18n,
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

// Redirect to given url
func (this *RequestHandler) Redirect(url string) {
	http.Redirect(this.Response, this.Request, url, http.StatusSeeOther)
}

// unset the entire session for the current request
func (this *RequestHandler) EmptySession() {
	delete(sessions, this.SessionId)
}
