package controllers

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"task-manager/config"
	"task-manager/models"
	"time"
)

const (
	DB_NAME         = "task_manager"
	COLLECTION_NAME = "tasks"
)

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Error decoding request body"}`, http.StatusInternalServerError)
		log.Printf("Error decoding request body: %v", err)
		return
	}
	// checks if all the keys are not empty
	if isValid := task.Validate(); !isValid {
		http.Error(w, `{"error": "The received task data is invalid"}`, http.StatusBadRequest)
		log.Println("The received task data is invalid")
		return
	}

	collection := config.Client.Database(DB_NAME).Collection(COLLECTION_NAME)
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
		http.Error(w, `{"error": "Error saving task"}`, http.StatusInternalServerError)
		log.Printf("Error saving task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Task added successfully",
		"id":      result.InsertedID,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Task added successfully with ID: %v", result.InsertedID)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := config.Client.Database(DB_NAME).Collection(COLLECTION_NAME)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, `{"error": "Error fetching tasks"}`, http.StatusInternalServerError)
		log.Printf("Error fetching tasks: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var tasks []bson.M
	if err = cursor.All(ctx, &tasks); err != nil {
		http.Error(w, `{"error": "Error decoding tasks"}`, http.StatusInternalServerError)
		log.Printf("Error decoding tasks: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
	log.Println("Tasks fetched successfully")
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error": "id should be provided as a path param"}`, http.StatusBadRequest)
		log.Println("id should be provided as a path param")
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Println("invalid id")
		return
	}
	collection := config.Client.Database(DB_NAME).Collection(COLLECTION_NAME)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task bson.M
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		http.Error(w, `{"error": "Error fetching task"}`, http.StatusInternalServerError)
		log.Printf("Error fetching task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
	log.Printf("Task fetched successfully with ID: %v", id)
}
func UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error": "id should be provided as a path param"}`, http.StatusBadRequest)
		log.Println("id should be provided as a path param")
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Println("invalid id")
		return
	}

	var taskUpdate bson.M
	err = json.NewDecoder(r.Body).Decode(&taskUpdate)
	if err != nil {
		http.Error(w, `{"error": "Error decoding request body"}`, http.StatusInternalServerError)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	collection := config.Client.Database(DB_NAME).Collection(COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": taskUpdate}
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, `{"error": "Error updating task"}`, http.StatusInternalServerError)
		log.Printf("Error updating task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "Task updated successfully",
		"id":      id,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Task updated successfully with ID: %v", id)
}
func DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, `{"error": "id should be provided as a path param"}`, http.StatusBadRequest)
		log.Println("id should be provided as a path param")
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, `{"error": "invalid id"}`, http.StatusBadRequest)
		log.Println("invalid id")
		return
	}

	collection := config.Client.Database(DB_NAME).Collection(COLLECTION_NAME)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		http.Error(w, `{"error": "Error deleting task"}`, http.StatusInternalServerError)
		log.Printf("Error deleting task: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":      "Task deleted successfully",
		"id":           id,
		"deletedCount": result.DeletedCount,
	}
	json.NewEncoder(w).Encode(response)
	log.Printf("Task deleted successfully with ID: %v", id)
}

