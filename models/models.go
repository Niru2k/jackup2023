package models

//User details
type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

//Token values for each user-id
type Authentication struct {
	Id    uint   `json:"id"`
	Token string `json:"token"`
}

//This slice is used instead of User table
var Database []User

//This slice is used instead of Token table
var Auth []Authentication
