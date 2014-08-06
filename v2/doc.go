/*
Kwiscale is a standard framework made in Go using some external librairies to manage routes, forms, sessions... It
tries to let developpers to implement web application as python WSGI handler set. Handler can implement HTTP verbs
as functions (Get(), Post(), Put()...).

The sequence is to:

	* create a server instance with NewServer()
	* implement Handlers
	* add routes from URL to handler using Server.Route()

Examples

Creating a server:

	server := kwiscale.NewServer()

A basic Handler:

	type GreetingHandler struct {
		kwiscale.Handler
	}

	// HTTP Verb should have map[string]string parameter to get route params
	func (g *GreetingHandler) Get (params map[string]string) {
		g.Write("Hello " + params["name"]) // write name param
	}

Now, create the route:

	// Route maps a URL route to a Factory that returns new IHandler
	server.Route("/greeting/{name}", func() IHandler{return new(GreetingHandler)})

	// Then, serve on 8000 port
	server.Listen(":8000")

Using Factory

Route() method needs to have a Factory as second parameter. Factory is a simple closure that will be called only 
when handler is needed. 

When you create a route (using Server.Route() method)
a handler stack is created (buffered chan) and a goroutine is launched.  This goroutine fill the chan calling this
factory. When a HTTP call is routed, the corresponding channel is read to get an instanciated handler from memory.
Then, the goroutine fill again the chan to replace the using handler.

That way, a certain number of handlers is
instanciated in memory *before* to be used by HTTP calls. By default, server create one stack of 10 handlers per route.

You can change this value after server creation:

	server.HandlerStackSize = 50


*/
package kwiscale
