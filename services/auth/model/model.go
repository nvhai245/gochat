package model

// User model
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool `json:"isAdmin"`
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
	Username    string `json:"username"`
	Password string `json:"password"`
}

// AuthorizedUser struct
type AuthorizedUser struct {
	IsAdmin  bool   `json:"isAdmin"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
