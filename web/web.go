package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AbhilashBalaji/abcd/config"
	"github.com/AbhilashBalaji/abcd/replication"

	"github.com/AbhilashBalaji/abcd/db"
)

// Server Contains HTTP handlers for DB
type Server struct {
	db     *db.Database
	shards *config.Shards
}

// NewServer creates a new Server instance w/ HTTP handlers to get/set lmao
func NewServer(db *db.Database, s *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: s,
	}
}

// Redirect to correct shard
func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shards.CurIdx, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error Redirecting request : %v\n", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

// GetHandler handles "GET" endpoint
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	shard := s.shards.Index(key)

	value, err := s.db.GetKey(key)
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v", shard, s.shards.CurIdx, s.shards.Addrs[shard], value, err)

}

// SetHandler handles "WRITE" endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.shards.Index(key)
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d, current shard = %d", err, shard, s.shards.CurIdx)

}

// DeleteExtraKeysHandler handles "Purge" endpoint
func (s *Server) DeleteExtraKeysHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error = %v", s.db.DeleteExtraKeys(func(key string) bool {
		return s.shards.Index(key) != s.shards.CurIdx
	}))
}

// GetNextKeyForReplication returns the next key for replication.
func (s *Server) GetNextKeyForReplication(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	k, v, err := s.db.GetNextKeyForReplication()
	enc.Encode(&replication.NextKeyValue{
		Key:   string(k),
		Value: string(v),
		Err:   err,
	})
}

// DeleteReplicationKey deletes the key from replica queue.
func (s *Server) DeleteReplicationKey(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	err := s.db.DeleteReplicationKey([]byte(key), []byte(value))
	if err != nil {
		w.WriteHeader(http.StatusExpectationFailed)
		fmt.Fprintf(w, "error: %v", err)
		return
	}

	fmt.Fprintf(w, "ok")
}
