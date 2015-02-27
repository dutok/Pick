package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/stretchr/graceful"
	"github.com/unrolled/secure"
	"time"
	"net/http"
	"os"
	"html/template"
)

var at string
var err error

type Home struct {
	AccessToken string
}

func main() {
	loadConfig()
	server := newServer()
	startServer(&server)
	httpServer(&server)
}

func httpServer(server *Server) {
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"dutok.koding.io", "udkk5833aa60.dutok.koding.io"},
		SSLRedirect:           false,
		SSLHost:               "ssl.dutok.koding.io",
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
	})

    at = os.Args[1]
    
	secureMux := mux.NewRouter()
    s := secureMux.PathPrefix("/" + at).Subrouter()
    
	loadRoutes(s, server)
	
	secureMux.HandleFunc("/" + at, rootHandler)
	secureMux.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("public"))))

	secure := negroni.New()
	secure.UseHandler(secureMux)
	secure.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))

	n := negroni.New()
	n.Use(gzip.Gzip(gzip.BestSpeed))

	n.UseHandler(secureMux)

	graceful.Run(":80", 10*time.Second, n)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("public/index.html")
	home := Home{
		AccessToken: at,
	}
	t.Execute(w, home)
}
