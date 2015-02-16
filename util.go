package main

import (
	"log"
)

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
