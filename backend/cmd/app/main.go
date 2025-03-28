package main

import (
	"fmt"
	"log"
	"net/http"

	"social-network/internal/api"
	"social-network/internal/repository"
	"social-network/pkg/middlewares"
)

func main() {
	db, err := repository.OpenDb()
	if err != nil {
		return
	}
	
	// Set flags to include file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if err := repository.ApplyMigrations(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: middlewares.CORS(middlewares.AuthMiddleware(api.Routes(db), db)),
	}

	fmt.Println("http://localhost:8080/")
	err = server.ListenAndServe()
	if err != nil {
		log.Println("Error in starting of server:", err)
		return
	}
}
