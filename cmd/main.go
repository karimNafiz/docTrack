package main

import (
	DB "docTrack/config"
	routes "docTrack/routes"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {

	dsn := "host=localhost user=postgres password=260322 dbname=DocTrack port=5432 sslmode=disable"
	if err := DB.InitDB(dsn); err != nil {
		log.Fatal("failed to connect to database ", err)
	}

	router := routes.SetupRouter()

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // your UI origin
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Chunk-Index"}),
		handlers.AllowCredentials(),
	)

	log.Println("Server running on :8080 ")
	log.Fatal(http.ListenAndServe(":8080", cors(router)))

}
