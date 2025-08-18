package dto

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