package kwiscale

import (
	"github.com/gorilla/sessions"
)

type ISession interface {
	Get(name string)
	Set(name string, value interface{})
}

func NewCookieStore(secret string) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(secret))
}

func NewFileStore(secret string) *sessions.FilesystemStore {
	return sessions.NewFilesystemStore("/tmp/sessions", []byte(secret))
}

func (s *Server) SetSessionStore(session sessions.Store) {
	s.SessionStore = session
}

func (h *Handler) GetSession(key string) (*sessions.Session, error) {
	return h.Server.SessionStore.Get(h.Request, key)
}

func (h *Handler) SaveSession(session *sessions.Session) {
	session.Save(h.Request, h.Response.ResponseWriter)
}
