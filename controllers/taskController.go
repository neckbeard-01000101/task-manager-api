package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if isValid := task.Validate(); !isValid {
		http.Error(w, "The received task data is invalid", http.StatusBadRequest)
		return
	}
	collection := config.Client.Database("task_manager").Collection("tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, bson.M{
		"title":       task.Title,
		"description": task.Description,
		"category":    task.Category,
		"priority":    task.Priority,
		"deadline":    task.Deadline,
		"status":      task.Status,
	})
	if err != nil {
		http.Error(w, "Error saving task: "+err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task added ID: %s", result.InsertedID)
}
