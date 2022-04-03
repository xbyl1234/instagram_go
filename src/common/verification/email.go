package verification

import (
	"sync"
	"time"
)

type EmailInfo struct {
	VerificationCodeProvider
	RetryTimeout         int
	RetryDelay           int
	Provider             string
	reqLock              sync.Mutex
	lastReqPhoneTime     time.Time
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

func (this *EmailInfo) GetType() string {
	return "email"
}

func (this *EmailInfo) GetProvider() string {
	return this.Provider
}

func (this *EmailInfo) GetArea() string {
	return ""
}
