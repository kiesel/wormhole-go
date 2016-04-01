package main

import (
  "os"
  "net"
  "errors"
  "bufio"
  "strings"

  "gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("wormhole")
var format =  logging.MustStringFormatter(
  "%{color}%{time:15:04:05.000} %{shortfunc} %{level:.5s} %{id:03x}%{color:reset} >> %{message}",
)

type Error interface {
  Error() string
}

func main() {

  // Setup logging
  logbackend := logging.NewLogBackend(os.Stdout, "", 0)
  logbackendformatter := logging.NewBackendFormatter(logbackend, format)
  logging.SetBackend(logbackend, logbackendformatter)


  log.Info("Wormhole server starting ...")

  l, err := net.Listen("tcp", ":2000")
  if err != nil {
    log.Fatal(err)
  }

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

func handleCommandStart(parts []string) (resp string, err Error)  {
  log.Info("START", parts)
  return "OK", nil
}

func handleCommandExplore(parts []string) (resp string, err Error)  {
  log.Info("EXPLORE", parts)
  return "OK", nil
}

func handleCommandShell(parts []string) (resp string, err Error)  {
  log.Info("SHELL", parts)
  return "OK", nil
}