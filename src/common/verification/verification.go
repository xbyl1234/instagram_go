package verification

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/verification/emali"
	"makemoney/common/verification/phone"
	"net/http"
	"time"
)

type VerificationCodeProvider interface {
	RequireAccount() (string, error)
	RequireCode(number string) (string, error)
	ReleaseAccount(number string) error
	BlackAccount(number string) error
	GetBalance() (string, error)
	GetProvider() string
	GetArea() string
	Login() error
	GetType() string
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
	VerificationProvider = make(map[string]VerificationCodeProvider)
	for _, item := range config {
		var Provider VerificationCodeProvider
		if item.ProviderType == "phone" {
			Provider, err = InitPhoneVerification(&item.Phone)
		} else {
			Provider, err = InitEmailVerification(&item.Email)
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

func InitPhoneVerification(phoneInfo *phone.PhoneInfo) (VerificationCodeProvider, error) {
	//if provider == "do889" {
	//	ret := &PhoneDo889{}
	//	err := common.LoadJsonFile("./config/phone.json", ret)
	//	if err == nil {
	//		ret.retryDelay = time.Duration(ret.RetryDelay) * time.Second
	//		ret.retryTimeout = time.Duration(ret.RetryTimeout) * time.Second
	//		ret.client = &http.Client{}
	//		//ret.reqLock = &sync.Mutex{}
	//		//common.DebugHttpClient(ret.client)
	//	}
	//	return ret, err
	//} else
	if phoneInfo.Provider == "taxin" {
		ret := &phone.PhoneTaxin{}
		ret.PhoneInfo = phoneInfo
		ret.RetryDelayDuration = time.Duration(ret.RetryDelay) * time.Second
		ret.RetryTimeoutDuration = time.Duration(ret.RetryTimeout) * time.Second
		ret.Client = &http.Client{}

		var err error
		if ret.Token == "" {
			err = ret.Login()
			if err != nil {
				return nil, err
			}
			//common.Dumps("./config/phone_taxin.json", ret)
		}
		return ret, err
	}
	return nil, nil
}

func InitEmailVerification(emailInfo *emali.EmailInfo) (VerificationCodeProvider, error) {
	var err error
	if emailInfo.Provider == "guerrilla" {
		ret := &emali.Guerrilla{}
		ret.EmailInfo = emailInfo
		ret.RetryDelayDuration = time.Duration(ret.RetryDelay) * time.Second
		ret.RetryTimeoutDuration = time.Duration(ret.RetryTimeout) * time.Second

		ret.MysqlDB, err = sqlx.Connect("mysql", ret.EmailMysqlUrl)
		if err != nil {
			return nil, err
		}
		return ret, err
	}

	return nil, &common.MakeMoneyError{ErrStr: "unknow provider"}
}
