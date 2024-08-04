package pkg

type NotificationRequest struct {
	ServiceId  string `json:"service_id" validate:"required"`
	RoutingKey string `json:"routing_key"`
	Subject    string `json:"subject"`
	Message    string `json:"message" validate:"required"`
}
