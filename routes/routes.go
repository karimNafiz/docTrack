package routes

import (
	user_handler "docTrack/handlers/users"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// regiser user route
	// no middleware required
	router.HandleFunc("/register", user_handler.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", user_handler.LoginHandler).Methods("POST")

	return router
}
