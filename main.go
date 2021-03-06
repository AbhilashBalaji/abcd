package main

import (
	"flag"
	"log"

	"github.com/AbhilashBalaji/abcd/config"
	"github.com/AbhilashBalaji/abcd/db"
	"github.com/AbhilashBalaji/abcd/replication"
	"github.com/AbhilashBalaji/abcd/web"

	http "net/http"
)

var (
	dbLocation = flag.String("db-location", "", "Bolt DB path")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "config file for static sharding")
	shard      = flag.String("shard", "", "Name of shard for the data")
	replica    = flag.Bool("replica", false, "run as master or replica")
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

	c, err := config.ParseFile(*configFile)
	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config %q: %v", *configFile, err)
	}

	log.Printf("Shard count is %d , current shard : %d", shards.Count, shards.CurIdx)

	db, close, err := db.NewDatabase(*dbLocation, *replica)

	if err != nil {
		log.Fatalf("error creating NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()

	if *replica {
		leaderAddr, ok := shards.Addrs[shards.CurIdx]
		if !ok {
			log.Fatalf("Could not find address for leader for shard %d", shards.CurIdx)
		}
		go replication.ClientLoop(db, leaderAddr)
	}

	srv := web.NewServer(db, shards)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeysHandler)
	http.HandleFunc("/next-replication-key", srv.GetNextKeyForReplication)
	http.HandleFunc("/delete-replication-key", srv.DeleteReplicationKey)

	// hash(key) % <count> := shardIdx

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
