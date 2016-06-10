kwiscale
========

[![Join the chat at https://gitter.im/kwiscale/framework](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/kwiscale/framework?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://drone.io/github.com/kwiscale/framework/status.png)](https://drone.io/github.com/kwiscale/framework/latest)
[![Coverage Status](https://coveralls.io/repos/kwiscale/framework/badge.svg?branch=master&service=github)](https://coveralls.io/github/kwiscale/framework?branch=master)
[![Documentation Status](https://readthedocs.org/projects/kwiscale/badge/?version=latest)](http://kwiscale.readthedocs.org/en/latest/?badge=latest)
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
    go get gopkg.in/kwiscale/framework.v1/cmd/kwiscale

Create a project:

    kwiscale new app myapp
    cd $GOPATH/myapp

Create a handler
    kwiscale new handler index / homepage

Open generated `handlers/index.go` and append "Get" method:

```go
package handlers

import (
	"gopkg.in/kwiscale/framework.v1"
)


func init(){
	kwiscale.Register(&IndexHandler{})
}

type IndexHandler struct { kwiscale.RequestHandler }

func (handler *IndexHandler) Get(){
    handler.WriteString("Hello you !")
}

```

And run the app !:

```bash
go run *.go
```

Now, go to http://127.0.0.1:8000 - you should see "Hello You!". If not, check "kwiscale.yml" file if port is "8000", check error log, and so on.

If really there is a problem, please submit an issue.


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
