I wanted to learn Go, so here we are. This repo is a sandbox of learning and experimenting.

Currently, what I have is a simple chat app using websockets.

The code has been heavily borrowed from Andrew Gerrand's incredible talk (https://vimeo.com/53221560),
but it will morph as I continue to learn more. 

Also, the websocket library I've used is deprecated, so I will be refactoring the code to account for that.

I also have more or less the same thing implemented in NodeJS. Wanted to compare the two implementations.

###Installation###

If you're new to Go, I highly recommend going here first: https://golang.org/doc/code.html. There's some stuff about directory structuring and environment variables that if you skip over will make life really frustrating.

Oh, and check out https://gobyexample.com/ to get up to speed on the language itself!

####Go####

(If this doesn't work exactly right, let me know)

```
// assuming you're on project root ./golang-server
export GOPATH=$HOME/....../<your_project_root>
go get
go build src/server.go
go run src/server.go // server now listening on port 4000
```

####NodeJS####

```
// assuming you're on project root ./golang-server
cd src/node-server
npm install
node server.js // server now listening on port 3000
```

###Directory Structure###

```
golang-server
  /pkg
  /src 
    server.go
    /code.google.com
    /node-server
      server.js
      index.html
      package.json
```


