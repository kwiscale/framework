package kwiscale

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Log(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}

func Error(v ...interface{}) {
	msg := []interface{}{"[ERROR]"}
	msg = append(msg, v...)
	log.Println(msg...)
}

// getMatchRoute returns the better handler that matched request url.
// It returne handlername, mux.Route and mux.RouteMatch.
func getBestRoute(app *App, r *http.Request) (string, *mux.Route, mux.RouteMatch) {

	var points int
	var handlerToInstanciate = ""
	var matchedRoute *mux.Route
	var routeMatch mux.RouteMatch
	for route, handler := range app.handlers {
		var match mux.RouteMatch
		if route.Match(r, &match) {
			plength := len(match.Vars)

			// if number of url path part is lower than last matched route
			// then try next...
			Log("Matches route vars", match.Vars)
			if plength < points {
				continue
			}
			points = plength

			// construct handler from its name
			Log(fmt.Sprintf("Route matches %#v\n", route))
			Log("Handler to fetch", handler)
			handlerToInstanciate = handler
			matchedRoute = route
			routeMatch = match
		}
	}

	return handlerToInstanciate, matchedRoute, routeMatch
}
