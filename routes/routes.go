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
	r.HandleFunc("/add-task", controllers.AddTask).Methods("POST")
}
