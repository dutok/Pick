package main

import (
    "bufio"
    "os/exec"
    "io"
    "log"
    "strings"
)

type Server struct {
    Command *exec.Cmd
    Db DB
    Stdoutpipe io.ReadCloser
    Stdinpipe io.WriteCloser
    Status bool
}

var currentserver Server

func newServer(db DB) Server {
	command := exec.Command("java", "-Xmx512M", "-Xms512M", "-jar", "minecraft_server.jar", "nogui")
	command.Dir = "server"
	stdoutPipe, err := command.StdoutPipe()
	check(err, "Minecraft server")
	stdinPipe, err := command.StdinPipe()
	check(err, "Minecraft server")
	server := Server{command, db, stdoutPipe, stdinPipe, false}
	return server
}

func startServer(server Server) {
    server.Command.Start()
	log.Println("Minecraft server: STARTED")
	server.Status = true
	
	go func() {
    	for {
    		select {
        		// Read from our channel then write it to the servers stdin
        		case cmd := <-mcStdIn:
        			io.WriteString(server.Stdinpipe, cmd+"\n")
    		}
    	}
	}()
	
	go func() {
    	rd := bufio.NewReader(server.Stdoutpipe)
    	for {
            str, _ := rd.ReadString('\n')
            if str != "" {
                if strings.Contains(str, "Saving chunks for level") {
                    server.Db.message(str)
                    server.Command.Process.Kill()
                    break
                } else {
        	        server.Db.message(str)
                }
            }
    	}
    	server.Command.Wait()
    	server.Command.Process.Release()
    	log.Println("Minecraft server: STOPPED")
    	server.Status = false
	}()
}

func (server *Server) stop() {
    cmd := "stop"
    mcStdIn <- cmd
}

func (server *Server) status() bool {
    return server.Status  
}