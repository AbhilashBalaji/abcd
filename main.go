package main

import (
	"flag"
	"log"
	http "net/http"

	db "./db"
	web "./web"
)

var (
	dbLocation = flag.String("db-location", "", "Bolt DB path")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
)

func parseFlags() {
	flag.Parse()

	if *dbLocation == "" {
		log.Fatal("Must Provide DB location")
	}
}

func main() {
	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)

	if err != nil {
		log.Fatalf("NewDatabase(%q): %v", *dbLocation, err)
	}

	defer close()
	srv := web.NewServer(db)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
