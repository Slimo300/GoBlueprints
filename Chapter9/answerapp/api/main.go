package main

import (
	"io"
	"net/http"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from App Engine")
}

func main() {
	http.HandleFunc("/", handleHello)

	http.ListenAndServe(":8080", nil)
}
