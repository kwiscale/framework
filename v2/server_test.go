package kwiscale

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// base server
var server *Server

// initialize server
func initServer() {
	server = NewServer()
}

// The test handler used in this test
type TestHandler struct {
	Handler
}

// The GET response that is ok
func (h *TestHandler) Get(c map[string]string) {
	h.Response.WriteString("test ok")
}

// Test if bad page returns 404 even if other page is ok
func Test_Respond_404(t *testing.T) {

	initServer()
	server.Route("/correct", func() IHandler { return new(TestHandler) })

	req, _ := http.NewRequest("GET", "/badpage", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatal("Bad status", w.Code)
	}

	req, _ = http.NewRequest("GET", "/correct", nil)
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal("Bad status", w.Code)
	}
}

// Check only a 200 OK
func Test_Respond_OK(t *testing.T) {
	initServer()
	server.Route("/home", func() IHandler { return new(TestHandler) })

	req, _ := http.NewRequest("GET", "/home", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal("Bad status", w.Code)
	}

}

func Test_NotImplemented(t *testing.T) {
	initServer()
	server.Route("/home", func() IHandler { return new(TestHandler) })
	req, _ := http.NewRequest("GET", "/home", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal("Bad status", w.Code)
	}

	req, _ = http.NewRequest("POST", "/home", nil)
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusNotImplemented {
		t.Fatal("Bad status", w.Code)
	}
}

func Test_BadRequest(t *testing.T) {
	initServer()
	server.Route("/home", func() IHandler { return new(TestHandler) })
	req, _ := http.NewRequest("GET", "/home", nil)
	w := httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatal("Bad status", w.Code)
	}

	req, _ = http.NewRequest("WTF", "/home", nil)
	w = httptest.NewRecorder()
	server.Router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatal("Bad status", w.Code)
	}
}
