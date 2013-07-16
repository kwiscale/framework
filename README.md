kwiscale
========

Web Middleware for Golang

At this time, Kwiscale is at the very begining of developpement. But you can test and give'em some pull-request to improve it.


Check documentation: http://godoc.org/github.com/metal3d/kwiscale

How to use
==========

Install with "go get" command::

    go get github.com/metal3d/kwiscale

Create a project::

    mkdir ~/myproject && cd ~/myproject

The common way to create handlers is to append a package::

    mkdir handlers
    vim handlers/index.go


Let's try an example::

    package handlers

    import "github.com/metal3d/kwiscale"

    // this is the Index Handler that
    // is composed by a RequestHandler
    type IndexHandler struct {
        // note this, we now use Tag to declare route Inside the Handler
        kwiscale.RequestHandler `route:"/home"`
    }

    func (this *IndexHandler) Get () {
        this.Write("Hello !")
    }

Then in you main.go::
    
    package main
    
    import (
        "github.com/metal3d/kwiscale"
        "./handlers"
    )

    func main() {
        h := handlers.IndexHandler{}
        kwiscale.AddHandler(&h)

        kwiscale.Serve(":8081") //listen :8081
    }


Then run::

    go run main.go


Visit http://127.0.0.1:8081/home and you should see "Hello"


What gives Kwiscale ?
=====================

Kwiscale let you declare Handler methods with the HTTP method. This allows you to declare:

* Get()
* Post()
* Head()
* Delete()
* Put()

Routes are regexp and captured elements are stocked in handler.UrlParams that is a []string.

The handler.Render method takes template path and context to implement. But it adds a nice way to "overide" templates.

Let's take an example. 

Append templates directory::
    
    mkdir templates

Then create templates/main.html::

    <!DOCTYPE html>
    <html>
        <head>
            <title>{{ if .title }}{{.title}}{{ else }} Default title {{ end }}</title>
        </head>
        <body>
            {{ template "CONTENT" . }}
        </body>
    </html>

Now create templates/home directory::
    
    mkdir templates/home

Create templates/home/welcome.html::
    
    {{ override "main.html" }}

    {{ define "CONTENT" }}
        This the welcome message {{ .msg }}
    {{ end }}

In handlers/index.go::

    func (this *IndexHandler) Get() {
        this.Render("home/welcome.html", map[string]string{
            "title" : "Welcome !!!",
            "msg"   : "Hello you",
        })
    }


