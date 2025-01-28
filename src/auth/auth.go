package auth

import (
	"sazardev.clean-menu-go/src/models"
)

var CurrentUser models.User

func SetCurrentUser(user models.User) {
	CurrentUser = user
}

func GetCurrentUser() models.User {
	return CurrentUser
}
