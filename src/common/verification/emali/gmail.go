package emali

import (
	"fmt"
	"github.com/emersion/go-imap"
	_ "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-imap/client"
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"net/mail"
	"strings"
	"sync"
	"time"
)

type GMail struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	RetryTimeout         int    `json:"retry_timeout"`
	RetryDelay           int    `json:"retry_delay"`
	client               *client.Client
	reqLock              sync.Mutex
	RetryTimeoutDuration time.Duration
	RetryDelayDuration   time.Duration
}

func (this *GMail) RequireAccount() (string, error) {
	username := strings.ReplaceAll(this.Username, "@gmail.com", "")
	return username + "+" + common.GenString(common.CharSet_All, 5) + fmt.Sprintf("%d", time.Now().Unix()) + "@gmail.com", nil
}

type RespEmail struct {
	Header mail.Header
	Body   []byte
}

func (this *GMail) RequireCode(email string) (*RespEmail, error) {
	this.reqLock.Lock()
	defer this.reqLock.Unlock()

	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		_, err := this.client.Select("INBOX", false)
		if err != nil {
			return nil, err
		}
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}
		criteria.Text = []string{email}
		ids, err := this.client.Search(criteria)
		if err != nil {
			return nil, err
		}
		if len(ids) > 0 {
			seqset := &imap.SeqSet{}
			seqset.AddNum(ids[0])
			messages := make(chan *imap.Message, 2)
			section := &imap.BodySectionName{}
			err = this.client.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
			if err != nil {
				return nil, err
			}
			for msg := range messages {
				r := msg.GetBody(section)
				m, err := mail.ReadMessage(r)
				if err != nil {
					return nil, err
				}
				body, err := ioutil.ReadAll(m.Body)
				if err != nil {
					return nil, err
				}
				return &RespEmail{
					Header: m.Header,
					Body:   body,
				}, nil
			}
		}

		log.Warn("wait for email %s code...", email)
		time.Sleep(this.RetryDelayDuration)
	}

	return nil, &common.MakeMoneyError{ErrStr: "require code timeout", ErrType: common.RecvPhoneCodeError}
}

func (this *GMail) ReleaseAccount(number string) error {
	return nil
}

func (this *GMail) BlackAccount(number string) error {
	return nil
}

func (this *GMail) GetBalance() (string, error) {
	return "", nil
}

func (this *GMail) Login() error {
	var err error
	this.client, err = client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return err
	}
	err = this.client.Login(this.Username, this.Password)
	if err != nil {
		return err
	}
	return nil
}
