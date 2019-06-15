package models

type User struct {
	User_ID  string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Email    string `json:"email"`
}