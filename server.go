package main

import (
    "os/exec"
    "os"
	"io"
	"bufio"
	"log"
	"time"
)

var firebaseurl string = "https://go-mine.firebaseio.com/"
var secret string = "tqOsGYhixWNyORaiO0g8AOcXEdO6JzNbPQhbJHNT"

func main() {
    os.Chdir("server")
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft_server.jar", "nogui")
	stdoutPipe, _ := command.StdoutPipe()
    _ = command.Start()
    db := DB{firebaseurl, secret}
    go stream(stdoutPipe, db)
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