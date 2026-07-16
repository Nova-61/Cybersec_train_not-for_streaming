package main

import (
	"errors"
	"fmt"
)

type Notifier interface {
	Send(message string) error
}

type EmailNotifier struct {
	To      string
	From    string
	Subject string
}

func (e EmailNotifier) Send(message string) error {
	if e.To == "" {
		return errors.New("email: recipient is empty")
	}
	fmt.Printf("Отправлено email на [%s] от [%s] с темой '%s': %s\n", e.To, e.From, e.Subject, message)
	return nil
}

type SMSNotifier struct {
	PhoneNumber string
}

func (s SMSNotifier) Send(message string) error {
	if s.PhoneNumber == "" {
		return errors.New("sms: phone number is empty")
	}
	fmt.Printf("Отправлено SMS на [%s]: %s\n", s.PhoneNumber, message)
	return nil
}

type PushNotifier struct {
	DeviceID string
}

func (p PushNotifier) Send(message string) error {
	if p.DeviceID == "" {
		return errors.New("push: device ID is empty")
	}
	fmt.Printf("Отправлен Push на [%s]: %s\n", p.DeviceID, message)
	return nil
}

func SendAll(notifiers []Notifier, message string) {
	for i, notifier := range notifiers {
		err := notifier.Send(message)
		if err != nil {
			fmt.Printf("Ошибка в уведомлении #%d: %v\n", i+1, err)
		}
	}
}

func main() {
	email := EmailNotifier{
		To:      "user@example.com",
		From:    "noreply@example.com",
		Subject: "Hello",
	}

	sms := SMSNotifier{
		PhoneNumber: "+79001234567",
	}

	push := PushNotifier{
		DeviceID: "device-123",
	}

	notifiers := []Notifier{email, sms, push}

	SendAll(notifiers, "Привет, мир!")
}
