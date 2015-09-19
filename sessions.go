package kwiscale

import (
	"errors"
	"log"

	"github.com/gorilla/sessions"
)

var sessionEngine = make(map[string]ISessionStore, 0)

type SessionEngineOptions map[string]interface{}

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

	// SetOptions set some optionnal values to session engine
	SetOptions(*SessionEngineOptions)

	// SetSecret should register a string to encode cookie (not mandatory
	// but you should implement this to respect interface)
	SetSecret([]byte)

	// Get a value from storage , interface param is the key
	Get(IBaseHandler, interface{}) (interface{}, error)

	// Set a value in the storage, first interface param is the key,
	// second interface is the value to store
	Set(IBaseHandler, interface{}, interface{})

	// Clean, should cleanup files
	Clean(IBaseHandler)
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

// SetOptions does nothing for the engine
func (*SessionStore) SetOptions(*SessionEngineOptions) {}

// Get a value from session by name.
func (s *SessionStore) Get(handler IBaseHandler, key interface{}) (interface{}, error) {
	session, err := s.store.Get(handler.getRequest(), s.name)
	if err != nil {
		return nil, err
	}
	log.Print("Getting session", key, session.Values[key])
	if session.Values[key] == nil {
		return nil, errors.New("empty session")
	}
	return session.Values[key], nil
}

// Set a named value in sessionstore.
func (s *SessionStore) Set(handler IBaseHandler, key interface{}, val interface{}) {
	log.Print("Writing session", key, val)
	session, _ := s.store.Get(handler.getRequest(), s.name)
	session.Values[key] = val
	session.Save(handler.getRequest(), handler.getResponse())
}

func (s *SessionStore) Clean(handler IBaseHandler) {

	session, _ := s.store.Get(handler.getRequest(), s.name)
	session.Values = make(map[interface{}]interface{})
	session.Save(handler.getRequest(), handler.getResponse())
}
