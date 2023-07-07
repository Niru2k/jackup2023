package models

import (
	"time"
)

// Login credentials
type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup credentials
type SignupReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Post Request
type PostReq struct {
	PostTitle   string `json:"post_title"`
	PostContent string `json:"post_content"`
	Catagory    string `json:"catagory"`
}

// User details
type User struct {
	// gorm.Model
	UserId   uint   `json:"-" gorm:"primarykey"`
	Username string `json:"username" gorm:"column:username;type:varchar(100)"`
	Email    string `json:"email" gorm:"column:email;type:varchar(100) unique"`
	Password string `json:"password" gorm:"column:password;type:varchar(100)"`
	Role     string `json:"role" gorm:"-:all"`
	RoleId   uint   `json:"-" gorm:"column:role_id;type:bigint references Roles(role_id)"`
}

// Roles table
type Roles struct {
	RoleId uint   `gorm:"column:role_id;type:bigint primary key"`
	Role   string `gorm:"column:role;type:varchar(50)"`
}

// Catagory Table
type Catagory struct {
	CatagoryId uint   `gorm:"column:catagory_id;type:bigint primary key"`
	Catagory   string `gorm:"column:catagory;type:varchar(50)"`
}

// Token values for each user-id
type Authentication struct {
	UserId uint   `json:"user_id" gorm:"column:user_id;type:bigint primary Key"`
	Token  string `json:"token" gorm:"column:token;type:varchar(200)"`
}

// Post details
type Post struct {
	PostId      uint      `json:"-" gorm:"primarykey"`
	PostTitle   string    `json:"post_title,omitempty" gorm:"column:post_title;type:varchar(100)"`
	PostContent string    `json:"post_content,omitempty" gorm:"column:post_content;type:varchar(500)"`
	Catagory    string    `json:"catagory,omitempty" gorm:"-"`
	CatagoryId  uint      `json:"-" gorm:"column:catagory_id;type:bigint references Catagories(catagory_id)"`
	UserId      uint      `json:"-" gorm:"column:user_id;type:bigint references Users(user_id)"`
	Date        time.Time `json:"-" gorm:"autoCreateTime"`
}

// Comments
type Comments struct {
	CommentId uint   `json:"comment_id"`
	UserId    uint   `json:"user_id"`
	Email     string `json:"email"`
	PostId    uint   `json:"post_id"`
	PostTitle string `json:"post_title"`
	Comment   string `json:"comment"`
}

// Comments request
type CommentReq struct {
	PostTitle string `json:"post_title"`
	Email     string `json:"email"`
	Comment   string `json:"comment"`
}

// This slice is used instead of Post table
var CommentTable []Comments
