package main

import (
	"log"
	"net/http"

	DB "docTrack/config"
	routes "docTrack/routes"
)

func main() {

	dsn := "host=localhost user=postgres password=260322 dbname=DocTrack port=5432 sslmode=disable"
	if err := DB.InitDB(dsn); err != nil {
		log.Fatal("failed to connect to database ", err)
	}

	router := routes.SetupRouter()

	log.Println("Server running on :8080 ")
	log.Fatal(http.ListenAndServe(":8080", router))

}
