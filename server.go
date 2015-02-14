package main
 
import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"runtime"
)
 
var firebaseurl string = "https://go-mine.firebaseio.com/"
var secret string = "tqOsGYhixWNyORaiO0g8AOcXEdO6JzNbPQhbJHNT"
 
var mcStdIn chan string
 
func init() {
	mcStdIn = make(chan string)
}
 
func main() {
    runtime.GOMAXPROCS(2)
	os.Chdir("server")
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft_server.jar", "nogui")
	stdoutPipe, _ := command.StdoutPipe()
	stdinPipe, _ := command.StdinPipe()
	_ = command.Start()
	db := DB{firebaseurl, secret}
	go stream(stdoutPipe, db)
	os.Chdir("..")
	loadConfig()
 
	// A for will block so put it in a goroutine
	go func() {
		for {
			select {
			// Read from our channel then write it to the servers stdin
			case cmd := <-mcStdIn:
				io.WriteString(stdinPipe, cmd+"\n")
			}
			// TODO add a handle server shutdown
		}
	}()
	defer command.Wait()
	httpServer()
}
 
func stream(stdoutPipe io.ReadCloser, db DB) {
	log.Println("Server started.")
	rd := bufio.NewReader(stdoutPipe)
	for {
		str, err := rd.ReadString('\n')
		if err != nil {
			log.Fatal("Read Error:", err)
		}
		
		db.message(str)
	}
}
 
func httpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/command/{command}/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		command := mux.Vars(r)["command"]
		db := DB{firebaseurl, secret}
		auth := db.check(token)
		if auth == 0 {
			log.Println("sendCommand: Invalid token.")
		} else {
			mcStdIn <- command
		}
	})
	r.HandleFunc("/configs/{token}", getConfigs)
	r.HandleFunc("/config/{id}/{token}", getConfig)
	r.HandleFunc("/update/{id}/{content}/{token}", setConfig)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	log.Println("HTTP server started on :9000")
	go http.ListenAndServe(":9000", nil)
}
 
func getConfigs(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	db := DB{firebaseurl, secret}
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
 
func getConfig(w http.ResponseWriter, r *http.Request) {
	var i int
	token := mux.Vars(r)["token"]
	id := mux.Vars(r)["id"]
	db := DB{firebaseurl, secret}
	auth := db.check(token)
	if auth == 0 {
		log.Println("getConfig: Invalid token.")
	} else {
		i, _ = strconv.Atoi(id)
		file := files[i]
		var buf, err = ioutil.ReadFile(file)
		if err != nil {
			log.Println(err)
		}
 
		splitstring := strings.SplitAfter(file, "/")
		filename := splitstring[len(splitstring)-1]
 
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(buf)
	}
}
 
func setConfig(w http.ResponseWriter, r *http.Request) {
	var i int
	token := mux.Vars(r)["token"]
	id := mux.Vars(r)["id"]
	content := mux.Vars(r)["content"]
	db := DB{firebaseurl, secret}
	auth := db.check(token)
	if auth == 0 {
		log.Println("setConfig: Invalid token.")
	} else {
		i, _ = strconv.Atoi(id)
		file := files[i]
		err := ioutil.WriteFile(file, []byte(content), 0644)
		if err != nil {
			log.Println(err)
			w.Write([]byte("There was an error updating the file."))
		}
		w.Write([]byte("The file was updated successfully."))
	}
}