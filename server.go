package main

import (
    "os/exec"
    "os"
	"io"
	"bufio"
	"log"
	"time"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

var firebaseurl string = "https://go-mine.firebaseio.com/"
var secret string = "tqOsGYhixWNyORaiO0g8AOcXEdO6JzNbPQhbJHNT"

func main() {
    os.Chdir("server")
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft_server.jar", "nogui")
	stdoutPipe, _ := command.StdoutPipe()
	stdinPipe, _ := command.StdinPipe()
    _ = command.Start()
    db := DB{firebaseurl, secret}
    go stream(stdoutPipe, db)
    os.Chdir("..")
    go loadConfig()
    httpServer(stdinPipe)
    defer command.Wait()
}


func stream(stdoutPipe io.ReadCloser, db DB) {
    log.Println("Server started.")
    rd := bufio.NewReader(stdoutPipe)
    for {
	    str, err := rd.ReadString('\n')
	    if err != nil {
	        log.Fatal("Read Error:", err)
	        return
	    }
	    t := time.Now().Local()
	    db.message(str, t.Format("20060102150405"))
    }
}

func httpServer(stdinPipe io.WriteCloser) {
    r := mux.NewRouter()
    r.HandleFunc("/command/{command}/{token}", func(w http.ResponseWriter, r *http.Request) {
         token := mux.Vars(r)["token"]
         command := mux.Vars(r)["command"]
         db := DB{firebaseurl, secret}
         auth := db.check(token)
         if (auth == 0){
            log.Println("sendCommand: Invalid token.")
         } else {
            io.WriteString(stdinPipe, command + "\n")
         }
    })
    r.HandleFunc("/config/{token}", getConfigs)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
    http.Handle("/", r)
    log.Println("HTTP server started on :9000")
    http.ListenAndServe(":9000", nil)
}

func getConfigs(w http.ResponseWriter, r *http.Request) {
    token := mux.Vars(r)["token"]
    db := DB{firebaseurl, secret}
    auth := db.check(token)
    if (auth == 0){
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