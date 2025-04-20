package routes

import (
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

	router.HandleFunc("/initUpload", upload_session_handler.InitUploadSession).Methods("POST")

	return router
}
