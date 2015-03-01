package kwiscale

// Manage error handling
// TODO: allow user to define error pages

import (
	"log"
	"net/http"
)

// HandleError write error code in header + message
func HandleError(code int, response http.ResponseWriter, req *http.Request, err error) {
	if debug {
		log.Print(code, response, req, err)
	}
	errstring := ""
	if err != nil {
		errstring = err.Error()
	}
	http.Error(response, errstring, code)
}
