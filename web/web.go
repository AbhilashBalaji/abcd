package web

import (
	"fmt"
	"net/http"

	"../db"
)

// Server Contains HTTP handlers for DB
type Server struct {
	db *db.Database
}

// NewServer creates a new Server instance w/ HTTP handlers to get/set lmao
func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

// GetHandler handles "GET" endpoint
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Value = %q , error = %v", value, err)

}

// SetHandler handles "WRITE" endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v", err)

}
