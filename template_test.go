package kwiscale

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// A basic request handler
type templateHandler struct{ RequestHandler }

// Respond to GET
func (h *templateHandler) Get() {
	h.Render("main.html", map[string]interface{}{
		"Foo": "bar",
	})
}

// Test a template rendering with a not found file.
func TestRenderError(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://www.test.com/", nil)
	w := httptest.NewRecorder()
	app := NewApp(&Config{
		TemplateDir: "./",
	})
	app.AddRoute("/", &templateHandler{})
	app.ServeHTTP(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatal("A non existing template should do a 500 error, but it returns ", w.Code)
	}

}

// Test a template rendering with found template.
func TestRenderVar(t *testing.T) {

	tpl := `<p>{{ .Foo }}</p>`
	d, err := ioutil.TempDir("", "kwiscale-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(d)

	err = ioutil.WriteFile(filepath.Join(d, "main.html"), []byte(tpl), 0644)
	if err != nil {
		t.Fatal(err)
	}

	r, _ := http.NewRequest("GET", "http://www.test.com/", nil)
	w := httptest.NewRecorder()
	app := NewApp(&Config{
		TemplateDir: d,
	})
	app.AddRoute("/", &templateHandler{})
	app.ServeHTTP(w, r)

	expected := "<p>bar</p>"
	body := w.Body.String()
	if body != expected {
		t.Errorf("rendered template failed: %s != %s", expected, body)
	}

}
