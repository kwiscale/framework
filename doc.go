/*
Kwiscale is a framework that provides Handling System. That means that you will be able to create
Handlers that handles HTTP Verbs methods. Kwiscale can handle basic HTTP Verbs as Get, Post, Delete, Patch, Head and Option.

Kwiscale can also serve Websockets.

It provides a basic template system based on http/template from Go SDK and has got a plugin system to provides other template
engines

See http://github.com/metal3d/kwiscale-template-pongo2 to use Pongo2.

To handle HTTP Verbs:

	type HomeHandler struct {kwiscale.RequetHandler}
	func (home *HomeHandler) Get(){
		// This will respond to GET call
		home.WriteString("Hello !")
	}

	func main(){
		app := kwiscale.NewApp(nil)
		app.AddRoute("/home", HomeHandler{})
		// Default listening on 8000 port
		app.ListenAndServe()
	}

*/
package kwiscale
