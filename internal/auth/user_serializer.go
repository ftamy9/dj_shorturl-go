package auth

type SignupRequest struct {
	ID           string `json:"id"`
	PasswordHash string `json:"password_hash"`
}

type SignupResponse struct {
	ID         string `json:"id"`
	PasswordHash string `json:"password_hash"`
}

type LoginResponse struct {
	Token string `json:"toke"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
