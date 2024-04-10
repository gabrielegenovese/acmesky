package flightMatcher

import (
	"acmesky/entities"
	"acmesky/services/notification/prontogram"
	"fmt"
)

type NotificationResult struct {
	Id        string
	ServiceId string
}
type Notification struct {
	subject string
	content string
}

func NotifyCustomer(pref entities.CustomerFlightSubscriptionRequest, notification Notification) ([]NotificationResult, []error) {
	var notifications []NotificationResult
	var errors []error
	if pref.ProntogramID != "" {
		res, err := strategyNotify_Prontogram(pref, notification)
		if err != nil {
			errors = append(errors, err)
		} else {
			notifications = append(notifications, NotificationResult{
				Id:        fmt.Sprintf("%v", res.MessageId),
				ServiceId: "PRONTOGRAM",
			})
		}
	}

	return notifications, errors
}

func strategyNotify_Prontogram(pref entities.CustomerFlightSubscriptionRequest, notification Notification) (prontogram.SendMessageResponse, error) {

	messageContent := notification.subject + "\n" + notification.content + "\n"
	fmt.Printf("Notify Prontogram customer (%s) with message:\n%s", pref.ProntogramID, messageContent)

	return prontogram.SendMessage(messageContent, pref.ProntogramID)
}
