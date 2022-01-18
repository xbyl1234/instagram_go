package emali

import (
	"net/http"
	"sync"
	"time"
)

type EmailInfo struct {
	Domain               string `json:"domain"`
	RetryTimeout         int    `json:"retry_timeout"`
	RetryDelay           int    `json:"retry_delay"`
	Provider             string `json:"provider"`
	EmailMysqlUrl        string `json:"email_mysql_url"`
	client               *http.Client
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
	return this.Domain
}
