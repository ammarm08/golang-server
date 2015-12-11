package main

import (
    "fmt"
    "log"
    "net/http"
)

const PORT = "4000"

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe("localhost:"+PORT, nil)
    if err != nil {
        log.Fatal(err)
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Alohamora")
}