package main

import (
    "github.com/melvinmt/firebase"
    "log"
    "net/http"
    "os"
    "io/ioutil"
)

var err error

type DB struct {
    url string
    secret string
}

type Message struct {
    Body string
    Time string
}

func (db *DB) message(body string, time string) {
    ref := firebase.NewReference(db.url + "console/messages").Auth(db.secret)

    message := Message{
        Body: body,
        Time:  time,
    }

    // Write the value to Firebase.
    if err = ref.Push(message); err != nil {
        panic(err)
    }
}

func (db *DB) check(token string) int {
    response, err := http.Get(db.url + "allowed/" + token + ".json")
    if err != nil {
        log.Printf("%s", err)
        os.Exit(1)
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Printf("%s", err)
            os.Exit(1)
        }
        
        if string(contents) == "null" {
            return 0 //no go
        } else {
            return 1
        }
    }
    
    return 0
}