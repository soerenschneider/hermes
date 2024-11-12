package domain

import (
	"strings"
	"time"
)

type Notification struct {
	Id                   int64
	Inserted             time.Time
	RetryDate            time.Time
	UnsuccessfulAttempts int

	ServiceId string
	Subject   string
	Message   string
}

func FromNotification(n NotificationRequest) []Notification {
	serviceIds := strings.Split(n.ServiceId, ",")
	var notifications []Notification
	for _, serviceId := range serviceIds {
		not := Notification{
			Inserted:             time.Now(),
			UnsuccessfulAttempts: 0,

			ServiceId: serviceId,
			Subject:   n.Subject,
			Message:   n.Message,
		}

		notifications = append(notifications, not)
	}

	return notifications
}
