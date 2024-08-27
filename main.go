package main

import (
	"fmt"
	"github.com/gofor-little/env"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"task-manager/config"
	"task-manager/routes"
)

func main() {
	err := env.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
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
