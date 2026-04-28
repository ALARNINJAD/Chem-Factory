package notification

import (
	"log"
	"os"

	"github.com/kavenegar/kavenegar-go"
)

type SimpleSenderSMS interface {
	SendSimpleSMS(sms SimpleSMS) error
}

type NotificationManager struct {
	Kavenegar *kavenegarManager
	Mediana *medianaManager
}

type SimpleSMS struct {
	Receptor []string `json:"receptor"`
	Massage  string   `json:"massage"`
}

func (n *NotificationManager) SendSMSWithProvider(sender SimpleSenderSMS, sms SimpleSMS) error {
	return sender.SendSimpleSMS(sms)
}

func Init() *NotificationManager {
	log.Println(os.Getenv("KAVENEGAR_API"))
	return &NotificationManager{
		Kavenegar: &kavenegarManager{
			API: kavenegar.New(os.Getenv("KAVENEGAR_API_KEY")),
			Sender: os.Getenv("KAVENEGAR_SENDER"),
		}, Mediana: &medianaManager{
			APIKey: os.Getenv("MEDIANA_API_KEY"),
			URL: os.Getenv("MEDIANA_URL"),
			Sender: os.Getenv("MEDIANA_SENDER"),
		},
	}
}
