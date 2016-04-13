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

func main() {

	// Setup logging
	logbackend := logging.NewLogBackend(os.Stdout, "", 0)
	logbackendformatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(logbackendformatter)

	// Read config
	log.Info("Trying to parse wormhole configuration from " + wormhole.GetDefaultConfig())
	source, err := ioutil.ReadFile(wormhole.GetDefaultConfig())
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
		writer.WriteString("[OK] ")
		writer.WriteString(resp)
	}

	writer.WriteString("\n")
	writer.Flush()
}

func handleLine(c net.Conn, line string) (resp string, err Error) {
	parts := strings.Split(strings.TrimSpace(line), " ")

	log.Debug("Extracted parts %s", parts)
	if len(parts) < 2 {
		log.Warning("Too little parts, quit.")
		return "", errors.New("Too few words, expected at least 2.")
	}

	switch strings.ToLower(parts[0]) {
	case "invoke":
		return handleInvocation(parts[1], parts[2:])
	case "exit":
		return handleExit()
	case "reload":
		return handleReload()
	}

	return "", errors.New("Unknown command, expected one of " + config.AvailableApps())
}

func handleInvocation(mapping string, args []string) (resp string, err Error) {
	app, err := config.GetApp(mapping)
	if err != nil {
		return "", err
	}

	args = app.MergeArguments(args)
	for index, arg := range args {
		args[index] = config.TranslatePath(arg)
	}

	log.Info("Invoking '%v' (mapped by %v) with args: %v", app.Executable, mapping, args)
	go wormhole.ExecuteCommand(app.Executable, args...)

	return "Started " + mapping, nil
}

func handleExit() (response string, err Error) {
	return "Bye!", nil
}

func handleReload() (response string, err Error) {
	return "Reloading...", nil
}
