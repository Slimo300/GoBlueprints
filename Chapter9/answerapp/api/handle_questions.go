package main

import (
	"net/http"

	"github.com/golang/appengine/datastore"
	"google.golang.org/appengine"
)

func handleQuestions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleQuestionCreate(w, r)
	case "GET":
		params := pathParams(r, "/api/questions/:id")
		questionID, ok := params[":id"]
		if ok { // GET /api/questions/ID
			handleQuestionsGet(w, r, questionID)
			return
		}
		handleTopQuestions(w, r) // GET /api/questions/
	default:
		http.NotFound(w, r)
	}
}

func handleQuestionCreate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var q Question
	err := decode(r, &q)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, q, http.StatusCreated)
}

func handleQuestionsGet(w http.ResponseWriter, r *http.Request, questionID string) {
	ctx := appengine.NewContext(r)
	questionKey, err := datastore.DecodeKey(questionID)
	if err != nil {
		respondErr(ctx, w, r, err, http.StatusBadRequest)
		return
	}
	question, err := GetQuestion(ctx, questionKey)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			respondErr(ctx, w, r, datastore.ErrNoSuchEntity, http.StatusNotFound)
			return
		}
		respondErr(ctx, w, r, err, http.StatusInternalServerError)
		return
	}
	respond(ctx, w, r, question, http.StatusOK)
}
