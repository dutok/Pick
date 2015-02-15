package main

import (
    "bufio"
    "os/exec"
    "io"
    "log"
    "strings"
    "github.com/lukevers/mcgoquery"
    "github.com/cloudfoundry/gosigar"
    "time"
)

type Server struct {
    Host      *string
	QueryPort *int
    Command *exec.Cmd
    Db DB
    Stdoutpipe *io.ReadCloser
    Stdinpipe *io.WriteCloser
    Status *int
    Cmdchan chan string
    Query     *mcgoquery.Client
    Stats *Stats
}

type MinecraftStats struct {
    Status     *int
    GameType   string
	GameId     string
	Version    string
	Map        string
	MaxPlayers int
	NumPlayers int
	Motd       string
}

type Stats struct {
    MinecraftStats
    ServerStats
}

type ServerStats struct {
   Memory sigar.Mem
}

func newServer(db DB) Server {
	command := exec.Command("java", "-Xmx512M", "-Xms512M", "-jar", "minecraft_server.jar", "nogui")
	command.Dir = "server"
	stdoutPipe, err := command.StdoutPipe()
	check(err, "Minecraft server")
	stdinPipe, err := command.StdinPipe()
	check(err, "Minecraft server")
	cmdchan := make(chan string)
	host := "localhost"
	queryport := 25565
	var c mcgoquery.Client
	var stats Stats
	status := 0
	server := Server{&host, &queryport, command, db, &stdoutPipe, &stdinPipe, &status, cmdchan, &c, &stats}
	return server
}

func startServer(server *Server) {
    if server.Command.Process != nil {
        server.Command.Process.Kill()
    }
    server.Command.Start()
	log.Println("Minecraft server: STARTED")
	*server.Status = 1
	
	go func() {
    	for {
        		select {
            		case cmd := <-server.Cmdchan:
            		    if *server.Status == 1 {
            			    io.WriteString(*server.Stdinpipe, cmd+"\n")
            		    }
        		}
    	}
	}()
	
	go func() {
    	rd := bufio.NewReader(*server.Stdoutpipe)
    	for {
            str, _ := rd.ReadString('\n')
            if str != "" {
                if strings.Contains(str, "Saving chunks for level") {
                    server.Db.message(str)
                    server.Command.Process.Kill()
                    break
                } else if strings.Contains(str, "Query running") {
                    server.Db.message(str)
                    server.Query, err = mcgoquery.Create(*server.Host, *server.QueryPort)
                    if err == nil {
                        updateStats(server)
                        go queryTimer(server)
                    } else {
                        log.Println(err)
                    }
                } else {
        	        server.Db.message(str)
                }
            }
    	}
    	server.Command.Wait()
    	server.Command.Process.Release()
    	log.Println("Minecraft server: STOPPED")
    	*server.Status = 0
	}()
}

func queryTimer(server *Server) {
    log.Println("Minecraft query: connected")
    ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				updateStats(server)
			}
		}
	}()
	time.Sleep(1 * time.Nanosecond)
}

func updateStats(server *Server) {
    server.Stats.ServerStats.Memory.Get()
    stat, err := server.Query.Full()
    
    if stat != nil {
        check(err, "Minecraft query")
    }
    
    mem := sigar.Mem{}
	mem.Get()
	
	var mcstats MinecraftStats

    if stat != nil {
        mcstats = MinecraftStats{
            server.Status,
            stat.GameType,
    		stat.GameID,
    		stat.Version,
    		stat.Map,
    		stat.MaxPlayers,
    		stat.NumPlayers,
    		stat.MOTD,
        }
    }
    ServerStats := ServerStats{
        mem,  
    }
    *server.Stats = Stats{
        mcstats,
        ServerStats,
    }
}

func (server *Server) stop() {
    server.sendCommand("stop")
}

func (server *Server) status() int {
    return *server.Status  
}

func (server *Server) sendCommand(command string) {
    log.Println("command recieved: " + command)
	server.Cmdchan <- command
}

func (server *Server) connect() {
	var err error

	// Try to connect to the Minecraft server Query
	server.Query, err = mcgoquery.Create(*server.Host, *server.QueryPort)
	if err != nil {
		// Try reconnecting in 15 seconds
	}
}