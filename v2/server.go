package kwiscale

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"kwiscale/db"
	"net/http"
)

// A Server handles routes and can respond to an address.
type Server struct {
	// mux.Router that handles routes.
	Router *mux.Router

	// Number of handler to stack (default: 10).
	HandlerStackSize int

	// Template directory (default: ./).
	TemplateDir string

	// Gorilla session store interface
	SessionStore sessions.Store

	// ORM Database connection
	DB *gorm.DB
}

// Constructor for new Server.
func NewServer() *Server {
	return &Server{
		HandlerStackSize: 10,
		Router:           mux.NewRouter(),
		TemplateDir:      "./",
	}
}

// Add a route (mux.Route) to respond.
// Factory is a simple closure that returns a IHandler able to respond to the route.
func (s *Server) Route(route string, factory Factory) {
	stack := stacker(factory, s.HandlerStackSize)
	s.Router.HandleFunc(route, func(w http.ResponseWriter, req *http.Request) {
		// get handler from stack and dispatch request
		h := <-stack
		h.setServer(s)
		dispatch(h, w, req)
	})
}

// Listen starts to lisen on address.
func (s *Server) Listen(address string) {
	http.ListenAndServe(address, s.Router)
}

// InitDB initialize database
func (s *Server) InitDB(driver, conn string) {
	dbconn, err := db.InitDB(driver, conn)
	if err != nil {
		panic(err)
	}
	s.DB = dbconn
}
