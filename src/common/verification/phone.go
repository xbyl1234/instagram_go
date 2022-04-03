package verification

import (
	"net/http"
	"sync"
	"time"
)

type PhoneInfo struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	UrlLogin           string `json:"url_login"`
	UrlReqPhoneNumber  string `json:"url_req_phone_number"`
	UrlReqPhoneCode    string `json:"url_req_phone_code"`
	UrlReqReleasePhone string `json:"url_req_release_phone"`
	UrlReqBlackPhone   string `json:"url_req_black_phone"`
	UrlReqBalance      string `json:"url_req_balance"`
	Token              string `json:"token"`
	ProjectID          string `json:"project_id"`
	RetryTimeout       int    `json:"retry_timeout"`
	RetryDelay         int    `json:"retry_delay"`
	Provider           string `json:"provider"`
	Area               string `json:"area"`
	City               string `json:"city"`

	Client               *http.Client
	reqLock              sync.Mutex
	lastReqPhoneTime     time.Time
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

func (this *PhoneInfo) GetType() string {
	return "phone"
}

func (this *PhoneInfo) GetProvider() string {
	return this.Provider
}

func (this *PhoneInfo) GetArea() string {
	return this.Area
}
