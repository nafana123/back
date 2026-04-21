package user

type RegistrationRequest struct {
	Login           string `json:"login" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	PasswordConfirm string `json:"password_confirm" binding:"required"`
}

type VerifyRequest struct {
	Code  string `json:"code" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
