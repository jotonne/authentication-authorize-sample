package totp

import (
	"crypto/rand"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func CreateQRCode(account string) *otp.Key {
	key, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      "SampleApp",
		AccountName: account,
		Period:      30, // default 30s
		SecretSize:  20,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
		Rand:        rand.Reader,
	})
	return key
}
