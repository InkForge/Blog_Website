package domain

type INotificationService interface {
	SendEmail(to string, subject string, body string) error
}
