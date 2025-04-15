package param

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
