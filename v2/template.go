package kwiscale

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

// overrideReg is the regexp that get overrides from templates.
var overrideReg *regexp.Regexp

func init() {
	// this regexp checks if {{ override "template_file" }} is found
	overrideReg = regexp.MustCompile(`\{\{\s*override\s*\"?(.*?)\"?\s*\}\}`)
}

// ITemplateEngine is the base interface to implement a templatse engine.
type ITemplateEngine interface {
	Render(r *Response, tplname string, ctx ContextParams)
}

// BaseTemplateEngine is a basic template engine that uses golang html/template package.
type BaseTemplateEngine struct {
}

// Render writes the template result to the ReponseWriter.
func (b *BaseTemplateEngine) Render(h *Response, tplname string, ctx ContextParams) {
	cur, err := os.Getwd()
	if err != nil {
		h.WriteHeader(http.StatusInternalServerError)
		h.WriteString(err.Error())
	}

	tplcontent, err := parseOverride(cur, h.Server.TemplateDir, tplname)
	if err != nil {
		h.WriteHeader(http.StatusInternalServerError)
		h.WriteString(err.Error())
	}

	log.Println(tplcontent)

	tpl, err := template.New(tplname).Parse(tplcontent)
	if err != nil {
		h.WriteHeader(http.StatusInternalServerError)
		h.WriteString(err.Error())
	}
	tpl.Execute(h.ResponseWriter, ctx)
}

// parseOverride returns full template with override contents.
func parseOverride(dir, templateDir, tplname string) (tplcontent string, err error) {

	tplname = path.Join(dir, templateDir, tplname)
	log.Println("open template:", tplname)
	if btplcontent, err := ioutil.ReadFile(tplname); err != nil {
		return "", err
	} else {
		tplcontent = string(btplcontent)
	}

	// get overrides
	overrides := overrideReg.FindAllStringSubmatch(string(tplcontent), -1)
	for _, ov := range overrides {
		log.Println("overrides:", ov[1])
		override, err := parseOverride(dir, templateDir, ov[1])
		if err != nil {
			return "", err
		}
		tplcontent = strings.Replace(
			tplcontent,
			ov[0],
			override,
			-1)
	}
	return
}
