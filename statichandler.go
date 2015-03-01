package kwiscale

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

// StaticHandler handle static files handlers. Use App.SetStatic(path) that create the static handler
type staticHandler struct {
	RequestHandler
	cacheEnabled bool
}

var cache = make(map[string][]byte)

func (s *staticHandler) putInCache(c []byte, f string) {
	if s.cacheEnabled {
		cache[f] = c
	}
}

func (s *staticHandler) Get() {
	file := s.Vars["file"]
	file = filepath.Join(s.app.Config.StaticDir, file)

	var content []byte
	var err error

	if s.cacheEnabled {
		content = cache[file]
	}
	if content == nil {
		content, err = ioutil.ReadFile(file)
		if err != nil {
			HandleError(http.StatusNotFound, s.getResponse(), s.getRequest(), err)
			return
		}
		// save in cache after all
		defer s.putInCache(content, file)
	}

	mimetype := mime.TypeByExtension(filepath.Ext(file))
	s.Response.Header().Add("Content-Type", mimetype)
	s.Write(content)
}
