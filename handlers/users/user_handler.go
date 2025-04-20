package users

import (
	session_service "docTrack/services/session"
	user_service "docTrack/services/user"
	utils "docTrack/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// the requestBody is the expected JSON structure for registration
type requestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// registering users through a post method

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// close the body when we are done
	// so the tcp connection can be re-used
	defer r.Body.Close()

	var body requestBody

	// Decode and validate incoming json
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {

		http.Error(w, "Invalid JSON body ", http.StatusBadRequest)
		return
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// important
	defer r.Body.Close()

	var body requestBody

	// decode the incoming JSON
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {

		http.Error(w, "Invalid JSON body ", http.StatusBadRequest)
		return
	}

	// search if the user exists or not
	user, err := user_service.FindUserByUsername(body.Username)
	if err != nil {
		fmt.Println("coudl not find user by user name ")
		http.Error(w, "Invalid Credentials ", http.StatusBadRequest)
		return
	}

	// if user exists we check the password
	flag, _ := utils.VerifyPassword(body.Password, user.Password)

	if !flag {
		fmt.Println("password doesnt match ")
		http.Error(w, "Invalid Credentials ", http.StatusBadRequest)
		return

	}

	//if all the credentials have been validated
	// we need to create a session
	session_ptr, err := session_service.CreateSession(user)

	if err != nil {
		fmt.Println("could not create session ")
		http.Error(w, "Server Error could not create session ", http.StatusInternalServerError)
		return
	}

	// after the session has been created we have a pointer to that session

	err = session_service.SetSessionCookie(w, session_ptr)

	if err != nil {
		fmt.Println("could not set session cookie ")
		http.Error(w, "Server Error could not create session ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user successfully logged in ",
	})

}
