package user

import "fmt"

func GetErrInvalidUserID(userID uint) error {
	return fmt.Errorf("the user id %d is invalid", userID)
}
