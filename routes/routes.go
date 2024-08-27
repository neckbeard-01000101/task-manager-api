package routes

import (
	"net/http"
	"task-manager/controllers"

	"github.com/gorilla/mux"
)

func ReqisterRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is listening to incoming requests..."))
	})
	r.HandleFunc("/tasks", controllers.AddTask).Methods("POST")
	r.HandleFunc("/tasks", controllers.GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", controllers.GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", controllers.UpdateTaskByID).Methods("PUT")
	r.HandleFunc("/tasks/{id}", controllers.DeleteTaskByID).Methods("DELETE")
}
