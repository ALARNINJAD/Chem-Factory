package bale

type OTPRequest struct {
	BotID       int         `json:"bot_id"`
	PhoneNumber string      `json:"phone_number"`
	RequestID   string      `json:"request_id,omitempty"`
	MessageData MessageData `json:"message_data,omitempty"`
}

type MessageData struct {
	OTPMessage OTPMessage `json:"otp_message,omitempty"`
}

type OTPMessage struct {
	OTP string `json:"otp"`
}
