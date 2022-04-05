package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
)

var UsePanic = true
var (
	InsAccountError_ChallengeRequired = "challenge_required"
	InsAccountError_LoginRequired     = "login_required"
	InsAccountError_Feedback          = "feedback_required"
	InsAccountError_RateLimitError    = "rate_limit_error"
)

type BaseApiResp struct {
	url  string
	inst *Instagram

	Status     string `json:"status"`
	ErrorType  string `json:"error_type"`
	Message    string `json:"message"`
	ErrorTitle string `json:"error_title"`
}

func (this *BaseApiResp) SetInfo(url string, inst *Instagram) {
	this.inst = inst
	this.url = url
}

func (this *BaseApiResp) isError() bool {
	return this.Status != "ok"
}

func (this *BaseApiResp) CheckError(err error) error {
	if err != nil {
		return err
	}
	if this == nil {
		return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.OtherError}
	}
	if this.Status != "ok" {
		log.Warn("account: %s, url: %s, api error: %s",
			this.inst.User,
			this.url,
			this.ErrorType+":"+this.Message)
		if UsePanic {
			switch this.Message {
			case InsAccountError_ChallengeRequired, InsAccountError_LoginRequired, InsAccountError_Feedback:
				this.inst.Status = this.Message
				panic(&common.MakeMoneyError{ErrStr: this.Message})
			default:
				return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ApiError}
			}
		} else {
			switch this.Message {
			case InsAccountError_ChallengeRequired:
				this.inst.Status = InsAccountError_ChallengeRequired
				return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ChallengeRequiredError}
			case InsAccountError_LoginRequired:
				this.inst.Status = InsAccountError_LoginRequired
				return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.LoginRequiredError}
			case InsAccountError_Feedback:
				this.inst.Status = InsAccountError_Feedback
				return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.FeedbackError}
			default:
				return &common.MakeMoneyError{ErrStr: this.Message, ErrType: common.ApiError}
			}
		}
	}
	return nil
}
