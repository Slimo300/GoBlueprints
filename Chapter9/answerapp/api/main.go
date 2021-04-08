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
	http.HandleFunc("/api/questions/", handleQuestions)
	http.HandleFunc("/api/answers/", handleAnswers)
	http.HandleFunc("/api/votes/", handleVotes)

	http.ListenAndServe(":8080", nil)
}
