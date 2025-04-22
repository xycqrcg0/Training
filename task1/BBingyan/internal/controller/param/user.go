package param

import "time"

type UserRequest struct {
	Code     string `json:"code"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Signature string `json:"signature"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Signature string    `json:"signature"`
	Likes     int       `json:"likes"`
	Follows   int       `json:"follows"`
	CreatedAt time.Time `json:"created-at"`
}
