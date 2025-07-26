package dto

type RegisterRequest struct {
	email                 string `json:"user" binding:"required,email"`
	password              string `json:"password" binding:"required"`
	password_confirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
}