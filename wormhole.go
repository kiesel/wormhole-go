package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"

	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("wormhole")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} %{level:.5s} %{id:03x}%{color:reset} >> %{message}",
)

type Error interface {
	Error() string
}

type WormholeConfig struct {
	Port    int               `yaml:"port,omitempty"`
	Mapping map[string]string `yaml:"mapping"`
	Editors map[string]string `yaml:"editors"`
}

func (this *WormholeConfig) String() string {
	str := "Config {\n"
	str += "  Port: " + fmt.Sprint(this.GetPort()) + "\n"

	str += "  Path mappings: {\n"
	for key, value := range this.Mapping {
		str += "    " + key + " -> " + value + "\n"
	}
	str += "  }\n"

	str += "  Editors: {\n"
	for name, path := range this.Editors {
		str += "    " + name + " -> " + path + "\n"
	}
	str += "  }\n"

	str += "}"

	return str
}

func (this *WormholeConfig) GetPort() int {
	if 0 == this.Port {
		return 5115
	}

	return this.Port
}

func main() {

	// Setup logging
	logbackend := logging.NewLogBackend(os.Stdout, "", 0)
	logbackendformatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(logbackend, logbackendformatter)

	// Read config
	log.Info("Parsing wormhole configuration ...")
	var config WormholeConfig
	source, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".wormhole.yml"))
	log.Debug("%s", source)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Configuration: %v", config.String())

	// Start main
	log.Info("Wormhole server starting ...")

	l, err := net.Listen("tcp4", ":"+strconv.Itoa(config.GetPort()))
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Listening at " + l.Addr().String())

	defer l.Close()
	for {
		// Wait for connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("Received connection from %s", conn.RemoteAddr().String())

		// Handle connection
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	line, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(c)

	log.Debug("[%s] %s", c.RemoteAddr().String(), line)
	resp, err := handleLine(c, line)

	if err != nil {
		writer.WriteString("Err ")
		writer.WriteString(err.Error())
	} else {
		writer.WriteString("Ok ")
		writer.WriteString(resp)
	}

	writer.Flush()
}

func handleLine(c net.Conn, line string) (resp string, err Error) {
	parts := strings.Split(strings.TrimSpace(line), " ")

	log.Debug("Extracted parts %s", parts)
	if len(parts) < 2 {
		log.Debug("Too little parts, quit.")
		return "", errors.New("Too few words, expected at least 2.")
	}

	switch parts[0] {
	case "EDIT":
		return handleCommandEdit(parts[1:])

	case "SHELL":
		return handleCommandShell(parts[1:])

	case "EXPLORE":
		return handleCommandExplore(parts[1:])

	case "START":
		return handleCommandStart(parts[1:])
	}

	return "", errors.New("Unknown command, expected one of [EDIT, SHELL, EXPLORE, START]")
}

func handleCommandEdit(parts []string) (resp string, err Error) {
	log.Info("EDIT", parts)
	return "OK", nil
}

func handleCommandStart(parts []string) (resp string, err Error) {
	log.Info("START", parts)
	return "OK", nil
}

func handleCommandExplore(parts []string) (resp string, err Error) {
	log.Info("EXPLORE", parts)
	return "OK", nil
}

func handleCommandShell(parts []string) (resp string, err Error) {
	log.Info("SHELL", parts)
	return "OK", nil
}
