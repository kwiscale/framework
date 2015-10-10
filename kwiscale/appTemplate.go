package main

const TPLAPP = `package main

import (
	"gopkg.in/kwiscale/framework.v0"
	"{{.Project}}/{{.HandlersPKG}}"
)

func main(){

	app := kwiscale.NewAppFromConfigFile()

	//@routes@ -- DO NOT REMOVE THIS COMMENT
	//@end routes@ -- DO NOT REMOVE THIS COMMENT

	app.ListenAndServe()
}
`

// the main.go template.
const TPLHANDLER = `package {{.HandlersPKG}}

import (
	"gopkg.in/kwiscale/framework.v0"
)

type {{.Handler}} struct { kwiscale.RequestHandler }

`

// Generate:
//	app.AddRoute(route, handler)
// or
//	app.AddNamedRoute(route, handler, alias)
const TPLADDNAMEDROUTE = "	app.Add{{if .Route.alias}}Named{{end}}Route(`{{.Route.route}}`,{{.Route.handler}}{}" +
	"{{if .Route.alias}}, \"{{ .Route.alias }}\"{{end}}" +
	")"
