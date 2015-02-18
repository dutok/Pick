package main

import (
	"github.com/gorilla/websocket"
	"net"
	"log"
	"sync"
	"net/http"
	
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn, messages *[]string) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
	for _, message := range *messages {
        if err := cc.websocket.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			return
		}
    }
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(message []byte, messages *[]string) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

    *messages = append(*messages, string(message))

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func sockServer(server *Server, messages *[]string, w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client}
	addClient(sockCli, messages)

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			deleteClient(sockCli)
			return
		}
		server.sendCommand(string(p))
	}
}