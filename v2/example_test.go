package kwiscale_test

import kwiscale "./"

type HomeHandler struct {
	kwiscale.Handler
}

func ExampleServer_Route() {

	server := kwiscale.NewServer()
	server.Route("/homepage", func() kwiscale.IHandler { return new(HomeHandler) })

}
