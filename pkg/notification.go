package pkg

import "time"

type Notification struct {
	Inserted             time.Time
	UnsuccessfulAttempts int

	ServiceId string
	Subject   string
	Message   string
}

func FromNotification(n NotificationRequest) Notification {
	return Notification{
		Inserted:             time.Now(),
		UnsuccessfulAttempts: 0,

		ServiceId: n.ServiceId,
		Subject:   n.Subject,
		Message:   n.Message,
	}
}
