kwiscale
========

[![GoDoc](https://godoc.org/gopkg.in/kwiscale/framework.v1?status.svg)](https://godoc.org/gopkg.in/kwiscale/framework.v1)


Web Middleware for Golang

At this time, Kwiscale is at the very begining of developpement. But you can test and give'em some pull-request to improve it.


Features
========

- Implement your handlers as structs with HTTP Verbs as method
- Plugin system for template engines, session engines and ORM
- Use gorilla to manipulate routes and mux
- Handler spawned with concurrency 


How to use
==========

Install with "go get" command:

    go get gopkg.in/kwiscale/framework.v1

Create a project:

    mkdir ~/myproject && cd ~/myproject

The common way to create handlers is to append a package::

    mkdir handlers
    vim handlers/index.go

Let's try an example:

```go
package handlers

import "gopkg.in/kwiscale/framework.v1"

// this is the Index Handler that
// is composed by a RequestHandler
type IndexHandler struct {
    // compose your handler with kwiscale.Handler
    kwiscale.RequestHandler
}

// Will respond to GET request. "params" are url params (not GET and POST data)
func (i *IndexHandler) Get () {
    i.WriteString("Hello !" + i.Vars["userid"])
}
```

Then in you main.go::

```go

package main

import (
    "gopkg.in/kwiscale/framework.v1"
    "./handlers"
)

// HomeHandler
type HomeHandler struct {
	kwiscale.RequestHandler
}

// Get respond to GET request
func (h *HomeHandler) Get (){
	h.WriteString("reponse to GET home")
}

// Another handler
type OtherHandler struct {
	kwiscale.RequestHandler
}

func (o *OtherHandler) Get (){
	// read url params
	// it always returns a string !
	userid := o.Vars["userid"]
	o.WriteString(fmt.Printf("Hello user %s", userid))
}


func main() {
	kwiscale.DEBUG = true
	app := kwiscale.NewApp(&kswicale.Config{
        Port: ":8000",
    })
	app.AddRoute("/", HomeHandler{})
	app.AddRoute("/user/{userid:[0-9]+}", OtherHandler{})
    app.ListenAndServe()

    // note: App respects http.Mux so you can use:
    // http.ListenAndServe(":9999", app)
    // to override behaviors, testing, or if your infrastructure
    // restricts this usage
}
```


Then run:

    go run main.go

Or build your project:
    
    go build main.go
    ./main


- Visit http://127.0.0.1:8000/ and you should see "Hello FOO"
- Visit http://127.0.0.1:8000/user/12345 and you should see "Hello user 12345"


The Kwiscale way ?
==================

Kwiscale let you declare Handler methods with the HTTP method. This allows you to declare:

* Get()
* Post()
* Head()
* Delete()
* Put()
* Patch()


Basic Templates
===============

Kwiscale provides a "basic" template engine that use `http/template`. Kwiscale only add a "very basic template override system".

If you plan to have a complete override system, please use http://gopkg.in/kwiscale/template-pongo2.v1 that implements pango2 template.

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
        {{/* Remember to use "." as context */}}
        {{ template "CONTENT" . }}
    </body>
</html>
```
Now create templates/home directory:
    
    mkdir templates/home

Create templates/home/welcome.html:

    {{/* override "main.html" */}}

    {{ define "CONTENT" }}
        This the welcome message {{ .msg }}
    {{ end }}

This template overrides "main.html" (in `./templates/` directory) and append "CONTENT" template definition. So, the "CONTENT" block will appear at `template "CONTENT"` in "main.html". That's all.  

In handlers/index.go you may now ask for template rendering:

```go
func (h *IndexHandler) Get() {
    h.Render("home/welcome.html", map[string]string{
        "title" : "Welcome !!!",
        "msg"   : "Hello you",
    })
}
```

You can override template directory using App configuration passed to the constructor:

```go
app := kwiscale.NewApp(&kswiscale.Config{
    TemplateDir: "./my-template-dir",
})

```


TODO
====

Features in progress:

- [ ] Database ORM interface
- [ ] Custom Error handler
