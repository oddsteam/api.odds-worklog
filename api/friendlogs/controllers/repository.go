package controllers

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type Repository struct {
	session *mongo.Session
}

type Event struct {
	Action    string    `bson:"action"`
	Message   string    `bson:"message"`
	CreatedAt time.Time `bson:"createdAt"`
	Version   string    `bson:"version"`
}

func NewRepository(session *mongo.Session) *Repository {
	return &Repository{session}
}

func (r *Repository) Create(action, message string) {
	e := Event{Action: action, Message: message, Version: "0.0.1"}
	e.CreatedAt = time.Now()
	coll := r.session.GetCollection("event")
	err := coll.Insert(e)
	utils.FailOnError(err, "Fail to save event")
}
