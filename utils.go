package kwiscale

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Log print logs on STDOUT if debug is activated.
func Log(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}

// Error prints error on STDOUT.
func Error(v ...interface{}) {
	msg := []interface{}{"[ERROR]"}
	msg = append(msg, v...)
	log.Println(msg...)
}

// getMatchRoute returns the better handler that matched request url.
// It returns handlername, mux.Route and mux.RouteMatch.
//
// Rule is simple:
//	- if A url path length is greater than B url path length, B wins
//  - if A url path length is equal that B url path length, then:
//		- if A number of path vars is greater than B number if path vars, B wins
//
// So:
//
//	- /path/A vs /B => B wins
//  - /path/A/{foo:.*} vs /path/B/bar => B wins
func getBestRoute(app *App, r *http.Request) (handlerName string, route *mux.Route, match mux.RouteMatch) {

	points := -1

	for handlerRoute, handler := range app.handlers {
		var routematch mux.RouteMatch
		if handlerRoute.Match(r, &routematch) {
			plength := len(strings.Split(r.URL.Path, "/")) + len(routematch.Vars)
			if points == -1 {
				points = plength
			} else if plength > points {
				continue
			}

			Log("Matches route vars", handler, routematch.Vars)
			points = plength

			Log("Handler to fetch", handler)
			handlerName = handler.handlername
			route = handlerRoute
			match = routematch
		}
	}

	return
}
