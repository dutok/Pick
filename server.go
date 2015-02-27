package main

import (
	"bufio"
	"github.com/cloudfoundry/gosigar"
	"github.com/lukevers/mcgoquery"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	Host       *string
	QueryPort  *int
	Command    *exec.Cmd
	Stdoutpipe *io.ReadCloser
	Stdinpipe  *io.WriteCloser
	Status     *int
	Cmdchan    chan string
	Messages   *[]string
	Query      *mcgoquery.Client
	Stats      *Stats
	TPS        *string
	StartTime  *int32
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
	Tps        *string
	StartTime  *int32
}

type Stats struct {
	MinecraftStats
	ServerStats
}

type ServerStats struct {
	Memory sigar.Mem
	CPU    float64
}

func newServer() Server {
	command := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft_server.jar", "nogui")
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
	var messages []string
	var tps = ""
	var starttime = int32(0)
	server := Server{&host, &queryport, command, &stdoutPipe, &stdinPipe, &status, cmdchan, &messages, &c, &stats, &tps, &starttime}
	return server
}

func startServer(server *Server) {
	if server.Command.Process != nil {
		server.Command.Process.Kill()
	}
	server.Command.Start()
	enableQuery()
	log.Println("Minecraft server: Attempting to start")

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
					broadcastMessage([]byte(str), server.Messages)
					server.Command.Process.Kill()
					break
				} else if strings.Contains(str, "Query running") {
					broadcastMessage([]byte(str), server.Messages)
					server.Query, err = mcgoquery.Create(*server.Host, *server.QueryPort)
					if err == nil {
						go queryTimer(*server)
					} else {
						log.Println(err)
					}
				} else if strings.Contains(str, "Done") && strings.Contains(str, "For help, type") {
					broadcastMessage([]byte(str), server.Messages)
					*server.Status = 1
					log.Println("Minecraft server: Started")
					*server.StartTime = int32(time.Now().Unix())
				} else if strings.Contains(str, "TPS from last") {
					output := str[len(str)-5:]
					tps := strings.Trim(output, "\n")
					*server.TPS = tps
				} else {
					broadcastMessage([]byte(str), server.Messages)
				}
			}
		}
		server.Command.Wait()
		server.Command.Process.Release()
		log.Println("Minecraft server: Stopped")
		*server.Status = 0
	}()
}

func queryTimer(server Server) {
	updateStats(server)
	log.Println("Minecraft query: connected")
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				server.sendCommand("tps")
				updateStats(server)
			}
		}
	}()
	time.Sleep(1 * time.Nanosecond)
}

func updateStats(server Server) {
	server.Stats.ServerStats.Memory.Get()
	stat, err := server.Query.Full()

	if stat != nil {
		check(err, "Minecraft query")
	}

	mem := sigar.Mem{}
	mem.Get()

	idle0, total0 := getCPUSample()
	time.Sleep(3 * time.Second)
	idle1, total1 := getCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)
	cpu := 100 * (totalTicks - idleTicks) / totalTicks

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
			server.TPS,
			server.StartTime,
		}
	}
	ServerStats := ServerStats{
		mem,
		cpu,
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

func enableQuery() {
	read, err := ioutil.ReadFile("server/server.properties")
	check(err, "Minecraft query")
	if strings.Contains(string(read), "enable-query=false") {
		r := string(read)
		Value := strings.Replace(r, "enable-query=false", "enable-query=true", -1)
		err = ioutil.WriteFile("server/server.properties", []byte(Value), 0644)
		check(err, "Minecraft query")
	}
}

func getCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					log.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}
