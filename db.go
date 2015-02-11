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
    body string
    url string
}

func (db *DB) write(message Message) {
    ref := firebase.NewReference(db.url + message.url).Auth(db.secret)

    // Write the value to Firebase.
    if err = ref.Write(message); err != nil {
        panic(err)
    }
}