package kwiscale

import "net/http"
import "regexp"
import "strings"
import "html/template"

// request handler represents an handler
type IRequestHandler interface {
	Get()
	Post()
	Head()
	Delete()
	Put()
	Render(template string, context interface{})

	getRoutes() []string
	setParams(w http.ResponseWriter, r *http.Request, u []string)
}

// object that will handler routes and request information
type RequestHandler struct {
	IRequestHandler
	Routes    []string
	Response  http.ResponseWriter
	Request   *http.Request
	UrlParams []string
}

// return routes that handler can answer
func (this *RequestHandler) getRoutes() []string {
	return this.Routes
}

// method that set request and response object
func (this *RequestHandler) setParams(w http.ResponseWriter, r *http.Request, urlparams []string) {
	this.Response = w
	this.Request = r
	this.UrlParams = urlparams
}

// alias to http.Response.Write + type conversion from string to []byte
/*func (this RequestHandler) Write(s string) {
	this.Response.Write([]byte(s))
}*/

// Render a template using override directive if any
func (this *RequestHandler) Render(tpl string, context interface{}) {

	content := getCachedTemplate(tpl)

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

	t, _ := template.New("main").Parse(content)
	t.Execute(this.Response, context)

}
