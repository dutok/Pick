package main

import (
	"github.com/codegangsta/negroni"
	oauth2 "github.com/goincremental/negroni-oauth2"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
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
	//secure.Use(oauth2.LoginRequired())
	secure.UseHandler(secureMux)

	n := negroni.New()
	n.Use(sessions.Sessions("my_session", cookiestore.New([]byte("secret123"))))
	n.Use(oauth2.Google(&oauth2.Config{
		ClientID:     "824000373870-148afj3scuj2fururtrn2ffn9vu48rfs.apps.googleusercontent.com",
		ClientSecret: "tB9cqq53V1H0yXjsp1SGKcDv",
		RedirectURL:  "http://dutok.koding.io/oauth2callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	}))
	n.Use(gzip.Gzip(gzip.BestSpeed))

	router := mux.NewRouter()

	//There is probably a nicer way to handle this than repeat the restricted routes again
	//of course, you could use something like gorilla/mux and define prefix / regex etc.
	router.PathPrefix("/").Handler(secure)

	n.UseHandler(router)

	n.Run(":80")
}
