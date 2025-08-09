package dto

type LoginUserResponse struct {
	PublicID string `json:"public_id"`
	Token    string `json:"token"`
}