package kwiscale

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHandler struct{ RequestHandler }

func (th *TestHandler) Get() {
	th.WriteString("Hello")
}

func initApp(t *testing.T) *App {
	conf := &Config{}
	app := NewApp(conf)
	return app
}

func TestCloser(t *testing.T) {
	app := initApp(t)
	app.AddRoute("/foo", TestHandler{})
	<-app.SoftStop()
}

// TestSimpleRequest should respond whit 200 and print "hello"
func TestSimpleRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/foo", TestHandler{})
	app.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("HTTP Status is not ok: ", w.Code)
	}

	resp, _ := ioutil.ReadAll(w.Body)
	if string(resp) != "Hello" {
		t.Error("Handler didn't respond with 'hello': ", string(resp))
	}
	<-app.SoftStop()
}

// Try to call a bad route
func TestBadRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/bad", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/foo", TestHandler{})
	app.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Error(`HTTP Status is not "not found": `, w.Code)
	}

	<-app.SoftStop()
}
