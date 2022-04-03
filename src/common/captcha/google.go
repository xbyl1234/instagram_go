package captcha

import (
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"time"
)

type GoogleCaptchaConfig struct {
	ClientKey    string `json:"clientKey"`
	WebsiteURL   string `json:"websiteURL"`
	WebsiteKey   string `json:"websiteKey"`
	Type         string `json:"type"`
	RetryTimeout int    `json:"retry_timeout"`
	RetryDelay   int    `json:"retry_delay"`
}

type GoogleCaptcha struct {
	ClientKey string `json:"clientKey"`
	Task      struct {
		WebsiteURL          string      `json:"websiteURL"`
		WebsiteKey          string      `json:"websiteKey"`
		WebsiteSToken       interface{} `json:"websiteSToken"`
		RecaptchaDataSValue interface{} `json:"recaptchaDataSValue"`
		Type                string      `json:"type"`
	} `json:"task"`
	SoftId               int `json:"softId"`
	client               *http.Client
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

var Google *GoogleCaptcha

func InitDefaultGoogleCaptcha(config *GoogleCaptchaConfig) *GoogleCaptcha {
	Google = NewGoogleCaptcha(config)
	return Google
}

func NewGoogleCaptcha(config *GoogleCaptchaConfig) *GoogleCaptcha {
	g := &GoogleCaptcha{
		ClientKey: config.ClientKey,
		Task: struct {
			WebsiteURL          string      `json:"websiteURL"`
			WebsiteKey          string      `json:"websiteKey"`
			WebsiteSToken       interface{} `json:"websiteSToken"`
			RecaptchaDataSValue interface{} `json:"recaptchaDataSValue"`
			Type                string      `json:"type"`
		}{
			WebsiteURL:          config.WebsiteURL,
			WebsiteKey:          config.WebsiteKey,
			WebsiteSToken:       nil,
			RecaptchaDataSValue: nil,
			Type:                config.Type,
		},
		RetryTimeoutDuration: time.Duration(config.RetryTimeout) * time.Second,
		RetryDelayDuration:   time.Duration(config.RetryDelay) * time.Second,
		SoftId:               802,
		client:               &http.Client{},
	}
	return g
}

type RespCreateTask struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	TaskId           string `json:"taskId"`
}

func (this *GoogleCaptcha) CreateTask() (string, error) {
	resp := &RespCreateTask{}
	err := common.HttpDoJson(this.client, &common.RequestOpt{
		IsPost:   true,
		ReqUrl:   "https://api.yescaptcha.com/createTask",
		JsonData: this,
	}, resp)
	if err != nil {
		return "", err
	}
	if resp.ErrorCode != "" {
		return "", &common.MakeMoneyError{ErrStr: resp.ErrorDescription}
	}
	return resp.TaskId, nil
}

type ReqTaskResult struct {
	ClientKey   string `json:"clientKey"`
	TaskId      string `json:"taskId"`
	CacheRecord string `json:"cacheRecord,omitempty"`
}

type RespTaskResult struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Solution         struct {
		GRecaptchaResponse string `json:"gRecaptchaResponse"`
	} `json:"solution"`
	Status string `json:"status"`
}

func (this *GoogleCaptcha) GetTaskResult(taskId string) (string, error) {
	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		result, err := this.getTaskResult(taskId)
		if err != nil {
			log.Error("req google code error: %v", err)
		} else {
			if result.Status == "ready" {
				return result.Solution.GRecaptchaResponse, nil
			}
			if result.Status != "processing" {
				log.Error("req google code error: %s", result.ErrorDescription)
			}
		}
		log.Warn("wait for google code...")
		time.Sleep(this.RetryDelayDuration)
	}
	return "", &common.MakeMoneyError{ErrStr: "require google code timeout", ErrType: common.RecvPhoneCodeError}
}

func (this *GoogleCaptcha) getTaskResult(taskId string) (*RespTaskResult, error) {
	resp := &RespTaskResult{}
	err := common.HttpDoJson(this.client, &common.RequestOpt{
		IsPost: true,
		ReqUrl: "https://api.yescaptcha.com/getTaskResult",
		JsonData: &ReqTaskResult{
			ClientKey: this.ClientKey,
			TaskId:    taskId,
			//CacheRecord: "",
		},
	}, resp)
	return resp, err
}
