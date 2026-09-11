package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"siem-agent/internal/models"
	"siem-agent/internal/retry"
)

type Sender struct {
	serverURL string
	apiKey    string
	client    *http.Client
	retrier   *retry.Retrier
}

func NewSender(serverURL, apiKey string, retryCount, retryDelaySeconds int) *Sender {
	return &Sender{
		serverURL: serverURL,
		apiKey:    apiKey,
		client:    &http.Client{Timeout: 10 * time.Second},
		retrier:   retry.NewRetrier(retryCount, time.Duration(retryDelaySeconds)*time.Second),
	}
}

// SendBatch отправляет пачку логов одним HTTP-запросом, с повторами при ошибке
func (s *Sender) SendBatch(batch []models.LogEntry) error {
	if len(batch) == 0 {
		return nil
	}

	data, err := json.Marshal(models.LogBatch{
		Entries: batch,
		Count:   len(batch),
		Source:  "siem-agent",
	})
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}

	return s.retrier.Do(func() error {
		req, err := http.NewRequest("POST", s.serverURL, bytes.NewReader(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("server returned %d", resp.StatusCode)
		}
		return nil
	})
}
