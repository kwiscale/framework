package kwiscale

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

// StaticHandler handle static files handlers. Use App.SetStatic(path) that create the static handler
type staticHandler struct {
	RequestHandler
	cacheEnabled bool
}

var cache = make(map[string][]byte)

func (s *staticHandler) putInCache(f string) {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		cache[f] = content
	}
}

func (s *staticHandler) Get() {
	file := s.Vars["file"]
	file = filepath.Join(s.app.Config.StaticDir, file)

	var (
		content []byte
		err     error
		exists  bool
	)

	if s.cacheEnabled {
		if content, exists = cache[file]; !exists {
			go s.putInCache(file)
		}
	}
	// not in cache or cache is disable
	if content == nil {
		if content, err = ioutil.ReadFile(file); err != nil {
			s.App().Error(http.StatusNotFound, s.getResponse(), err)
			return
		}
	}

	// control or add etag
	if etag, err := eTag(file); err == nil {
		if match, ok := s.Request.Header["If-None-Match"]; ok {
			for _, m := range match {
				if etag == m {
					s.Response.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
		s.Response.Header().Add("ETag", etag)
	}

	mimetype := mime.TypeByExtension(filepath.Ext(file))
	s.Response.Header().Add("Content-Type", mimetype)
	s.Response.Header().Add("Content-Length", fmt.Sprintf("%d", len(content)))
	s.Write(content)
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
