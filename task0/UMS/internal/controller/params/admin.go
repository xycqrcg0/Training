package params

type AdminRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type HandleUserRequest struct {
	Email string `json:"email" param:"email"`
}
