package kwiscale

import (
	"github.com/gorilla/sessions"
)

var sessionEngine = make(map[string]ISessionStore, 0)

// RegisterSessionEngine can register session engine that implements
// ISessionStore. The name is used to let configuration to select it.
func RegisterSessionEngine(name string, engine ISessionStore) {
	sessionEngine[name] = engine
}

// ISessionStore to implement to give a session storage
type ISessionStore interface {
	// Init is called when store is initialized while App is initialized
	Init()

	// Name should set the session name
	Name(string)

	// SetSecret should register a string to encode cookie (not mandatory
	// but you should implement this to respect interface)
	SetSecret([]byte)

	// Get a value from storage , interface param is the key
	Get(IBaseHandler, interface{}) interface{}

	// Set a value in the storage, first interface param is the key,
	// second interface is the value to store
	Set(IBaseHandler, interface{}, interface{})
}

// SessionStore is a basic cookie based on gorilla.session.
type SessionStore struct {
	store  *sessions.CookieStore
	name   string
	secret []byte
}

// Init prepare the cookie storage.
func (s *SessionStore) Init() {
	s.store = sessions.NewCookieStore(s.secret)
}

// SetSecret record a string to encode cookie
func (s *SessionStore) SetSecret(secret []byte) {
	s.secret = secret
}

// Name set session name
func (s *SessionStore) Name(name string) {
	s.name = name
}

// Get a value from session by name.
func (s *SessionStore) Get(handler IBaseHandler, key interface{}) interface{} {
	session, _ := s.store.Get(handler.getRequest(), s.name)
	return session.Values[key]
}

// Set a named value in sessionstore.
func (s *SessionStore) Set(handler IBaseHandler, key interface{}, val interface{}) {
	session, _ := s.store.Get(handler.getRequest(), s.name)
	session.Values[key] = val
	session.Save(handler.getRequest(), handler.getResponse())
}
