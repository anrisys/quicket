package dto

type ResponseSuccess struct {
	Code    string `json:"code" example:"SUCCESS"`
	Message string `json:"message" example:"Operation successful"`
}

type CreateEventSuccessResponse struct {
	ResponseSuccess `json:",inline"`
	Event           SimpleEventDTO `json:"event"`
}