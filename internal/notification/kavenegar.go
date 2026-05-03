package notification

import (
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type kavenegarManager struct {
	api    *kavenegar.Kavenegar
	sender string
}

func NewKavenegarManager(api *kavenegar.Kavenegar, sender string) *kavenegarManager {
	return &kavenegarManager{
		api: api,
		sender: sender,
	}
}

func (k *kavenegarManager) SendSimpleSMS(sms SimpleSMS) error {
	if _, err := k.api.Message.Send(k.sender, sms.Receptor, sms.Massage, nil); err != nil {
		return fmt.Errorf("Notification kavenegar, send simple sms, api send massage: %w", err)
	}
	return nil
}
