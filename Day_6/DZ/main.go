// я хз как это делать вообще если честно потому что
package main

import (
	"errors"
	"fmt"
	"os"
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

type Logger struct{}

func (l Logger) Log(message string) {
	fmt.Println("[LOG]", message)
}

func (l Logger) Send(message string) error {
	l.Log("Отправка уведомления: " + message)
	fmt.Printf("Лог отправлен: %s\n", message)
	return nil
}

type Storage interface {
	Save(data string) error
	Load() (string, error)
}

type FileStorage struct {
	Filename string
}

func (f FileStorage) Save(data string) error {
	return os.WriteFile(f.Filename, []byte(data), 0644)
}

func (f FileStorage) Load() (string, error) {
	data, err := os.ReadFile(f.Filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type MemoryStorage struct {
	data string
}

func (m *MemoryStorage) Save(data string) error {
	m.data = data
	return nil
}

func (m *MemoryStorage) Load() (string, error) {
	if m.data == "" {
		return "", errors.New("no data in memory")
	}
	return m.data, nil
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

	logger := Logger{}

	notifiers := []Notifier{email, sms, push, logger}

	SendAll(notifiers, "Привет, мир!")

	fmt.Println("\n--- Storage ---")

	fileStorage := FileStorage{Filename: "data.txt"}
	fileStorage.Save("Hello from file!")
	data, _ := fileStorage.Load()
	fmt.Println("FileStorage:", data)

	memoryStorage := &MemoryStorage{}
	memoryStorage.Save("Hello from memory!")
	data, _ = memoryStorage.Load()
	fmt.Println("MemoryStorage:", data)
}
