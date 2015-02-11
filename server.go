package main

import (
	"os/exec"
	"io"
	"bufio"
	"log"
	"github.com/melvinmt/firebase"
	"net/http"
)

func main() {
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "server/minecraft_server.jar", "nogui")
	stdoutPipe, _ := command.StdoutPipe()
    _ = command.Start()
    go stream(stdoutPipe, messages)
    defer command.Wait()
    log.Println("Server started!!")
}

func stream(stdoutPipe io.ReadCloser, messages chan string) {
    rd := bufio.NewReader(stdoutPipe)
    for {
	    str, err := rd.ReadString('\n')
	    if err != nil {
	        log.Fatal("Read Error:", err)
	        return
	    }
	    // Do stuff with str
	}
}

func firebase(){
    var err error

    url := "https://go-mine.firebaseio.com/users/fred/name"

    // Can also be your Firebase secret:
    authToken := "tqOsGYhixWNyORaiO0g8AOcXEdO6JzNbPQhbJHNT"

    // Auth is optional:
    ref := firebase.NewReference(url).Auth(authToken)

    // Create the value.
    personName := PersonName{
        First: "Fred",
        Last:  "Swanson",
    }

    // Write the value to Firebase.
    if err = ref.Write(personName); err != nil {
        panic(err)
    }

    // Now, we're going to retrieve the person.
    personUrl := "https://go-mine.firebaseio.com/users/fred"

    personRef := firebase.NewReference(personUrl).Export(false)

    fred := Person{}

    if err = personRef.Value(fred); err != nil {
        panic(err)
    }

    fmt.Println(fred.Name.First, fred.Name.Last) // prints: Fred Swanson
}