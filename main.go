package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
)

var firebaseurl string = "https://go-mine.firebaseio.com/"
var secret string = "IcVM9hRKCqz4GTkUpiHphbNBHg7y4hW62FJTM5bz"

var err error

func main() {
	runtime.GOMAXPROCS(2)
	db := DB{firebaseurl, secret}
	loadConfig()
	server := newServer(db)
	startServer(&server)
	httpServer(db, server)
}

func httpServer(db DB, server Server) {
	r := mux.NewRouter()
	r.HandleFunc("/command/{command}/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		command := mux.Vars(r)["command"]
		auth := db.check(token)
		if auth == 0 {
			log.Println("sendCommand: Invalid token.")
		} else {
			server.sendCommand(command)
		}
	})
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
			log.Println("stop: Invalid token.")
		} else {
			server.stop()
		}
	})
	r.HandleFunc("/server/start/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		auth := db.check(token)
		if auth == 0 {
			log.Println("start: Invalid token.")
		} else {
			server = newServer(db)
			startServer(&server)
		}
	})
	r.HandleFunc("/server", func(w http.ResponseWriter, r *http.Request) {
		getStats(w, r, &server)
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	log.Println("HTTP server started on :9000")
	http.ListenAndServe(":9000", nil)
}

func getStats(w http.ResponseWriter, r *http.Request, server *Server) {
	statsjson, err := json.Marshal(&server.Stats)
	check(err, "Minecraft stats")

	w.Header().Set("Content-Type", "application/json")
	w.Write(statsjson)
}

func getConfigs(w http.ResponseWriter, r *http.Request, db DB) {
	token := mux.Vars(r)["token"]
	auth := db.check(token)
	if auth == 0 {
		log.Println("getConfigs: Invalid token.")
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
		log.Println("getConfig: Invalid token.")
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
		log.Println("setConfig: Invalid token.")
	} else {
		i, _ = strconv.Atoi(id)
		file := files[i]
		newcontent := strings.Replace(content, "&#47;", "/", -1)
		err := ioutil.WriteFile(file, []byte(newcontent), 0644)
		check(err, "HTTP server")
		w.Write([]byte("The file was updated successfully."))
	}
}

func check(err error, source string) {
	if err != nil {
		log.Println("[" + source + "] " + err.Error())
	}
}

func fatalcheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
