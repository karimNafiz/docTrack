package routes

import (
	users "docTrack/handlers/users"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// regiser user route
	// no middleware required
	router.HandleFunc("/register", users.RegisterHandler).Methods("POST")

	return router
}
