package kwiscale

import (
	htpl "html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

var tpldir string

func SetTemplateDir(path string) (string, error) {
	t, err := filepath.Abs(path)

	if err != nil {
		log.Fatalln("Path not ok", err)
	}
	tpldir = t

	log.Println("Template dir set to ", tpldir)

	return t, err
}

type ITemplate interface {
	// Render method to implement to compile and run template
	// then write to Reponse. Should returns number of
	// written byte and error is any.
	Render(io.Writer, string, interface{}) error
}

// Basic template engine that use html/template
type Template struct {
	files []string
}

func (tpl *Template) Render(w io.Writer, file string, ctx interface{}) error {

	file = filepath.Join(tpldir, file)

	tpl.files = make([]string, 0)
	content, err := ioutil.ReadFile(file)

	if err != nil {
		if DEBUG {
			log.Println(err)
		}
		return err
	}

	tpl.parseOverride(content)

	tpl.files = append(tpl.files, file)

	if DEBUG {
		log.Println(tpl.files)
	}

	t, err := htpl.ParseFiles(tpl.files...)
	if err != nil {
		if DEBUG {
			log.Println(err)
		}
		return err
	}

	return t.Execute(w, ctx)
}

func (tpl *Template) parseOverride(content []byte) {

	re := regexp.MustCompile(`\{\{/\*\s*override\s*\"?(.*?)\"?\s*\*/\}\}`)
	matches := re.FindAllSubmatch(content, -1)

	for _, m := range matches {
		// find bottom templates
		tplfile := filepath.Join(tpldir, string(m[1]))
		c, _ := ioutil.ReadFile(tplfile)
		tpl.parseOverride(c)

		tpl.files = append(tpl.files, tplfile)
	}
}
