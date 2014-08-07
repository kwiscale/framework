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

func main() {
    // Get a server
    server := kwiscale.NewServer()

    // Link route /home/XXX to IndexHander factory
    // "username" will be passed inside the map params to Get method.
    server.Route("/home/{username:.*}", func IHandler {return &handlers.IndexHandler})
    
    // Start to listen on 0.0.0.0:8081
    server.Listen(":8081")
}
```

Note: Route() prototype is:

    kwiscale.Route(url string, factory kwiscale.Factory)

kwiscale.Factory type is "func IHandler {}". The factory is only a function that returns pointer on your Handler. In above example:

```go
// explain factory system
var factory kwiscale.Factory 
factory = func IHandler {
	return new(IndexHandler) // or return &IndexHandler{}
}

// Then:
kwiscale.Route("/home/{username:.*}", factory)
```

We use factory to allows new handler at time. Kwiscale use a factory stack that spawn some handler and keep it in a channel.

Then run:

    go run main.go

Or build your project:
    
    go build main.go
    ./main


Visit http://127.0.0.1:8081/home/FOO and you should see "Hello FOO"


The Kwiscale way ?
==================

Kwiscale let you declare Handler methods with the HTTP method. This allows you to declare:

* Get(map[string]string)
* Post(map[string]string)
* Head(map[string]string)
* Delete(map[string]string)
* Put(map[string]string)


The map[string]string parameter represents what is given to the URL. It uses Gorilla "mux.Route" to implement that route syntax. Note that parameters are not GET (and not POST) values, but are parts of the path. Get/Post values can be read from handler.Request object (TODO: documentation).

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


