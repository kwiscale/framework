package kwiscale

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var templateEngine = make(map[string]interface{})

// Options to pass to template engines if needed
type TplOptions map[string]interface{}

// RegisterTemplateEngine records template engine that implements ITemplate
// interface. The name is used to let config select the template engine.
func RegisterTemplateEngine(name string, t interface{}) {
	templateEngine[name] = t
}

// ITemplate should be implemented by other template implementation to
// allow RequestHandlers to use Render() method
type Template interface {
	// Render method to implement to compile and run template
	// then write to RequestHandler "w" that is a io.Writer.
	Render(w io.Writer, template string, ctx interface{}) error

	// SetTemplateDir should set the template base directory
	SetTemplateDir(string)

	// SetOptions pass TplOptions to template engine
	SetTemplateOptions(TplOptions)
}

// Basic template engine that use html/template
type BuiltInTemplate struct {
	files  []string
	tpldir string
}

func (tpl *BuiltInTemplate) SetTemplateDir(path string) {
	t, err := filepath.Abs(path)

	if err != nil {
		panic(err)
	}
	tpl.tpldir = t
	Log("Template dir set to ", tpl.tpldir)
}

// Render method for the basic Template system.
// Allow {{/* override "path.html" */}}, no cache, very basic.
func (tpl *BuiltInTemplate) Render(w io.Writer, file string, ctx interface{}) error {
	var err error
	defer func() {
		if err != nil {
			Error(err)
		}
	}()

	file = filepath.Join(tpl.tpldir, file)

	tpl.files = make([]string, 0)
	content, err := ioutil.ReadFile(file)

	// panic if read file breaks
	if err != nil {
		panic(err)
	}

	tpl.parseOverride(content)

	tpl.files = append(tpl.files, file)

	Log(tpl.files)

	funcmap := template.FuncMap{
		"static": func(file string) string {
			app := w.(WebHandler).App()
			url, err := app.GetRoute("statics").URL("file", file)
			if err != nil {
				return err.Error()
			}
			return url.String()
		},
		"url": func(handler string, args ...interface{}) string {
			pairs := make([]string, 0)
			for _, p := range args {
				pairs = append(pairs, fmt.Sprintf("%v", p))
			}
			h := w.(WebHandler).App().GetRoutes(handler)
			base := make([]string, 0)
			for _, r := range h {
				url, err := r.URL(pairs...)
				if err != nil {
					continue
				}
				route := strings.Split(url.String(), "/")
				if len(route) >= len(base) {
					base = route
				}
			}
			if len(base) == 0 {
				return "handler url not realized - please check"
			}
			return strings.Join(base, "/")
		},
	}

	t, err := template.
		New(filepath.Base(file)).
		Funcs(funcmap).
		ParseFiles(tpl.files...)

	// panic if template breaks in parse
	if err != nil {
		panic(err)
	}

	err = t.Execute(w, ctx)
	// return error there !
	return err
}

// Template is basic, it doesn't need options
func (tpl *BuiltInTemplate) SetTemplateOptions(TplOptions) {}

// parseOverride will append overriden templates to be integrating in the
// template list to render
func (tpl *BuiltInTemplate) parseOverride(content []byte) {

	re := regexp.MustCompile(`\{\{/\*\s*override\s*\"?(.*?)\"?\s*\*/\}\}`)
	matches := re.FindAllSubmatch(content, -1)

	for _, m := range matches {
		// find bottom templates
		tplfile := filepath.Join(tpl.tpldir, string(m[1]))
		c, _ := ioutil.ReadFile(tplfile)
		tpl.parseOverride(c)

		tpl.files = append(tpl.files, tplfile)
	}
}
