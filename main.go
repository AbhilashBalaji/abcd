package main

import (
	"flag"
	"log"

	"github.com/AbhilashBalaji/abcd/config"
	"github.com/AbhilashBalaji/abcd/db"
	"github.com/AbhilashBalaji/abcd/web"
	"github.com/BurntSushi/toml"

	http "net/http"
)

var (
	dbLocation = flag.String("db-location", "", "Bolt DB path")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "config file for static sharding")
	shard      = flag.String("shard", "", "Name of shard for the data")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must Provide DB location")
	}
	if *shard == "" {
		log.Fatal("Must Provide Shard")
	}
}

func main() {
	parseFlags()

	var c config.Config

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.DecodeFile(%q): %v", *configFile, err)
	}

	var shardCount int
	var shardIdx int = -1
	var addrs = make(map[int]string)

	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address

		if s.Idx+1 > shardCount {
			shardCount = s.Idx + 1
		}
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Shard %q was not found", *shard)
	}

	log.Printf("Shard count is %d , current shard : %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbLocation)

	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()
	srv := web.NewServer(db, shardIdx, shardCount, addrs)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	// hash(key) % <count> := shardIdx

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
