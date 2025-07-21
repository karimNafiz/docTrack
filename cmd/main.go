package main

import (
	"context"
	DB "docTrack/config"
	"docTrack/file_upload_service"
	"docTrack/global_configs"
	logger "docTrack/logger"
	routes "docTrack/routes"

	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {

	logger.Init()
	dsn := "host=localhost user=postgres password=260322 dbname=DocTrack port=5432 sslmode=disable"
	if err := DB.InitDB(dsn); err != nil {
		log.Fatal("failed to connect to database ", err)
	}
	fUploadData, err := connectToFileUploadService(context.Background())
	if err != nil {
		log.Fatal("failed to connect to file upload ", err)
		return
	}
	router := routes.SetupRouter(fUploadData)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // your UI origin
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Chunk-Index"}),
		handlers.AllowCredentials(),
	)

	log.Println("Server running on :8080 ")
	log.Fatal(http.ListenAndServe(":8080", cors(router)))

}

func connectToFileUploadService(pContext context.Context) (*file_upload_service.FileUploadServiceInfo, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}
	payload := struct {
		Host                    string `json:"host"`
		Scheme                  string `json:"scheme"`
		Port                    string `json:"port"`
		UploadStatusCallBackUrl string `json:"uploadStatusCallBackUrl"`
	}{
		Host:                    global_configs.MAINSERVICEDOMAIN,
		Scheme:                  global_configs.MAINSERVICESCHEME,
		Port:                    global_configs.MAINSERVICEPORT,
		UploadStatusCallBackUrl: global_configs.FILEUPLOADSERVICECALLBACKURL,
	}
	fUploadSerData, err := file_upload_service.RegisterToFileUploadService(pContext, global_configs.FILEUPLOADSERVICESCHEME, global_configs.FILEUPLOADSERVICEDOMAIN, global_configs.FILEUPLOADSERVICEPORT, global_configs.FILEUPLOADSERVICEREGISTERENDPOINT, header, payload)

	return fUploadSerData, err
}

// need a function to subscribe file upload service
// we will get a service id back
