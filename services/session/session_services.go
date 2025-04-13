package session

import (
	db "docTrack/config"
	session_model "docTrack/models/sessions"
	user_model "docTrack/models/users"
	util "docTrack/utils"
	"errors"
	"time"
)

// creating a session for the user who jus logged in
// we will store it in the sessions table in the dataBase
func CreateSession(user *user_model.User) (*session_model.Session, error) {
	var session session_model.Session
	// we need to generate a token
	token, err := util.GenerateTokens(user.Username, user.Role)
	if err != nil {
		// problem in generating a token
		return nil, errors.New("could not generate token ")
	}

	session = session_model.Session{
		ID:        token,
		UserID:    user.ID,
		Username:  session.Username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(2 * time.Hour), // ill change this later no hard coding
	}

	if err := db.DB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil

}

func FindSession(sessionID string) (*session_model.Session, error) {
	var session session_model.Session
	err := db.DB.Where("session_id = ?", sessionID).First(&session).Error
	return &session, err
}

// func VerifySession (sessionID string) (*session_model.Session, error) {
// 	session , err  := FindSession(sessionID)
// 	if err != nil {
// 		return nil, err
// 	}

// }
