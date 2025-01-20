package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Image    string `json:"image"`
}

const (
	ADMINISTRATOR = "administrator"
	WAITER        = "waiter"
	CHEF          = "chef"
)

var Users []User
