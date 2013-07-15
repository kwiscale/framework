package kwiscale

import (
    "testing"
	"fmt"
	"net/http"
	"net/http/httptest"
)

type TestHandler struct {
    RequestHandler
}

func (t *TestHandler) Get() {
    t.Write("Response ok")
}

func Init () {

    h := TestHandler{}
    h.Routes = []string{"/foo"}

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

	fmt.Printf("%d - %s", w.Code, w.Body.String())
    if w.Code != http.StatusOK {
        t.Fatalf("Status is not 200: %v", w.Code)
    }
}

// test an unknown url
func Test404 (t *testing.T) {
    Init()

    req, err := http.NewRequest("GET", "http://example.com/bar", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()
	dispatch(w, req)

	fmt.Printf("%d - %s", w.Code, w.Body.String())
    if w.Code != http.StatusNotFound {
        t.Fatalf("Status is not 404: %v", w.Code)
    }
}
