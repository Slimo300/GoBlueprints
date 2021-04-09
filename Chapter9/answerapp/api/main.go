package api

import (
	"io"
	"log"
	"net/http"
	"os"
)

func handleHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello from App Engine")
}

func main() {
	http.HandleFunc("/", handleHello)
	http.HandleFunc("/api/questions/", handleQuestions)
	http.HandleFunc("/api/answers/", handleAnswers)
	http.HandleFunc("/api/votes/", handleVotes)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
