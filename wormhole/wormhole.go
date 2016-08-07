package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/kiesel/wormhole-go/lib"
	"gopkg.in/op/go-logging.v1"
)

type Error interface {
	Error() string
}

var (
	// Compile time values
	VersionString string

	// Command line flags
	quiet          bool
	displayVersion bool
	configFilename string
	logFilename    string
	injectVia      string

	// Volatiles
	config wormhole.WormholeConfig
	log    = logging.MustGetLogger("wormhole")
	format = logging.MustStringFormatter(
		"%{color}%{time:15:04:05.000} %{level:.1s} %{id:03x}%{color:reset} >> %{message}",
	)
)

func init() {
	flag.BoolVar(&displayVersion, "version", false, "Show version number, then exit.")
	flag.BoolVar(&quiet, "quiet", false, "Enable quiet mode")
	flag.StringVar(&configFilename, "configfile", wormhole.GetDefaultConfig(), "Set configuration path")
	flag.StringVar(&logFilename, "log", wormhole.GetDefaultLog(), "Set log path")
	flag.StringVar(&injectVia, "inject", ":environment", "Inject wormhole via :environment or file")
}

func Version() string {
	if "" != VersionString {
		return VersionString
	}

	return "<snapshot>"
}

func readConfiguration() (err Error) {
	var newConfig *wormhole.WormholeConfig

	log.Info("Trying to parse wormhole configuration from " + configFilename)
	if newConfig, err = wormhole.ReadConfigurationFrom(configFilename); err != nil {
		return err
	}

	// Now replace existing config with new
	log.Debug("New configuration %v", newConfig)
	config = *newConfig

	return nil
}

func main() {
	flag.Parse()

	if displayVersion {
		fmt.Println("Wormhole Version " + Version())
		os.Exit(1)
	}

	// Setup logging
	var logbackend *logging.LogBackend

	if quiet {
		fmt.Println("Quiet mode, logging to " + logFilename)

		logfile, err := os.OpenFile(logFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		logbackend = logging.NewLogBackend(logfile, "", 0)
	} else {
		logbackend = logging.NewLogBackend(os.Stdout, "", 0)
	}

	logbackendformatter := logging.NewBackendFormatter(logbackend, format)
	logging.SetBackend(logbackendformatter)

	// Read config
	readConfiguration()

	// Start server
	log.Info("Wormhole %s server starting ...", Version())

	l, err := net.Listen("tcp4", config.GetAddr())
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	go listenOn(l)

	args := flag.Args()
	if len(args) == 0 {
		select{}
	} else {
		if err := runCommand(args, injectVia, l.Addr().(*net.TCPAddr)); err != nil {
			log.Fatal(err)
		}
	}
}

func runCommand(args []string, injectVia string, addr *net.TCPAddr) error {
	c := exec.Command(args[0], args[1:len(args)]...)
	log.Info("Wormhole command %s starting ...", c)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if injectVia == ":environment" {
		c.Env = os.Environ()
		c.Env = append(c.Env, fmt.Sprintf("WORMHOLE_PORT=%d", addr.Port))
		c.Env = append(c.Env, fmt.Sprintf("WORMHOLE_IP=%s", addr.IP))
	} else {
		injectFile, err := os.OpenFile(injectVia, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		injectFile.WriteString(fmt.Sprintf("WORMHOLE_PORT=%d\n", addr.Port))
		injectFile.WriteString(fmt.Sprintf("WORMHOLE_IP=%s\n", addr.IP))
		injectFile.Close()

		defer os.Remove(injectFile.Name())
	}

	if err := c.Start(); err != nil {
		return err
	}

	if err := c.Wait(); err != nil {
		return err
	}

	return nil
}

func listenOn(l net.Listener) {
	log.Info("Listening at " + l.Addr().String())

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
	case "version":
		return handleVersion()
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

	log.Info("Invoking '%v' (mapped by %v) with args: %v", app.Executable, mapping, app.Args)
	go wormhole.ExecuteCommand(app.Executable, app.Args...)

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

func handleVersion() (response string, err Error) {
	return Version(), nil
}
