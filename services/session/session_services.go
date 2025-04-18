package session

import (
	"crypto/rand"
	db "docTrack/config"
	session_model "docTrack/models/sessions"
	user_model "docTrack/models/users"
	"encoding/base64"
	"errors"
	"net/http"
	"time"
)

// need to generate a session ID
// I will generate random bytes in a 32 byte slice

func generateSessionID() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil

}

// creating a session for the user who jus logged in
// we will store it in the sessions table in the dataBase
func CreateSession(user *user_model.User) (*session_model.Session, error) {
	var session session_model.Session
	// we need to generate a token
	token, err := generateSessionID()
	if err != nil {
		// problem in generating a token
		return nil, errors.New("could not generate token ")
	}

	session = session_model.Session{
		ID:        token,
		User_ID:   user.ID,
		Username:  user.Username,
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

func SetSessionCookie(w http.ResponseWriter, session *session_model.Session) error {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.ID,
		Path:     "/",
		Expires:  session.ExpiresAt,
		HttpOnly: true,                    // inaccessible to JS (mitigates XSS)
		Secure:   true,                    // only sent over HTTPS
		SameSite: http.SameSiteStrictMode, // prevents CSRF in most cases

	})

	// todo need to finish
	return nil
}

// func VerifySession (sessionID string) (*session_model.Session, error) {
// 	session , err  := FindSession(sessionID)
// 	if err != nil {
// 		return nil, err
// 	}

// }
