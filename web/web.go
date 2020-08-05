package web

import (
	"fmt"
	"hash/fnv"
	"net/http"

	"../db"
)

// Server Contains HTTP handlers for DB
type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
}

// NewServer creates a new Server instance w/ HTTP handlers to get/set lmao
func NewServer(db *db.Database, shardIdx, shardCount int) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
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

	h := fnv.New64()
	h.Write([]byte(key))
	shardIdx := int(h.Sum64() % uint64(s.shardCount))

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v , hash = %d , shardIDx = %d\n", err, h.Sum64(), shardIdx)

}
