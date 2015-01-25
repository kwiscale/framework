kwiscale
========

Web Middleware for Golang

At this time, Kwiscale is at the very begining of developpement. But you can test and give'em some pull-request to improve it.

Check documentation: http://godoc.org/github.com/metal3d/kwiscale

How to use
==========

Install with "go get" command:

    go get github.com/metal3d/kwiscale

Create a project:

    mkdir ~/myproject && cd ~/myproject

The common way to create handlers is to append a package::

    mkdir handlers
    vim handlers/index.go

Let's try an example:

```go
package handlers

import "github.com/metal3d/kwiscale"

// this is the Index Handler that
// is composed by a RequestHandler
type IndexHandler struct {

    // compose your handler with kwiscale.Handler
    kwiscale.Handler
}

// Will respond to GET request. "params" are url params (not GET and POST data)
func (this *IndexHandler) Get (params map[string]string) {
    this.Write("Hello !" + params["username"])
}
```

Then in you main.go::

```go

package main

import (
    "github.com/metal3d/kwiscale"
    "./handlers"
)

// HomeHandler
type HomeHandler struct {
	kwiscale.RequestHandler
}

// Get respond to GET request
func (h *HomeHandler) Get (){
	h.Response.WriteString("reponse to GET home")
}

// Another handler
type OtherHandler struct {
	kwiscale.RequestHandler
}

func (o *OtherHandler) Get (){
	// read url params
	// it always return a string !
	userid := o.Vars["userid"]
	o.Response.WriteString(userid)
}


func main() {
	kwiscale.DEBUG = true
	app := kwiscale.NewApp()
	app.AddRoute("/", HomeHandler{})
	app.AddRoute("/user/{userid:[0-9]+}", OtherHandler{})
	http.ListenAndServe(":8000", app)
}
```


Then run:

    go run main.go

Or build your project:
    
    go build main.go
    ./main


Visit http://127.0.0.1:8000/ and you should see "Hello FOO"


The Kwiscale way ?
==================

Kwiscale let you declare Handler methods with the HTTP method. This allows you to declare:

* Get()
* Post()
* Head()
* Delete()
* Put()
* Patch()


Templates
=========

Right now, basic templates (that are not so basics...) from golang sdk is the default template engine. But you can use Pango2 or other template engine. 

Kwiscale adds "override" system to import other templates. That way, you can call a subtemplate.

See the following example.

Append templates directory:
    
    mkdir templates

Then create templates/main.html:

```html
<!DOCTYPE html>
<html>
    <head>
        <title>{{ if .title }}{{.title}}{{ else }} Default title {{ end }}</title>
    </head>
    <body>
        {{ template "CONTENT" . }}
    </body>
</html>
```
Now create templates/home directory:
    
    mkdir templates/home

Create templates/home/welcome.html:
    
    {{ override "main.html" }}

    {{ define "CONTENT" }}
        This the welcome message {{ .msg }}
    {{ end }}

This template override "main.html". 


In handlers/index.go:

```go
func (h *IndexHandler) Get() {
    h.Response.Render("home/welcome.html", kwiscale.ContextParams{
        "title" : "Welcome !!!",
        "msg"   : "Hello you",
    })
}
```

To use other template, you should create struct that interface kwiscale.ITemplateEngine. You only have to implement one method: Render(h kwiscale.Handler, templatename string, ctx kwiscale.ContextParams). This method should write template with h.Write(...).

Note that ContextParams is a simple mam[interface{}]interface{}.


