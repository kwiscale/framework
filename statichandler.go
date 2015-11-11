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
	abs, _ := filepath.Abs(s.app.Config.StaticDir)
	file = filepath.Join(abs, file)

	// control or add etag
	if etag, err := eTag(file); err == nil {
		s.response.Header().Add("ETag", etag)
	}

	// create a fileserver for the static dir
	fs := http.FileServer(http.Dir(s.app.Config.StaticDir))
	// stip directory name and serve the file
	http.StripPrefix("/"+filepath.Base(s.app.Config.StaticDir), fs).
		ServeHTTP(s.Response(), s.Request())
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
