package main

import (
    "github.com/melvinmt/firebase"
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