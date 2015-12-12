
package main

import (
    "fmt"
    "io"
    "path"
    "html/template"
    "log"
    "net/http"
    "code.google.com/p/go.net/websocket"
)

const PORT = "localhost:4000"

func main() {
    http.HandleFunc("/", rootHandler)
    http.Handle("/socket", websocket.Handler(socketHandler))
    err := http.ListenAndServe(PORT, nil)
    if err != nil {
        log.Fatal(err)
    }
}

/*
Request handlers
*/

func rootHandler(w http.ResponseWriter, r *http.Request) {
    fp := path.Join("./src/index.html")
    tmpl, err := template.ParseFiles(fp)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, PORT); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func socketHandler(ws *websocket.Conn) {
    s := socket{ws, make(chan bool)}
    go match(s)
    <-s.done
}

/*
Custom type "socket" and its Close helper function
*/

type socket struct {
    io.ReadWriter
    done chan bool
}

func (s socket) Close() error {
    s.done <- true
    return nil
}

/*
Handling channels and data between connection
*/

var partner = make(chan io.ReadWriteCloser)

func match(c io.ReadWriteCloser) {
    fmt.Fprint(c, "Waiting on someone else to join chat...")
    select {
    case partner <- c:
        // connection gets passed through the partner channel
    case p := <-partner:
        chat(p, c)
    }
}

func chat(a, b io.ReadWriteCloser) {
    fmt.Fprintln(a, "Someone is ready to chat with you!")
    fmt.Fprintln(b, "Someone is ready to chat with you!")
    errc := make(chan error, 1)
    sendData(a, b, errc)
    sendData(b, a, errc)
    if err := <-errc; err != nil {
        log.Println(err)
    }
    a.Close()
    b.Close()
}

func sendData(w io.Writer, r io.Reader, errc chan<- error) {
    _, err := io.Copy(w, r)
    errc <- err
}




