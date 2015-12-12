
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "code.google.com/p/go.net/websocket"
    "html/template"
)

const PORT = "localhost:4000"

/*
1. rootHandler serves up HTML/JS
2. socketHandler deals with channel logic between concurrent connections
3. server listens on PORT (4000 in this case)
*/
func main() {
    http.HandleFunc("/", rootHandler)
    http.Handle("/socket", websocket.Handler(socketHandler))
    err := http.ListenAndServe(PORT, nil)
    if err != nil {
        log.Fatal(err)
    }
}

/*
We build a socket type that has access to all the functions that a typical
io.ReadWriter has.

It also has a "done" channel that sends either a true or a false. 
This will help our handler know whether a certain connection has been closed or not
*/
type socket struct {
    io.ReadWriter
    done chan bool
}

/*
A helper function that our socket uses to, when the time comes, fire a "TRUE" down
through to its "done" channel, signifying that the connection should be closed 
*/
func (s socket) Close() error {
    s.done <- true
    return nil
}

/*
A websocket connection gets translated into our custom socket type.
It then gets delegated to a goroutine that invokes the function "match"

Note that <-s.done is blocking, so the socketHandler function doesn't
return until an "I'm Done" event happens
*/
func socketHandler(ws *websocket.Conn) {
    s := socket{ws, make(chan bool)}
    go match(s)
    <-s.done
}

// the channel/"piping" through which connections communicate w/ each other
var partner = make(chan io.ReadWriteCloser)

// simultaneous sending/receiving of connections
func match(c io.ReadWriteCloser) {
    fmt.Fprint(c, "Waiting on someone else to join chat...")
    select {
    case partner <- c:
        // connection gets passed through the partner channel
    case p := <-partner:
        // receives connection being passed through partner channel,
        chat(p, c)
    }
}

/*
A sends messages to B
B sends messages to A
Error handling also happens
*/
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

// Templating
func rootHandler(w http.ResponseWriter, r *http.Request) {
    rootTemplate.Execute(w, PORT)
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<script>

var input, output, websocket;

function showMessage(m) {
        var p = document.createElement("p");
        p.innerHTML = m;
        output.appendChild(p);
}

function onMessage(e) {
        showMessage(e.data);
}

function onClose() {
        showMessage("Connection closed.");
}

function sendMessage() {
        var m = input.value;
        input.value = "";
        websocket.send(m);
        showMessage(m);
}

function onKey(e) {
        if (e.keyCode == 13) {
                sendMessage();
        }
}

function init() {
        input = document.getElementById("input");
        input.addEventListener("keyup", onKey, false);

        output = document.getElementById("output");

        websocket = new WebSocket("ws://{{.}}/socket");
        websocket.onmessage = onMessage;
        websocket.onclose = onClose;
}

window.addEventListener("load", init, false);

</script>
</head>
<body>
<input id="input" type="text">
<div id="output"></div>
</body>
</html>
`))


