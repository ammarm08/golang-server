package main

import (
    "fmt"
    "io"
    "log"
    "net"
)

const PORT = "localhost:4000"

func main() {
    server, err := net.Listen("tcp", PORT)
    if err != nil {
      log.Fatal(err)
    }
    for {
      conn, err := server.Accept()
      if err != nil {
        log.Fatal(err)
      }
      go match(conn)
    }
}

var partner = make(chan io.ReadWriteCloser)

func match(conn io.ReadWriteCloser) {
  fmt.Fprint(conn, "Waiting for someone to join chat...")
  select {
  case partner <- conn:
    // data sent from input conn to the "partner" channel
  case p := <- partner:
    // data sent from existing partner to newly initialized p
    // which implicity gets a type io.ReadWriteCloser as well
    // looks like we have a chat on our hands now!
    chat(p, conn)
  }
}

func chat(a, b io.ReadWriteCloser) {
  fmt.Fprintln(a, "Gotchu a friend!")
  fmt.Fprintln(b, "Gotchu a friend!")
  go io.Copy(a, b)
  io.Copy(b, a)
}