package kwiscale

import (
	"log"
	"net/http"
	"net/url"
)

// Deprecated functions that will disapear

// GetApp returns the app that holds this handler.
//
// DEPRECATED -- see App()
func (b *BaseHandler) GetApp() *App {
	log.Println("[WARN] GetApp() is deprecated, please use App() method instead.")
	return b.App()
}

// GetResponse returns the current response.
//
// DEPRECATED -- see Response()
func (b *BaseHandler) GetResponse() http.ResponseWriter {
	log.Println("[WARN] GetResponse() is deprecated, please use Response() method instead.")
	return b.Response()
}

// GetRequest returns the current request.
//
// DEPRECATED -- see Request()
func (b *BaseHandler) GetRequest() *http.Request {
	log.Println("[WARN] GetRequest() is deprecated, please use Request() method instead.")
	return b.Request()
}

// GetPost return the post data for the given "name" argument.
//
// DEPRECATED -- see PostVar()
func (b *BaseHandler) GetPost(name string) string {
	log.Println("[WARN] GetPost() is deprecated, please use PostVar() method instead.")
	return b.PostValue(name, "")
}

// GetURL return an url based on the declared route and given string pair.
//
// DEPRECATED -- see URL()
func (b *BaseHandler) GetURL(s ...string) (*url.URL, error) {
	log.Println("[WARN] GetURL() is deprecated, please use URL() method instead.")
	return b.URL(s...)
}

// GetPayload returns the Body content.
//
// DEPRECATED - see Payload()
func (b *BaseHandler) GetPayload() []byte {
	log.Println("[WARN] GetPauload() is deprecated, please use Payload() method instead.")
	return b.Payload()
}

// GetPostValues returns the entire posted values.
//
// DEPRECATED - see PostValues()
func (b *BaseHandler) GetPostValues() url.Values {
	log.Println("[WARN] GetPostValues() is deprecated, please use PostValues() method instead.")
	return b.PostValues()
}

// GetJSONPayload unmarshal body to the "v" interface.
//
// DEPRECATED - see JSONPayload()
func (b *BaseHandler) GetJSONPayload(v interface{}) error {
	log.Println("[WARN] GetJSONPayload() is deprecated, please use JSONPayload() method instead.")
	return b.JSONPayload(v)
}

// GlobalCtx Returns global template context.
//
// Deprecated: use handler.App().Context instead
func (b *RequestHandler) GlobalCtx() map[string]interface{} {
	log.Println("[WARN] GlobalCtx() is deprecated, please use App().Context instead.")
	return b.App().Context
}
