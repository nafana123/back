package auth

type UserData struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	PhotoURL     string `json:"photo_url"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Login string `json:"login"`
	Role  string `json:"role"`
}

type GoogleUserResponse struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

type PendingRegistration struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Code      string `json:"code"`
	ExpiresAt int64  `json:"expires_at"`
}
