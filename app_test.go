package kwiscale

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHandler struct{ RequestHandler }

func (th *TestHandler) Get() {
	th.WriteString("Hello")
}

type TestReverseRoute struct{ RequestHandler }

func (th *TestReverseRoute) Get() {
	route, err := th.GetURL("category", "foo")
	if err != nil {
		fmt.Println(err)
	}
	th.WriteString(route.String())
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

func TestRouteReverse(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/product/test", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/product/{category:.+}", TestReverseRoute{})
	app.ServeHTTP(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	if string(resp) != "/product/foo" {
		t.Fatal(resp, "!=", "/product/foo")
	}
	<-app.SoftStop()
}

func TestReverseURLFromApp(t *testing.T) {
	app := initApp(t)
	app.AddRoute("/product/{category:.+}", TestReverseRoute{})
	r, _ := http.NewRequest("GET", "http://example.com/product/test", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)

	u, err := app.GetRoute("kwiscale.TestReverseRoute").URL("category", "second")
	if err != nil {
		t.Fatal(err)
	}
	if u.String() != "/product/second" {
		t.Fatal(u.String(), "!=", "/product/second")
	}
	<-app.SoftStop()
}
