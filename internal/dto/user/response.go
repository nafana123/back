package user

type ErrorResponse struct {
	Error string `json:"error"`
	Field string `json:"field,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}