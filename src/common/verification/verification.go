package verification

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

func InitTaxin(phoneInfo *PhoneInfo) (VerificationCodeProvider, error) {
	ret := &PhoneTaxin{}
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

type GuerrillaConfig struct {
	RetryTimeout int    `json:"retry_timeout"`
	RetryDelay   int    `json:"retry_delay"`
	Domain       string `json:"domain"`
	MysqlUrl     string `json:"mysql_url"`
}

func InitGuerrilla(config *GuerrillaConfig) (VerificationCodeProvider, error) {
	var err error
	ret := &Guerrilla{
		MysqlDB:  nil,
		Domain:   config.Domain,
		MysqlUrl: config.MysqlUrl,
		EmailInfo: EmailInfo{
			RetryDelayDuration:   time.Duration(config.RetryDelay) * time.Second,
			RetryTimeoutDuration: time.Duration(config.RetryTimeout) * time.Second,
		},
	}

	ret.MysqlDB, err = sqlx.Connect("mysql", ret.MysqlUrl)
	if err != nil {
		return nil, err
	}
	return ret, err
}
