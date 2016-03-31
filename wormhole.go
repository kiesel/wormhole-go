package main

import (
  "log"
  "fmt"
  "net"
  "bufio"
  "strings"
)

func main() {
  fmt.Println("Wormhole server starting ...")

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

    // Handle connection
    go func(c net.Conn) {

      line, err := bufio.NewReader(conn).ReadString('\n')
      if err != nil {
        log.Fatal(err)
      }

      fmt.Println(c.RemoteAddr(), ">", line)
      handleLine(c, line)

      c.Close()
    }(conn)
  }
}

func handleLine(c net.Conn, line string) {
  parts := strings.Split(line, " ")

  if len(parts) < 2 {
    return
  }

  switch parts[0] {
    case "EDIT":
      handleCommandEdit(parts[1:])
      return

    case "SHELL":
      handleCommandShell(parts[1:])
      return

    case "EXPLORE":
      handleCommandExplore(parts[1:])
      return

    case "START":
      handleCommandStart(parts[1:])
      return
  }


}

func handleCommandEdit(parts []string) {
  fmt.Println("EDIT", parts)
}

func handleCommandStart(parts []string) {
  fmt.Println("START", parts)

}

func handleCommandExplore(parts []string) {
  fmt.Println("EXPLORE", parts)

}

func handleCommandShell(parts []string) {
  fmt.Println("SHELL", parts)

}