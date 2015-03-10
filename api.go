package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Token struct {
	Token string
}

func loadRoutes(secureMux *mux.Router, server *Server) {
	secureMux.HandleFunc("/sock", func(w http.ResponseWriter, r *http.Request) {
		sockServer(server, server.Messages, w, r)
	})

	secureMux.HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		getConfigs(w, r)
	})
	secureMux.HandleFunc("/config/{id}", func(w http.ResponseWriter, r *http.Request) {
		getConfig(w, r)
	})

	secureMux.HandleFunc("/update/{id}/{content}", func(w http.ResponseWriter, r *http.Request) {
		setConfig(w, r)
	})

	secureMux.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		getStats(w, r, server)
	})
	secureMux.HandleFunc("/server/start", func(w http.ResponseWriter, r *http.Request) {
		start(w, r, server)
	})

	secureMux.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
}

func getConfigs(w http.ResponseWriter, r *http.Request) {
	configjson, err := json.Marshal(files)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(configjson)
}

func getConfig(w http.ResponseWriter, r *http.Request) {
	var i int
	id := mux.Vars(r)["id"]
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

func getStats(w http.ResponseWriter, r *http.Request, server *Server) {
	statsjson, err := json.Marshal(&server.Stats)
	check(err, "Minecraft stats")

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsjson)
}

func setConfig(w http.ResponseWriter, r *http.Request) {
	var i int
	id := mux.Vars(r)["id"]
	content := mux.Vars(r)["content"]
	i, _ = strconv.Atoi(id)
	file := files[i]
	newcontent := strings.Replace(content, "&#47;", "/", -1)
	err := ioutil.WriteFile(file, []byte(newcontent), 0644)
	check(err, "HTTP server")
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("The file was updated successfully."))
}

func start(w http.ResponseWriter, r *http.Request, server *Server) {
	*server = newServer()
	startServer(server)
}
