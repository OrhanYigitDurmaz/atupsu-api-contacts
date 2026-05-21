//go:build windows

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"atupsu-api/api"
	"atupsu-api/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	dsn := flag.String("dsn", "", "ODBC data source name (e.g. atilim_64bit)")
	mdbPath := flag.String("mdb", `C:\atilim\atilim.mdb`, "path to atilim.mdb (ignored if -dsn is set)")
	addr := flag.String("addr", ":8080", "listen address")
	debug := flag.Bool("debug", false, "run interactively (skip service wrapper)")
	flag.Parse()

	var connStr string
	if *dsn != "" {
		connStr = db.ConnStringFromDSN(*dsn)
	} else {
		connStr = db.ConnStringFromFile(*mdbPath)
	}

	database, err := db.Open(connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	h := &api.Handler{DB: database}
	r := chi.NewRouter()
	r.Mount("/customers", h.Routes())

	server := &http.Server{Addr: *addr, Handler: r}

	interactive, err := isInteractive()
	if err != nil {
		log.Fatalf("failed to detect session type: %v", err)
	}

	if *debug || interactive {
		log.Printf("starting server on %s (interactive mode)", *addr)
		go server.ListenAndServe()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt)
		<-stop

		log.Println("shutting down...")
		server.Shutdown(context.Background())
	} else {
		if err := runService(server); err != nil {
			log.Fatalf("service failed: %v", err)
		}
	}
}
