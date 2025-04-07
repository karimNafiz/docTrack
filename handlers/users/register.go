package users

import (
	user_service "docTrack/services/user"
	"encoding/json"
	"net/http"
)

// the requestBody is the expected JSON structure for registration
type requestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:role`
}

// registering users through a post method

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var body requestBody

	// Decode and validate incoming json
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON body ", http.StatusBadRequest)
	}

	if body.Role == "" {
		body.Role = "user"
	}

	err := user_service.CreateUser(body.Username, body.Password, body.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})

}
