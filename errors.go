package kwiscale

// Manage error handling
// TODO: allow user to define error pages

import (
	"log"
	"net/http"
)

// HandleError write error code in header + message
func HandleError(code int, response http.ResponseWriter, req *http.Request, err error) {
	log.Print(code, response, req, err)
	/*response.WriteHeader(code)
	response.Write([]byte(http.StatusText(code)))
	if err != nil {
		response.Write([]byte(fmt.Sprintf("\n<pre>%+v</pre>", err)))
	}
	*/
	errstring := ""
	if err != nil {
		errstring = err.Error()
	}
	log.Print("ERROR ?", req)
	http.Error(response, errstring, code)
}
