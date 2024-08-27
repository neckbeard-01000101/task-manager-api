package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" josn:"description"`
	Category    string             `bson:"category" json:"category"`
	Priority    string             `bson:"priority" json:"priority"`
	Deadline    time.Time          `bson:"deadline" json:"deadline"`
	Status      string             `bson:"status" json:"status"`
}

func (t *Task) Validate() bool {
	return t.Title != "" &&
		t.Description != "" &&
		t.Category != "" &&
		t.Priority != "" &&
		!t.Deadline.IsZero() &&
		t.Status != ""
}
