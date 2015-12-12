package main

import (
    "fmt"
    "io"
    "log"
    "net"
)

const PORT = "localhost:4000"

func main() {
    // establish connection on given port
    server, err := net.Listen("tcp", PORT)

    // break out of connection if something fishy happens
    if err != nil {
      log.Fatal(err)
    }

    // if we made it this far, that means good stuff happened.
    // we start an infinite loop that accepts incoming connections,
    // handles errors,
    // and passes the incoming connection to our "match" function, then
    // lets that execute in a new goroutine 
    for {
      conn, err := server.Accept()
      if err != nil {
        log.Fatal(err)
      }
      go match(conn)
    }
}

// this a channel for whatever conversation is or is not about to happen.
// "partner" will appear in our "match" function (accessed as a global variable)
var partner = make(chan io.ReadWriteCloser)


func match(conn io.ReadWriteCloser) {
  fmt.Fprint(conn, "Waiting for someone to join chat...")

  // a connection will wait here
  select {
  case partner <- conn:
    // data sent from input conn to the "partner" channel
  case p := <- partner:
    // whatever data was already sent through partner gets 
    // received and initialized in p.
    // we now have access to both the input connection (conn)
    // as well as the previous connection (p)
    chat(p, conn)
  }
}

func chat(a, b io.ReadWriteCloser) {
  fmt.Fprintln(a, "Gotchu a friend!")
  fmt.Fprintln(b, "Gotchu a friend!")

  errconn := make(chan error, 1) // a channel just for errors!
  chatHandle(a, b, errconn)
  chatHandle(b, a, errconn)

  // blocking! only executes this conditional if an error is
  // ever sent and initialized as variable "err"
  if err := <-errconn; err != nil {
    log.Println(err)
  }
  a.Close()
  b.Close()
}

// There's a writer, and there's a reader. And there's an error chan too.
// (the "chan<- error" is Go's way of stating that it's a 
// "write-only" channel of type Error)
func chatHandle(writer io.Writer, reader io.Reader, errconn chan<- error) {
  _, err := io.Copy(writer, reader) // we only care about the err return value, if it exists
  errconn <- err // send err to errconn chan if and when it surfaces
}




