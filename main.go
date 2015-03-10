package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/stretchr/graceful"
	"time"
)

var err error

func main() {
	loadConfig()
	server := newServer()
	startServer(&server)
	httpServer(&server)
}

func httpServer(server *Server) {
	secureMux := mux.NewRouter()

	loadRoutes(secureMux, server)

	secure := negroni.New()
	secure.UseHandler(secureMux)

	n := negroni.New()
	n.Use(gzip.Gzip(gzip.BestSpeed))

	router := mux.NewRouter()

	//There is probably a nicer way to handle this than repeat the restricted routes again
	//of course, you could use something like gorilla/mux and define prefix / regex etc.
	router.PathPrefix("/").Handler(secure)

	n.UseHandler(router)

	graceful.Run(":3000", 10*time.Second, n)
}
