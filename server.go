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
var secret string = "IcVM9hRKCqz4GTkUpiHphbNBHg7y4hW62FJTM5bz"
 
var mcStdIn chan string
 
func init() {
	mcStdIn = make(chan string)
}

var err error
 
func main() {
    runtime.GOMAXPROCS(2)
	os.Chdir("server")
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft_server.jar", "nogui")
	stdoutPipe, err := command.StdoutPipe()
	check(err)
	stdinPipe, err := command.StdinPipe()
	check(err)
	err = command.Start()
	check(err)
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
	httpServer(db)
}
 
func stream(stdoutPipe io.ReadCloser, db DB) {
	log.Println("Server started.")
	rd := bufio.NewReader(stdoutPipe)
	for {
		str, err := rd.ReadString('\n')
		fatalcheck(err)
		
		db.message(str)
	}
}
 
func httpServer(db DB) {
	r := mux.NewRouter()
	r.HandleFunc("/command/{command}/{token}", func(w http.ResponseWriter, r *http.Request) {
		token := mux.Vars(r)["token"]
		command := mux.Vars(r)["command"]
		auth := db.check(token)
		if auth == 0 {
			log.Println("sendCommand: Invalid token.")
		} else {
			mcStdIn <- command
		}
	})
	r.HandleFunc("/configs/{token}", func(w http.ResponseWriter, r *http.Request){
	    getConfigs(w,r,db)   
	})
	r.HandleFunc("/config/{id}/{token}", func(w http.ResponseWriter, r *http.Request){
	    getConfig(w,r,db)   
	})
	r.HandleFunc("/update/{id}/{content}/{token}", func(w http.ResponseWriter, r *http.Request){
	    setConfig(w,r,db)   
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.Handle("/", r)
	log.Println("HTTP server started on :9000")
	go http.ListenAndServe(":9000", nil)
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
		check(err)
		file := files[i]
		var buf, err = ioutil.ReadFile(file)
		check(err)
 
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
		err := ioutil.WriteFile(file, []byte(content), 0644)
		check(err)
		w.Write([]byte("The file was updated successfully."))
	}
}

func check(err error) {
    if err != nil {
        log.Println(err)
    }
}

func fatalcheck(err error) {
    if err != nil {
        log.Fatal(err)
    }
}