package verification

import (
	"makemoney/common"
	"makemoney/common/log"
	"sync"
	"time"
)

var GMails []*GMail

type GmailAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Repet    int    `json:"repet"`
}

type GmailConfig struct {
	Accounts     []GmailAccount `json:"accounts"`
	RetryTimeout int            `json:"retry_timeout"`
	RetryDelay   int            `json:"retry_delay"`
}

var GetGmailLock = sync.Mutex{}

func GetGMails() *GMail {
	var g *GMail
	//var index = -1
	for true {
		GetGmailLock.Lock()
		for idx := range GMails {
			if GMails[idx] == nil {
				continue
			}
			//index = idx
			g = GMails[idx]
			GMails[idx] = nil
			break
		}
		GetGmailLock.Unlock()
		if g == nil {
			return nil
		}

		err := g.Login()
		if err != nil {
			log.Error("gmail %s login error: %v", g.Username, err)
			g = nil
			continue
		}
		return g
	}

	return nil
}

func InitDefaultGMail(config *GmailConfig) error {
	index := 0
	GMails = make([]*GMail, len(config.Accounts)*15)
	for _, item := range config.Accounts {
		if item.Repet > 15 {
			item.Repet = 15
		}
		for idx := 0; idx < item.Repet; idx++ {
			GMails[index] = &GMail{
				Username: item.Username,
				Password: item.Password,
				EmailInfo: EmailInfo{
					RetryTimeout:         config.RetryTimeout,
					RetryDelay:           config.RetryDelay,
					RetryTimeoutDuration: time.Duration(config.RetryTimeout) * time.Second,
					RetryDelayDuration:   time.Duration(config.RetryDelay) * time.Second,
				},
			}
			index++
		}
	}
	GMails = GMails[:index]
	if index == 0 {
		return &common.MakeMoneyError{ErrStr: "no gmail!"}
	}
	return nil
}
