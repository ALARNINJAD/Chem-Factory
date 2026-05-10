package bale

import (
	"bytes"
	baleDTO "chem-factory/internal/dto/bale"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func (bale *manager) SendOTP(phone, otp string) error {
    
	// 989123456789
	if len(phone) != 12 && len(phone) != 13 {
		return errors.New("Bale, send otp: phone number is incorrect.")
	}
	if (phone[0] != '9' || phone[1] != '8') && (phone[0] != '+' || phone[1] != '9' || phone[2] != '8') {
		return errors.New("Bale, send otp: phone number is incorrect. it must be 98********* or +98*********.")
	}

	// 4444 / 55555 / 666666 / 7777777 / 88888888
	if len(otp) > 8 || len(otp) < 4 {
		return errors.New("Bale, send otp: unacceptable otp length.")
	}
	for _, char := range otp {
		if !unicode.IsDigit(char) {
			return errors.New("Bale, send otp: otp is not numeric.")
		}
	}

    requestID := uuid.New().String()

    var lastErr error
    for attempt := range 3 {

		var statCode int
        statCode, lastErr = bale.doSendOTP(requestID, phone, otp)
        if lastErr == nil {
            return nil
        }
		if statCode == http.StatusInternalServerError || statCode == 1 {
        	time.Sleep(time.Duration(attempt+1) * time.Second)
			continue
		}
		return lastErr
    }
    return fmt.Errorf("Bale, send otp: all attempts failed: %w", lastErr)
}

func (bale *manager) doSendOTP(requestID, phone, otp string) (int, error) {

	reqBody := baleDTO.OTPRequest{
		RequestID:   requestID,
		BotID:       bale.botID,
		PhoneNumber: phone,
		MessageData: baleDTO.MessageData{OTPMessage: baleDTO.OTPMessage{OTP: otp}},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return -1, fmt.Errorf("Bale, send otp, json marshal: %w", err)
	}

	req, err := http.NewRequest("POST", bale.url, bytes.NewBuffer(data))
	if err != nil {
		return -1, fmt.Errorf("Bale, send otp, new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-access-key", bale.apiKey)

	resp, err := bale.client.Do(req)
	if err != nil {
		return 1, fmt.Errorf("Bale, send otp, do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, fmt.Errorf("Bale, send otp, read response body: %w", err)
	}

	switch resp.StatusCode {

	case http.StatusOK:
		return http.StatusOK, nil

	case http.StatusBadRequest:
		return http.StatusBadRequest, fmt.Errorf("Bale, send otp, bad Request (400): %s", body)
	
	case http.StatusUnauthorized:
		return http.StatusUnauthorized, errors.New("Bale, send otp, unauthorized (401): invalid api key.")
	
	case http.StatusInternalServerError:
		return http.StatusInternalServerError, errors.New("Bale, send otp, internal Server Error (500)")
	
	default:
		return resp.StatusCode, fmt.Errorf("Bale, send otp, error: %d, %s", resp.StatusCode, body)
	}
}