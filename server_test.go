package kwiscale

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestHandler struct {
	RequestHandler `route:"/foo"`
}

func (t *TestHandler) Get() {
	t.Write("Response ok")
}

func (t *TestHandler) New () IRequestHandler {
    return new(TestHandler)
}

func Init() {

	h := TestHandler{}
	AddHandler(&h)
}

// test a correct call
func Test200(t *testing.T) {
	Init()

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	dispatch(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Status is not 200: %v\n", w.Code)
	}
}

// test an unknown url
func Test404(t *testing.T) {
	Init()

	req, err := http.NewRequest("GET", "http://example.com/bar", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	dispatch(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Status is not 404: %v\n", w.Code)
	}
}
