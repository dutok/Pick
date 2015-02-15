package main

import (
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
    "time"
)

var lastMessage string

type DB struct {
    url string
    secret string
}

type Message struct {
    Body string
}

func (db *DB) message(msg string) {
    url := db.url + "console/messages.json?auth=" + secret
    message := Message{Body: msg}
    jsonmessage, err := json.Marshal(message)
    check(err, "Database")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonmessage))
    req.Close = true

    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        time.Sleep(time.Second)
        resp, err = client.Do(req)
        if err != nil {
            log.Println(err)
        }
    }
    
    defer resp.Body.Close()
}

func (db *DB) check(token string) int {
    response, err := http.Get(db.url + "allowed/" + token + ".json?auth=" + secret)
    defer response.Body.Close()
    if err != nil {
        log.Printf("%s", err)
    } else {
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            log.Printf("%s", err)
        }
        
        if string(contents) == "null" {
            return 0 //no go
        } else {
            return 1
        }
    }
    
    return 0
}