package kwiscale

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// StaticHandler handle static files handlers. Use App.SetStatic(path) that create the static handler
type staticHandler struct {
	RequestHandler
	cacheEnabled bool
	prefix       string
}

var cache = make(map[string][]byte)

func (s *staticHandler) putInCache(f string) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		cache[f] = content
	}
}

// Use http.FileServer to serve file after adding ETag.
func (s *staticHandler) Get() {
	file := s.Vars["file"]
	file = filepath.Join(s.app.Config.StaticDir, file)

	// control or add etag
	if etag, err := eTag(file); err == nil {
		s.response.Header().Add("ETag", etag)
	}

	fs := http.FileServer(http.Dir(s.prefix))
	fs.ServeHTTP(s.Response(), s.Request())
}

// Get a etag for the file. It's constuct with a md5 sum of
// <filename> + "." + <modification-time>
func eTag(file string) (string, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return "", err
	}

	s := md5.Sum([]byte(stat.Name() + "." + stat.ModTime().String()))
	return fmt.Sprintf("%x", s), nil
}
