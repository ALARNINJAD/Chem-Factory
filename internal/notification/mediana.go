package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type medianaManager struct {
	url    string
	apiKey string
	sender string
}

func NewMedianaManager(url, apiKey, sender string) *medianaManager {
	return &medianaManager{
		url: url,
		apiKey: apiKey,
		sender: sender,
	}
}

type medianaSimpleRequest struct {
	SendingNumber string   `json:"sendingNumber"`
	Recipients    []string `json:"recipients"`
	MessageText   string   `json:"messageText"`
}

func (m *medianaManager) SendSimpleSMS(sms SimpleSMS) error {

	reqBody := medianaSimpleRequest{
		SendingNumber: m.sender,
		Recipients:    sms.Receptor,
		MessageText:   sms.Massage,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("Notification mediana, send simple sms, json marshal: %w ", err)
	}

	req, err := http.NewRequest("POST", m.url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Notification mediana, send simple sms, new request: %w ", err)
	}

	req.Header.Set("Authorization", "ApiKey "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Notification mediana, send simple sms, do request: %w ", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("Notification mediana, send simple sms, bad Request (400): %s", body)
	case http.StatusUnauthorized:
		return errors.New("Notification mediana, send simple sms, unauthorized (401): invalid api key.")
	case http.StatusInternalServerError:
		return errors.New("Notification mediana, send simple sms, internal Server Error (500)")
	default:
		return fmt.Errorf("Notification mediana, send simple sms, error: %d - %s", resp.StatusCode, body)
	}
}
