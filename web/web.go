package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/AbhilashBalaji/abcd/config"

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

	shard := s.getShard(key)

	value, err := s.db.GetKey(key)
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Shard = %d, Current shard = %d,addr = %q Value = %q , error = %v", shard, s.shards.CurIdx, s.shards.Addrs[shard], value, err)

}

func (s *Server) getShard(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.shards.Count))
}

// SetHandler handles "WRITE" endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.getShard(key)
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d, current shard = %d", err, shard, s.shards.CurIdx)

}
