package dto

type RegisterUserRequest struct {
	Email                string `json:"user" binding:"required,email"`
	Password             string `json:"password" binding:"required,alphanumunicode,min=6"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	Role                 string `json:"role" binding:"omitempty,oneof=user organizer admin"`
}

type LoginUserRequest struct {
	Email    string `json:"user" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}