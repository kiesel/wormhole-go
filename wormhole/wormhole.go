package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/kiesel/wormhole-go/lib"

	"gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v2"
)

type Error interface {
	Error() string
}

var log = logging.MustGetLogger("wormhole")
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{level:.1s} %{id:03x}%{color:reset} >> %{message}",
)
var config wormhole.WormholeConfig
var VersionString string

func Version() string {
	if "" != VersionString {
		return VersionString
	}

	return "<snapshot>"
}

func readConfiguration() (err Error) {
	var newConfig wormhole.WormholeConfig

	log.Info("Trying to parse wormhole configuration from " + wormhole.GetDefaultConfig())

	source, err := ioutil.ReadFile(wormhole.GetDefaultConfig())
	if err != nil {
		log.Critical(err.Error())

		return err
	}

	err = yaml.Unmarshal(source, &newConfig)
	if err != nil {
		log.Critical(err.Error())
		return err
	}

	// Now replace existing config with new
	log.Debug("New configuration %v", newConfig)
	config = newConfig

	return nil
}

func main() {

	// Setup logging
	logbackend := logging.NewLogBackend(os.Stdout, "", 0)
	logbackendformatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(logbackendformatter)

	// Read config
	readConfiguration()

	// Start main
	log.Info("Wormhole %s server starting ...", Version())

	l, err := net.Listen("tcp4", config.GetAddr())
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Listening at " + l.Addr().String())

	defer l.Close()
	for {
		// Wait for connection
		conn, err := l.Accept()

		if err != nil {
			log.Critical(err.Error())
			continue
		}

		// Handle connection
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	log.Debug("Received connection from %s", c.RemoteAddr().String())

	defer c.Close()
	defer log.Debug("Closed connection to %s", c.RemoteAddr().String())

	line, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		log.Critical(err.Error())
		return
	}

	writer := bufio.NewWriter(c)

	log.Debug("[%s] >> %s", c.RemoteAddr().String(), line)
	resp, err := handleLine(c, line)

	if err != nil {
		log.Warning("[%s] << %s", c.RemoteAddr().String(), err.Error())

		writer.WriteString("[ERR] ")
		writer.WriteString(err.Error())
	} else {
		log.Info("[%s] << %s", c.RemoteAddr().String(), resp)

		writer.WriteString("[OK] ")
		writer.WriteString(resp)
	}

	writer.WriteString("\n")
	writer.Flush()
}

func handleLine(c net.Conn, line string) (resp string, err Error) {
	parts := strings.Split(strings.TrimSpace(line), " ")

	log.Debug("Extracted parts %s", parts)
	if len(parts) < 1 {
		log.Warning("Too few parts.")
		return "", errors.New("Too few words, expected at least 1.")
	}

	switch strings.ToLower(parts[0]) {
	case "invoke":
		if len(parts) < 2 {
			log.Warning("Too few parts.")
			return "", errors.New("Too few words, expected at least 2.")
		}

		return handleInvocation(parts[1], parts[2:])
	case "exit":
		return handleExit()
	case "reload":
		return handleReload()
	}

	return "", errors.New("Unknown command, expected one of " + config.AvailableApps())
}

func handleInvocation(mapping string, args []string) (resp string, err Error) {
	app, err := config.GetAppWith(mapping, args)
	if err != nil {
		return "", err
	}

	log.Info("Invoking '%v' (mapped by %v) with args: %v", app.Executable, mapping, args)
	go wormhole.ExecuteCommand(app.Executable, args...)

	return "Started " + mapping, nil
}

func handleExit() (response string, err Error) {
	log.Warning("Client requested exit, exiting.")

	os.Exit(0)
	return "Bye!", nil
}

func handleReload() (response string, err Error) {
	if err := readConfiguration(); err != nil {
		return "", err
	}

	return "Re-read configuration.", nil
}
