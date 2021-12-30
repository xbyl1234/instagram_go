package phone

import (
	"makemoney/common"
	"net/http"
	"strings"
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
	City               string `json:"city"`

	client           *http.Client
	reqLock          sync.Mutex
	lastReqPhoneTime time.Time
	retryTimeout     time.Duration
	retryDelay       time.Duration
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
		err := common.LoadJsonFile("./config/phone.json", ret)
		if err == nil {
			ret.retryDelay = time.Duration(ret.RetryDelay) * time.Second
			ret.retryTimeout = time.Duration(ret.RetryTimeout) * time.Second
			ret.client = &http.Client{}
			//ret.reqLock = &sync.Mutex{}
			//common.DebugHttpClient(ret.client)
		}
		return ret, err
	} else if provider == "taxin" {
		ret := &PhoneTaxin{}
		err := common.LoadJsonFile("./config/phone_taxin.json", ret)
		if err == nil {
			ret.retryDelay = time.Duration(ret.RetryDelay) * time.Second
			ret.retryTimeout = time.Duration(ret.RetryTimeout) * time.Second
			ret.client = &http.Client{}
			//ret.reqLock = &sync.Mutex{}
			//common.DebugHttpClient(ret.client)

			if ret.Token == "" {
				err = ret.Login()
				if err != nil {
					return nil, err
				}
				common.Dumps("./config/phone_taxin.json", ret)
			}
		}
		return ret, err
	}
	return nil, nil
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
