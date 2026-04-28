package notification

import (
	"fmt"

	"github.com/kavenegar/kavenegar-go"
)

type kavenegarManager struct {
	API    *kavenegar.Kavenegar
	Sender string
}

func (k *kavenegarManager) SendSimpleSMS(sms SimpleSMS) error {

	if _, err := k.API.Message.Send(k.Sender, sms.Receptor, sms.Massage, nil); err != nil {
		return fmt.Errorf("Notification kavenegar, send simple sms, api send massage: %w ", err)
	}

	return nil
}
