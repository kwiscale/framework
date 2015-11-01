package kwiscale

import (
	"errors"
	"fmt"
	"text/template"
)

var (
	ErrNotFound       = errors.New("Not found")
	ErrNotImplemented = errors.New("Not implemented")
	ErrInternalError  = errors.New("Internal server error")
)

// HTTPErrorHandler interface.
type HTTPErrorHandler interface {
	// Error returns the error.
	GetError() error
	// Details returns some detail inteface.
	Details() interface{}
	// Status returns the http status code.
	Status() int

	setStatus(int)
	setError(error)
	setDetails(interface{})
}

// ErrorHandler is a basic error handler that
// displays error in a basic webpage.
type ErrorHandler struct {
	RequestHandler
	status  int
	err     error
	details interface{}
}

func (dh *ErrorHandler) setStatus(s int) {
	dh.status = s
}

func (dh *ErrorHandler) setError(err error) {
	dh.err = err
}

func (dh *ErrorHandler) setDetails(d interface{}) {
	dh.details = d
}

// GetError returns error that was set by handlers.
func (dh *ErrorHandler) GetError() error {
	return dh.err
}

// Details returns details or nil if none.
func (dh *ErrorHandler) Details() interface{} {
	return dh.details
}

// Status returns the error HTTP Status set by handlers.
func (dh *ErrorHandler) Status() int {
	return dh.status
}

// Get shows a standard error in HTML.
func (dh *ErrorHandler) Get() {

	tpl := `<!doctype html>
<html>
	<head>
		<title>ERROR {{.Status}}</title>
		<style>
			html, body {
				font-family: Sans, sans-serif;
			}
			main {
				margin: auto;
				width: 80%;
				border: 2px solid #880000;
				padding: 2em;
			}
		</style>
	</head>
	<body>
	<main>
        <h1>ERROR {{ .Status}}</h1>
        <p>{{ .Error }}</p>
        <pre>{{ range .Details }}{{ . }}{{ end }}</pre>
	</main>
	</body>
</html>`

	t, err := template.New("error").Parse(tpl)
	if err != nil {
		fmt.Println(err)
		return
	}

	dh.Response().WriteHeader(dh.Status())
	t.Execute(dh.Response(), map[string]interface{}{
		"Status":  dh.Status(),
		"Error":   dh.GetError(),
		"Details": dh.Details(),
	})
}
