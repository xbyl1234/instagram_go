package verification

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"net/mail"
	"strings"
	"time"
)

type GMail struct {
	Username string `json:"username"`
	Password string `json:"password"`
	client   *client.Client

	EmailInfo
}

func (this *GMail) RequireAccount() (string, error) {
	username := strings.ReplaceAll(this.Username, "@gmail.com", "")
	return username + "+" + common.GenString(common.CharSet_All, 5) + fmt.Sprintf("%d", time.Now().Unix()) + "@gmail.com", nil
}

type RespEmail struct {
	Header mail.Header
	Body   []byte
}

func (this *GMail) RequireCode(email string) (string, error) {
	this.reqLock.Lock()
	defer this.reqLock.Unlock()

	start := time.Now()
	for time.Since(start) < this.RetryTimeoutDuration {
		_, err := this.client.Select("INBOX", false)
		if err != nil {
			return "", err
		}
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}
		criteria.Text = []string{email}
		ids, err := this.client.Search(criteria)
		if err != nil {
			return "", err
		}
		if len(ids) > 0 {
			seqset := &imap.SeqSet{}
			seqset.AddNum(ids[0])
			messages := make(chan *imap.Message, 2)
			section := &imap.BodySectionName{}
			err = this.client.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, messages)
			if err != nil {
				return "", err
			}
			for msg := range messages {
				r := msg.GetBody(section)
				m, err := mail.ReadMessage(r)
				if err != nil {
					return "", err
				}
				body, err := ioutil.ReadAll(m.Body)
				if err != nil {
					return "", err
				}
				//return &RespEmail{
				//	Header: m.Header,
				//	Body:   body,
				//}, nil

				bodyStr := strings.ReplaceAll(string(body), "=\r\n", "")
				tag := "padding-bottom:25px;\">"
				p := strings.Index(bodyStr, tag) + len(tag)
				if p == -1 {
					return "", &common.MakeMoneyError{ErrStr: "not find code"}
				}
				code := bodyStr[p : p+6]
				return code, nil
			}
		}

		log.Warn("wait for email %s code...", email)
		time.Sleep(this.RetryDelayDuration)
	}

	return "", &common.MakeMoneyError{ErrStr: "require code timeout", ErrType: common.RecvPhoneCodeError}
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
