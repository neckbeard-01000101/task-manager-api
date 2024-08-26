package main

import (
	"fmt"
	"log"
	"net/http"
	"task-manager/config"
	"task-manager/routes"

	"github.com/gofor-little/env"
	"github.com/gorilla/mux"
)

func main() {
	err := env.Load(".env")
	r := mux.NewRouter()
	routes.ReqisterRoutes(r)
	server := &http.Server{
		Handler: r,
		Addr:    "localhost:8080",
	}
	config.InitDB()
	fmt.Printf("Server is running on %v", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
