package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v2"
)

type Error interface {
	Error() string
}

type WormholeConfig struct {
	Port    int               `yaml:"port,omitempty"`
	Mapping map[string]string `yaml:"mapping"`
	Editors map[string]string `yaml:"editors"`
}

func (this *WormholeConfig) GetPort() int {
	if 0 == this.Port {
		return 5115
	}

	return this.Port
}

var log = logging.MustGetLogger("wormhole")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} %{level:.5s} %{id:03x}%{color:reset} >> %{message}",
)
var config WormholeConfig

func main() {

	// Setup logging
	logbackend := logging.NewLogBackend(os.Stdout, "", 0)
	logbackendformatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(logbackendformatter)

	// Read config
	log.Info("Parsing wormhole configuration ...")
	source, err := ioutil.ReadFile(path.Join(os.Getenv("HOME"), ".wormhole.yml"))
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("Configuration: %v", config)

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
		writer.WriteString("[ERR] ")
		writer.WriteString(err.Error())
	} else {
		writer.WriteString("[OK]")
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
	case "INVOKE":
		return handleInvocation(parts[1:])
	}

	return "", errors.New("Unknown command, expected one of [EDIT, SHELL, EXPLORE, START]")
}

func handleInvocation(parts []string) (resp string, err Error) {
	log.Info("Invoking ", parts)

	go executeCommand("/bin/sleep", "10")
	return "OK", nil
}

func executeCommand(executable string, args ...string) (err Error) {
	cmd := exec.Command(executable, args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// cmd.StdoutPipe().close()
	// cmd.StderrPipe().close()
	// cmd.StdinPipe().close()

	if err := cmd.Start(); err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Started '%s' w/ PID %d", executable, cmd.Process.Pid)

	cmd.Wait()

	log.Info("PID %d has quit.", cmd.Process.Pid)

	return nil
}
