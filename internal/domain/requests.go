package domain

type NotificationRequest struct {
	ServiceId string `json:"service_id" validate:"required"`
	Subject   string `json:"subject"`
	Message   string `json:"message" validate:"required"`
}
