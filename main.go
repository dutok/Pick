package main

import (
	//*"encoding/json"
	"github.com/codegangsta/negroni"
    oauth2 "github.com/goincremental/negroni-oauth2"
    sessions "github.com/goincremental/negroni-sessions"
    "github.com/goincremental/negroni-sessions/cookiestore"
	"net/http"
	/*"io/ioutil"
	"net/http"
	"strconv"
	"strings"*/
)

var err error

func main() {
	loadConfig()
	server := newServer()
	startServer(&server)
	httpServer(&server)
}

func httpServer(server *Server) {
	secureMux := http.NewServeMux()
    
    secureMux.HandleFunc("/sock", func(w http.ResponseWriter, r *http.Request) {
        sockServer(server, server.Messages, w, r)
    })
    
    secureMux.Handle("/", http.FileServer(http.Dir("public")))

    secure := negroni.New()
    secure.Use(oauth2.LoginRequired())
    secure.UseHandler(secureMux)

    n := negroni.New()
    n.Use(sessions.Sessions("my_session", cookiestore.New([]byte("secret123"))))
    n.Use(oauth2.Google(&oauth2.Config{
        ClientID:     "824000373870-148afj3scuj2fururtrn2ffn9vu48rfs.apps.googleusercontent.com",
        ClientSecret: "tB9cqq53V1H0yXjsp1SGKcDv",
        RedirectURL:  "http://dutok.koding.io/oauth2callback",
        Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
    }))

    router := http.NewServeMux()

    //There is probably a nicer way to handle this than repeat the restricted routes again
    //of course, you could use something like gorilla/mux and define prefix / regex etc.
    router.Handle("/", secure)
    router.Handle("/sock", secure)

    n.UseHandler(router)

    n.Run(":80")
	
    /*
	r.HandleFunc("/configs/{token}", func(w http.ResponseWriter, r *http.Request) {
		getConfigs(w, r, db)
	})
	r.HandleFunc("/config/{id}/{token}", func(w http.ResponseWriter, r *http.Request) {
		getConfig(w, r, db)
	})
	r.HandleFunc("/update/{id}/{content}/{token}", func(w http.ResponseWriter, r *http.Request) {
		setConfig(w, r, db)
	})
	r.HandleFunc("/server/stop/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		auth := db.check(token)
		if auth == 0 {
			log.Println("Auth: stop - Invalid token.")
		} else {
			server.stop()
		}
	})
	r.HandleFunc("/server/start/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		auth := db.check(token)
		if auth == 0 {
			log.Println("Auth: start - Invalid token.")
		} else {
			server = newServer(db)
			startServer(&server)
		}
	})
	r.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		getStats(w, r, &server)
	})
	
	http.Handle("/", r)
	log.Println("HTTP server: STARTED on :9000")
	http.ListenAndServe(":9000", nil)*/
}

/*func getStats(w http.ResponseWriter, r *http.Request, server *Server) {
	statsjson, err := json.Marshal(&server.Stats)
	check(err, "Minecraft stats")

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsjson)
}

func getConfigs(w http.ResponseWriter, r *http.Request, db DB) {
	token := mux.Vars(r)["token"]
	auth := db.check(token)
	if auth == 0 {
		log.Println("Auth: getConfigs - Invalid token.")
	} else {
		configjson, err := json.Marshal(files)
		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(configjson)
	}
}

func getConfig(w http.ResponseWriter, r *http.Request, db DB) {
	var i int
	token := mux.Vars(r)["token"]
	id := mux.Vars(r)["id"]
	auth := db.check(token)
	if auth == 0 {
		log.Println("Auth: getConfig - Invalid token.")
	} else {
		i, err = strconv.Atoi(id)
		check(err, "HTTP server")
		file := files[i]
		var buf, err = ioutil.ReadFile(file)
		check(err, "HTTP server")

		splitstring := strings.SplitAfter(file, "/")
		filename := splitstring[len(splitstring)-1]

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(buf)
	}
}

func setConfig(w http.ResponseWriter, r *http.Request, db DB) {
	var i int
	token := mux.Vars(r)["token"]
	id := mux.Vars(r)["id"]
	content := mux.Vars(r)["content"]
	auth := db.check(token)
	if auth == 0 {
		log.Println("Auth: setConfig - Invalid token.")
	} else {
		i, _ = strconv.Atoi(id)
		file := files[i]
		newcontent := strings.Replace(content, "&#47;", "/", -1)
		err := ioutil.WriteFile(file, []byte(newcontent), 0644)
		check(err, "HTTP server")
		w.Write([]byte("The file was updated successfully."))
	}
} */