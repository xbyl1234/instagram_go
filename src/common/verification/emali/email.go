package emali

import (
	"makemoney/common"
	"makemoney/common/verification"
	"net/http"
	"sync"
	"time"
)

type EmailInfo struct {
	Domain           string `json:"domain"`
	RetryTimeout     int    `json:"retry_timeout"`
	RetryDelay       int    `json:"retry_delay"`
	Provider         string `json:"provider"`
	EmailRedisUrl    string `json:"email_redis_url"`
	client           *http.Client
	reqLock          sync.Mutex
	lastReqPhoneTime time.Time
	retryTimeout     time.Duration
	retryDelay       time.Duration
}

func (this *EmailInfo) GetType() verification.VerificationType {
	return verification.TypeEmail
}

func (this *EmailInfo) GetProvider() string {
	return this.Provider
}

func (this *EmailInfo) GetArea() string {
	return this.Domain
}

func InitEmailVerification(email *EmailInfo) (verification.VerificationCodeProvider, error) {
	if email.Provider == "guerrilla" {

	}

	return nil, &common.MakeMoneyError{ErrStr: "unknow provider"}
}
