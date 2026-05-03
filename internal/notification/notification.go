package notification

import (
	"os"

	"github.com/kavenegar/kavenegar-go"
)

type SimpleSMSSender interface {
	SendSimpleSMS(sms SimpleSMS) error
}

type Manager struct {
	Kavenegar *kavenegarManager
	Mediana   *medianaManager
}

type SimpleSMS struct {
	Receptor []string `json:"receptor"`
	Massage  string   `json:"massage"`
}

func (m *Manager) SendSMSWithProvider(sender SimpleSMSSender, sms SimpleSMS) error {
	return sender.SendSimpleSMS(sms)
}

func New() *Manager {
	return &Manager{
		Kavenegar: NewKavenegarManager(
			kavenegar.New(os.Getenv("KAVENEGAR_API_KEY")),
			os.Getenv("KAVENEGAR_SENDER"),
		),
		Mediana: NewMedianaManager(
			os.Getenv("MEDIANA_URL"),
			os.Getenv("MEDIANA_API_KEY"),
			os.Getenv("MEDIANA_SENDER"),
		),
	}
}
