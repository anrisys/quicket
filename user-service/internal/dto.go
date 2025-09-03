package internal

type UserDTO struct {
	ID       int
	Email    string
	PublicID string
	Role     string
}

type RegisterUserRequest struct {
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8,password"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	Role                 string `json:"role" binding:"omitempty,oneof=user organizer admin"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResponseSuccess struct {
	Code    string `json:"code" example:"SUCCESS"`
	Message string `json:"message" example:"Operation successful"`
}

type LoginUserDTO struct {
	PublicID string `json:"public_id" example:"user_123"`
	Token    string `json:"token" example:"jwt.token.here"`
}

type LoginUserSuccess struct {
	ResponseSuccess `json:",inline"`
	Data            LoginUserDTO `json:"data"`
}

type RegisterUserSuccess struct {
	ResponseSuccess `json:",inline"`
}