package main

// The main package template.
const TPLAPP = `package main

import (
	"gopkg.in/kwiscale/framework.v0"
	_ "{{.Project}}/{{.HandlersPKG}}"
)

func main(){
	app := kwiscale.NewAppFromConfigFile()
	app.ListenAndServe()
}
`

// the main.go template.
const TPLHANDLER = `package {{.HandlersPKG}}

import (
	"gopkg.in/kwiscale/framework.v0"
)


func init(){
	kwiscale.Register(&{{.Handler}}{})
}

type {{.Handler}} struct { kwiscale.RequestHandler }

`

// Generate a handler register call.
const TPLADDNAMEDROUTE = "	kwiscale.Register(&{{.Route.handler}}{})"
