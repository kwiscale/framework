package kwiscale

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// StaticHandler handle static files handlers. Use App.SetStatic(path) that create the static handler
type staticHandler struct {
	RequestHandler
}

// Use http.FileServer to serve file after adding ETag.
func (s *staticHandler) Get() {
	file := s.Vars["file"]
	file = filepath.Join(s.app.Config.StaticDir, file)

	// control or add etag
	if etag, err := eTag(file); err == nil {
		s.response.Header().Add("ETag", etag)
	}

	fs := http.FileServer(http.Dir(s.app.Config.StaticDir))
	fs.ServeHTTP(s.response, s.request)
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
