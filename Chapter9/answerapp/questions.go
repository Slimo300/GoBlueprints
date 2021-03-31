package answerapp

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type Question struct {
	Key      *datastore.Key `json:"id" datastore:"-"`
	CTime    time.Time      `json:"created"`
	Question string         `json:"question"`
	User     UserCard       `json:"answers_count"`
}
