package auth_handler_v1

type SigninPayload struct {
	Login          string `json:"login" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6,max=16"`
	RememberMe     bool   `json:"remember_me"`
	TurnstileToken string `json:"turnstile_token" validate:"required"`
}

type PasswordResetPayload struct {
	Password       string `json:"password" validate:"required,min=6,max=16"`
	TurnstileToken string `json:"turnstile_token" validate:"required"`
}
