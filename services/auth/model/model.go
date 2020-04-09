package model

import (
	"time"
)

// User model
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

// Users request for all users information
type Users struct {
	Users []User `json:"users"`
}

// Message model
type Message struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Messages array of an user
type Messages struct {
	Messages []Message `json:"messages"`
}

// LoginData struct
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthorizedUser struct
type AuthorizedUser struct {
	IsAdmin  bool      `db:"isadmin"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	Avatar   string    `db:"avatar"`
	Phone    string    `db:"phone"`
	Birthday time.Time `db:"birthday"`
	Fb       string    `db:"fb"`
	Insta    string    `db:"insta"`
}
