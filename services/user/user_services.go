package user

import (
	db "docTrack/config"
	users_model "docTrack/models/users"
	"errors"

	"golang.org/x/crypto/bcrypt"
	// "net/http"
	// "encoding/json"
)

func CreateUser(username, password, role string) error {
	// self explanatory
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}
	// check for existing user
	var existingUser users_model.User
	if err := db.DB.Where("username = ?", username).First(&existingUser).Error; err == nil {
		// the user already exists or same user name
		return errors.New("username already taken ")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create user record
	user := users_model.User{
		Username: username,
		Password: string(hashedPass),
		Role:     role,
	}

	return db.DB.Create(&user).Error
}

func FindUserByUsername(username string) (*users_model.User, error) {
	var user users_model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}
