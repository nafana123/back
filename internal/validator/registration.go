package validator

import (
	"back/internal/dto/user"
	"strings"
)

func ValidateRegistration(req *user.RegistrationRequest) (bool, string, string) {
	if strings.TrimSpace(req.Login) == "" {
		return false, "Логин не может быть пустым", "login"
	}
	if len(req.Login) < 3 || len(req.Login) > 20 {
		return false, "Логин должен быть от 3 до 20 символов", "login"
	}
	
	if req.Email == "" {
		return false, "Email не может быть пустым", "email"
	}
	
	if req.Password == "" {
		return false, "Пароль не может быть пустым", "password"
	}
	if len(req.Password) < 6 {
		return false, "Пароль должен быть не менее 6 символов", "password"
	}
	
	if req.Password != req.PasswordConfirm {
		return false, "Пароль и подтверждение пароля не совпадают", "password_confirm"
	}
	
	return true, "", ""
}

func ValidateVerify(req *user.VerifyRequest) (bool, string, string) {
	if strings.TrimSpace(req.Login) == "" {
		return false, "Логин не может быть пустым", "login"
	}
	if len(req.Login) < 3 || len(req.Login) > 20 {
		return false, "Логин должен быть от 3 до 20 символов", "login"
	}

	if req.Email == "" {
		return false, "Email не может быть пустым", "email"
	}

	if req.Password == "" {
		return false, "Пароль не может быть пустым", "password"
	}
	if len(req.Password) < 6 {
		return false, "Пароль должен быть не менее 6 символов", "password"
	}

	code := strings.TrimSpace(req.Code)
	if code == "" {
		return false, "Введите код подтверждения", "code"
	}
	if len(code) != 6 {
		return false, "Код должен состоять из 6 цифр", "code"
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			return false, "Код должен содержать только цифры", "code"
		}
	}

	return true, "", ""
}

func ValidateLogin(req *user.LoginRequest) (bool, string, string) {
	if strings.TrimSpace(req.Email) == "" {
		return false, "Почта не может быть пустой", "email"
	}
	
	if req.Password == "" {
		return false, "Пароль не может быть пустым", "password"
	}
	
	return true, "", ""
}