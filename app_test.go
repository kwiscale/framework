package kwiscale

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Let each handler to get "*testing.T"
var T = make(map[*App]*testing.T)

// A test handler (simple).
type TestHandler struct {
	RequestHandler
}

func (th *TestHandler) Get() {
	th.WriteString("Hello")
}

// Handler to test reversed route.
type TestReverseRoute struct{ RequestHandler }

func (th *TestReverseRoute) Get() {
	// Test to get route from app
	app := th.App()
	t := T[app]

	u, err := app.GetRoute("kwiscale.TestReverseRoute").URL("category", "test")
	if err != nil {
		t.Error("Route from app returns error:", err)
	}
	if u.String() != "/product/test" {
		t.Error("Route from app is not /product/test: ", u)
	}

	route, err := th.GetURL("category", "foo")
	if err != nil {
		fmt.Println(err)
	}
	th.WriteString(route.String())
}

// Create app.
func initApp(t *testing.T) *App {
	conf := &Config{}
	app := NewApp(conf)
	T[app] = t
	return app
}

// Test the app "soft close".
func TestCloser(t *testing.T) {
	app := initApp(t)
	app.AddRoute("/foo", &TestHandler{})
	<-app.SoftStop()
}

// TestSimpleRequest should respond whit 200 and print "hello".
func TestSimpleRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/foo", &TestHandler{})
	app.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("HTTP Status is not ok: ", w.Code)
	}

	resp, _ := ioutil.ReadAll(w.Body)
	if string(resp) != "Hello" {
		t.Error("Handler didn't respond with 'hello': ", string(resp))
	}
}

// Try to call a bad route.
func TestBadRequest(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/bad", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/foo", &TestHandler{})
	app.ServeHTTP(w, r)
	if w.Code != http.StatusNotFound {
		t.Error(`HTTP Status is not "not found": `, w.Code)
	}
}

// Test the reverse route to get url.
func TestRouteReverse(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/product/test", nil)
	w := httptest.NewRecorder()

	app := initApp(t)
	app.AddRoute("/product/{category:.+}", &TestReverseRoute{})
	app.ServeHTTP(w, r)

	resp, _ := ioutil.ReadAll(w.Body)

	if string(resp) != "/product/foo" {
		t.Fatal(resp, "!=", "/product/foo")
	}
}

// Test to get reverse route from app.
func TestReverseURLFromApp(t *testing.T) {
	app := initApp(t)
	app.AddRoute("/product/{category:.+}", &TestReverseRoute{})

	u, err := app.GetRoute("kwiscale.TestReverseRoute").URL("category", "second")
	if err != nil {
		t.Fatal(err)
	}
	if u.String() != "/product/second" {
		t.Fatal(u.String(), "!=", "/product/second")
	}
}

// BUG: This is a limit case we really need to study
func _TestBestRoute(t *testing.T) {

	r, _ := http.NewRequest("GET", "http://example.com/test/route", nil)
	//w := httptest.NewRecorder()

	app := initApp(t)

	app.AddRoute("/test/route", &TestHandler{})
	app.AddRoute("/{p:.*}", &TestReverseRoute{})

	name, _, _ := getBestRoute(app, r)

	if name != "kwiscale.TestReverseRoute" {
		t.Fatal("For /test/route, the handler that matches should be kwiscale.TestReverseRoute and not", name)
	}
}
