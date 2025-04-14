package params

import (
	"time"
)

type UserLoginRequest struct {
	Email    string
	Password string
}

type UserRegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type UserResponse struct {
	Name      string
	Email     string
	CreatedAt time.Time
}
