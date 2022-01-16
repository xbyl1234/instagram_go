package verification

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/verification/emali"
	"makemoney/common/verification/phone"
	"strings"
)

type VerificationType int

var TypeEmail VerificationType = 0
var TypePhone VerificationType = 1

type VerificationCodeProvider interface {
	RequireAccount() (string, error)
	RequireCode(number string) (string, error)
	ReleaseAccount(number string) error
	BlackAccount(number string) error
	GetBalance() (string, error)
	GetProvider() string
	GetArea() string
	Login() error
	GetType() VerificationType
}

func GetCode(msg string) string {
	var index = 0
	find := false
	for index = range msg {
		if msg[index] >= '0' && msg[index] <= '9' {
			find = true
			break
		}
	}
	if find {
		code := strings.ReplaceAll(msg[index:index+7], " ", "")
		if len(code) != 6 {
			return ""
		}
		return code
	} else {
		return ""
	}
}

type Provider struct {
	ProviderType string          `json:"provider_type"`
	ProviderName string          `json:"provider_name"`
	Phone        phone.PhoneInfo `json:"phone"`
	Email        emali.EmailInfo `json:"email"`
}

var VerificationProvider map[string]VerificationCodeProvider

func InitVerificationProviderByJson(config []*Provider) error {
	var err error
	for _, item := range config {
		var Provider VerificationCodeProvider
		if item.ProviderType == "email" {
			Provider, err = emali.InitEmailVerification(&item.Email)
		} else {
			Provider, err = phone.InitPhoneVerification(&item.Phone)
		}
		if err != nil {
			log.Error("init provider %s error: %v", item.ProviderName, err)
		} else {
			VerificationProvider[item.ProviderName] = Provider
		}
	}
	if len(VerificationProvider) == 0 {
		return &common.MakeMoneyError{
			ErrStr: "no verification code provider",
		}
	}
	return nil
}
