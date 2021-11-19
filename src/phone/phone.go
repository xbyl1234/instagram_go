package phone

import (
	"makemoney/tools"
	"net/http"
	"sync"
	"time"
)

type PhoneVerificationCode interface {
	RequirePhoneNumber() (string, error)
	RequirePhoneCode(number string) (string, error)
	ReleasePhone(number string) error
	BlackPhone(number string) error
	GetBalance() (string, error)
	GetProvider() string
	GetArea() string
	Login() error
}

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
	client             *http.Client
	reqLock            sync.Mutex
	lastReqPhoneTime   time.Time
	retryTimeout       time.Duration
	retryDelay         time.Duration
}

func (this *PhoneInfo) GetProvider() string {
	return this.Provider
}

func (this *PhoneInfo) GetArea() string {
	return this.Area
}

func NewPhoneVerificationCode(provider string) (PhoneVerificationCode, error) {
	if provider == "do889" {
		ret := &PhoneDo889{}
		err := tools.LoadJsonFile("./config/phone.json", ret)
		if err == nil {
			ret.retryDelay = time.Duration(ret.RetryDelay) * time.Second
			ret.retryTimeout = time.Duration(ret.RetryTimeout) * time.Second
			ret.client = &http.Client{}
			//ret.reqLock = &sync.Mutex{}
			tools.DebugHttpClient(ret.client)
		}
		return ret, err
	}
	return nil, nil
}
