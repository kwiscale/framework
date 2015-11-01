/*
Kwiscale is a framework that provides Handling System. That means that you will be able to create
Handlers that handles HTTP Verbs methods. Kwiscale can handle basic HTTP Verbs as Get, Post, Delete, Patch, Head and Option.

Kwiscale can also serve Websockets.

Kwiscale provides a basic template system based on http/template from Go SDK and has got a plugin system to provides other template engines.


The built-in template allows to override templates. Example:

	// templates/main.go
	<body>
	{{ template "CONTENT" .}}
	</body>

	// template/home/index.go (remove the "\" in override directive)
	{{/* override "main.go" *\/}}

	{{ define "CONTENT" }}
	<p>Hello workd</p>
	{{ end }}


Also, built-in template provides two functions:

	- url: gives url of a handler with parameters
	- static: gives static resource url

Example:

	<link rel="stylesheet" href="{{static "/css/styles.css"}}" />
	<a href="{{ url "home" }}">Home</a>
	<a href="{{ url "user" "id" 12345  }}">User</a>
	<a href="{{ url "handlers.UserHandler" "id" 12345  }}">User</a>

For "url" function, the first argument can be:

	- alias for handler
	- handler name

Others arguments are the pair "key value".

See http://gopkg.in/kwiscale/template-pongo2.v1 to use Pongo2.

To handle HTTP Verbs:

	type HomeHandler struct {kwiscale.RequestHandler}
	func (home *HomeHandler) Get(){
		// This will respond to GET call
		home.WriteString("Hello !")
	}

	func main(){
		app := kwiscale.NewApp(nil)
		app.AddRoute("/home", &HomeHandler{})
		// Default listening on 8000 port
		app.ListenAndServe()
	}


To be able to use configuration file (yaml), you MUST register handlers. The common way to do is to use "init()" function in you handlers package:

	package handlers

	import "gopkg.in/kwiscale/framework.v1"

	func init(){
		kwiscale.Register(&HomeHandler{})
	}

	type HomeHandler struct {kwiscale.HomeHandler}
	//...

Note: if you're using kwiscale CLI, `kwiscale new handler` command create the register call for you.

Kwiscale provide a way to have method with parameter that are mapped from the given route.

Example:

	app.AddRoute(`/user/{name:.+/}{id:\d+}/`, UserHandler{})

	//...
	type UserHandler struct {kwiscale.RequestHandler}
	func (handler *UserHandler) Get(name string, id int) {
		// ...
	}

Note that parameters names for Get() method are not used to map url values. Kwiscale maps values repecting the order found in the route. So you may declare Get() method like this:


	func (handler *UserHandler) Get(a string, b int) {
		// ...
	}

Anyway, you always may use "UserHandler.Vars":  UserHandler.Vars["id"] and UserHandler.Vars["name"] that are `string` typed.


You may use Init() and Destroy() method that are called before and after HTTP verb invocation. You may, for example, open database connection in "Init" and close the connection in "Destroy".


Kwiscale provides a CLI:

	go get gopkg.in/framework/kwiscale
	kwiscale --help

	NAME:
	   kwiscale - tool to manage kwiscale application

	USAGE:
	   kwiscale [global options] command [command options] [arguments...]

	VERSION:
	   0.0.1

	COMMANDS:
	   new		Generate resources (application, handlers...)
	   generate	Parse configuration and generate handlers, main file...
	   help, h	Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --project "kwiscale-app"	project name, will set $GOPATH/src/[projectname] [$KWISCALE_PROJECT]
	   --handlers "handlers"	handlers package name [$KWISCALE_HANDLERS]
	   --help, -h			show help
	   --generate-bash-completion
	   --version, -v		print the version

See http://readthedocs.org/projects/kwiscale/

*/
package kwiscale
