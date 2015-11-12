package kwiscale

import (
	"errors"

	"github.com/gorilla/sessions"
)

var sessionEngine = make(map[string]SessionStore, 0)

// SessionEngineOptions set options for session engine.
type SessionEngineOptions map[string]interface{}

// RegisterSessionEngine can register session engine that implements
// ISessionStore. The name is used to let configuration to select it.
func RegisterSessionEngine(name string, engine SessionStore) {
	sessionEngine[name] = engine
}

// Register cookiesessionstore by default.
func init() {
	RegisterSessionEngine("default", &CookieSessionStore{})
}

// SessionStore to implement to give a session storage
type SessionStore interface {
	// Init is called when store is initialized while App is initialized
	Init()

	// Name should set the session name
	Name(string)

	// SetOptions set some optionnal values to session engine
	SetOptions(SessionEngineOptions)

	// SetSecret should register a string to encode cookie (not mandatory
	// but you should implement this to respect interface)
	SetSecret([]byte)

	// Get a value from storage , interface param is the key
	Get(WebHandler, interface{}) (interface{}, error)

	// Set a value in the storage, first interface param is the key,
	// second interface is the value to store
	Set(WebHandler, interface{}, interface{})

	// Clean, should cleanup files
	Clean(WebHandler)
}

// CookieSessionStore is a basic cookie based on gorilla.session.
type CookieSessionStore struct {
	store  *sessions.CookieStore
	name   string
	secret []byte
}

// Init prepare the cookie storage.
func (s *CookieSessionStore) Init() {
	s.store = sessions.NewCookieStore(s.secret)
}

// SetSecret record a string to encode cookie
func (s *CookieSessionStore) SetSecret(secret []byte) {
	s.secret = secret
}

// Name set session name
func (s *CookieSessionStore) Name(name string) {
	s.name = name
}

// SetOptions does nothing for the engine
func (*CookieSessionStore) SetOptions(SessionEngineOptions) {}

// Get a value from session by name.
func (s *CookieSessionStore) Get(handler WebHandler, key interface{}) (interface{}, error) {
	session, err := s.store.Get(handler.getRequest(), s.name)
	if err != nil {
		return nil, err
	}
	Log("Getting session", key, session.Values[key])
	if session.Values[key] == nil {
		return nil, errors.New("empty session")
	}
	return session.Values[key], nil
}

// Set a named value in sessionstore.
func (s *CookieSessionStore) Set(handler WebHandler, key interface{}, val interface{}) {
	Log("Writing session", key, val)
	session, _ := s.store.Get(handler.getRequest(), s.name)
	session.Values[key] = val
	session.Save(handler.getRequest(), handler.getResponse())
}

// Clean removes the entire session values for current session.
func (s *CookieSessionStore) Clean(handler WebHandler) {
	session, _ := s.store.Get(handler.getRequest(), s.name)
	session.Values = make(map[interface{}]interface{})
	session.Save(handler.getRequest(), handler.getResponse())
}
