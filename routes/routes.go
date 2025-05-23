package routes

import (
	folder_handler "docTrack/handlers/folder"
	upload_session_handler "docTrack/handlers/upload_session"
	user_handler "docTrack/handlers/users"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// regiser user route
	// no middleware required
	router.HandleFunc("/register", user_handler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", user_handler.LoginHandler).Methods("POST")

	router.HandleFunc("/folder", folder_handler.CreateFolderHandler).Methods("POST")

	router.HandleFunc("/upload", upload_session_handler.InitUploadSession).Methods("POST", "OPTIONS")

	// we need a router for /upload/{uploadID}/chunk?index = smth
	// we need to use regex
	router.HandleFunc("/upload/{uploadID:[0-9a-fA-F\\-]+}/chunk", upload_session_handler.UplodaChunk).Methods("POST", "OPTIONS")
	router.HandleFunc("/upload/{uploadID:[0-9a-fA-F\\-]+}/complete", upload_session_handler.CompleteUploadSession).Methods("POST", "OPTIONS")

	// need a handler for uploading chunks

	router.HandleFunc("/copy", folder_handler.CopyFolderHandler).Methods("POST")
	return router
}
