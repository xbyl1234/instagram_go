package phone

import (
	"makemoney/http_helper"
	"makemoney/log"
	"makemoney/tools"
	"strconv"
	"time"
)

//http://h5.do889.com:81/info
//741852
type PhoneDo889 struct {
	PhoneInfo
	RemainCount int
}

type PhoneDo889_Login struct {
	Token string `json:"token"`
	Data  []struct {
		Money   string `json:"money"`
		Money_1 string `json:"money_1"`
		Id      string `json:"id"`
		Leixing string `json:"leixing"`
	} `json:"data"`
}

func (this *PhoneDo889) Login() error {
	var respJson PhoneDo889_Login
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		ReqUrl: this.UrlLogin,
		Params: map[string]string{
			"username": this.Username,
			"password": this.Password},
	}, &respJson)
	if err != nil {
		return err
	}
	this.Token = respJson.Token

	log.Info("phone token %s", respJson.Token)
	return nil
}

//{
//"message": "ok",
//"mobile": "16532643928",
//"data": [],
//"1分钟内剩余取卡数": "298"
//}
type PhoneDo889_RequirePhone struct {
	Message     string        `json:"message"`
	Mobile      string        `json:"mobile"`
	Data        []interface{} `json:"data"`
	RemainCount string        `json:"1分钟内剩余取卡数"`
}

func (this *PhoneDo889) RequirePhoneNumber() (string, error) {
	//http://api.fghfdf.cn/api/get_mobile?token=你的token&project_id=项目ID
	this.reqLock.Lock()
	defer this.reqLock.Unlock()

	if this.RemainCount <= 10 && time.Since(this.lastReqPhoneTime).Minutes() < 1 {
		time.Sleep(time.Minute - time.Since(this.lastReqPhoneTime))
	}

	var respJson PhoneDo889_RequirePhone
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		ReqUrl: this.UrlReqPhoneNumber,
		Params: map[string]string{
			"token":      this.Token,
			"project_id": this.ProjectID},
	}, &respJson)
	this.lastReqPhoneTime = time.Now()
	if err != nil {
		return "", err
	}
	if respJson.Message != "ok" {
		return "", &tools.MakeMoneyError{ErrStr: respJson.Message}
	}
	this.RemainCount, _ = strconv.Atoi(respJson.RemainCount)
	return respJson.Mobile, err
}

type PhoneDo889_ReqPhoneCode struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Data    []struct {
		ProjectId   string `json:""`
		Modle       string `json:"modle"`
		Phone       string `json:"phone"`
		ProjectType string `json:"project_type"`
	} `json:"data"`
}

func (this *PhoneDo889) RequirePhoneCode(number string) (string, error) {
	//http://api.fghfdf.cn/api/get_message
	start := time.Now()
	for ; time.Since(start) < this.retryTimeout; {
		var respJson PhoneDo889_ReqPhoneCode
		err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{ReqUrl: this.UrlReqPhoneCode,
			Params: map[string]string{
				"token":      this.Token,
				"project_id": this.ProjectID,
				"phone_num":  number},
		}, &respJson)

		if err != nil {
			return "", err
		}
		if respJson.Message != "ok" {
			return "", &tools.MakeMoneyError{ErrStr: respJson.Message}
		}
		if respJson.Code != "" {
			return respJson.Code, nil
		}

		time.Sleep(this.retryDelay)
	}

	return "", &tools.MakeMoneyError{ErrStr: "require code timeout"}
}

type PhoneDo889_ReleasePhone struct {
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

func (this *PhoneDo889) ReleasePhone(number string) error {
	//http://api.fghfdf.cn/api/free_mobile
	var respJson PhoneDo889_ReleasePhone
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{
		ReqUrl: this.UrlReqReleasePhone,
		Params: map[string]string{
			"token":     this.Token,
			"phone_num": number},
	}, &respJson)

	if err != nil {
		return err
	}
	if respJson.Message != "ok" {
		return &tools.MakeMoneyError{ErrStr: respJson.Message}
	}
	return nil
}

type PhoneDo889_BlackPhone struct {
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

func (this *PhoneDo889) BlackPhone(number string) error {
	//http://api.fghfdf.cn/api/add_blacklist
	var respJson PhoneDo889_ReleasePhone
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{ReqUrl: this.UrlReqBlackPhone,
		Params: map[string]string{
			"token":      this.Token,
			"project_id": this.ProjectID,
			"phone_num":  number}},
		&respJson)

	if err != nil {
		return err
	}
	if respJson.Message != "ok" {
		return &tools.MakeMoneyError{ErrStr: respJson.Message}
	}
	return nil
}

type PhoneDo889_Balance struct {
	Message string `json:"message"`
	Data    []struct {
		Money   string `json:"money"`
		Money_1 string `json:"money_1"`
	} `json:"data"`
}

func (this *PhoneDo889) GetBalance() (string, error) {
	//http://api.fghfdf.cn/api/add_blacklist
	var respJson PhoneDo889_Balance
	err := http_helper.HttpDoJson(this.client, &http_helper.RequestOpt{ReqUrl: this.UrlReqBalance,
		Params: map[string]string{
			"token": this.Token},
	}, &respJson)

	if err != nil {
		return "", err
	}
	if respJson.Message != "ok" {
		return "", &tools.MakeMoneyError{ErrStr: respJson.Message}
	}
	return respJson.Data[0].Money, nil
}