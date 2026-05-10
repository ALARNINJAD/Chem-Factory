package otp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type manager struct {}

func New() *manager { return &manager{} }

func (otp manager) Generate(numberOfDigits int8) (string, error) {
	if numberOfDigits < 1 || numberOfDigits > 10 {
		return "", errors.New("OTP, generate: number of digits is unsupported.")
	}
	min := int32(1)
	for i := int8(1); i < numberOfDigits; i++ {
		min *= 10
	}
	max := min*10 - 1
	rangeSize := big.NewInt(int64(max - min + 1))
	randomBig, err := rand.Int(rand.Reader, rangeSize)
	if err != nil {
		return "", fmt.Errorf("OTP, failed to generate random otp: %w", err)
	}
	otpCode := min + int32(randomBig.Int64())
	return string(otpCode), nil
}
