package goinsta

import (
	"makemoney/goinsta/dbhelper"
	"makemoney/phone"
	"math/rand"
	"strconv"
	"time"
)

type Register struct {
	inst   *Instagram
	area   string
	number string
	phone  phone.PhoneVerificationCode
}

func NewRegister(area string, _phone phone.PhoneVerificationCode) *Register {
	register := &Register{}
	register.phone = _phone
	return register
}

func (this *Register) Do(username string, firstname string, password string) (*Instagram, error) {
	inst, err := this.do(username, firstname, password)
	return inst, err
}

func (this *Register) do(username string, firstname string, password string) (*Instagram, error) {
	this.inst = New(username, password)
	err := this.inst.Prepare()
	if err != nil {
		return nil, err
	}

	number, err := this.phone.RequirePhoneNumber()
	if err != nil {
		return nil, err
	}

	this.number = number
	err = this.checkPhoneNumber()
	//if err != nil {
	//	return nil, err
	//}
	respSendSignupSmsCode, err := this.sendSignupSmsCode()
	err = CheckApiError(respSendSignupSmsCode, err)
	if err != nil {
		return nil, err
	}
	dbhelper.UpdatePhoneSendOnce(this.phone.GetProvider(), this.area, this.number)
	var flag = false
	defer func() {
		if flag {
			dbhelper.UpdatePhoneRegisterOnce(this.area, this.number)
		}
	}()

	code, err := this.phone.RequirePhoneCode(number)
	if err != nil {
		return nil, err
	}

	validateSignupSmsCode, err := this.validateSignupSmsCode(code)
	err = CheckApiError(validateSignupSmsCode, err)
	if err != nil {
		return nil, err
	}

	realUsername, err := this.genUsername(username)
	if err != nil {
		return nil, err
	}
	this.inst.User = username

	createValidated, err := this.createValidated(realUsername, firstname, password, code, respSendSignupSmsCode.TosVersion)
	err = CheckApiError(createValidated, err)
	if err != nil {
		return nil, err
	}

	flag = true
	return this.inst, err
}

func (this *Register) genUsername(username string) (string, error) {
	usernameSuggestions, err := this.usernameSuggestions(username)
	err = CheckApiError(usernameSuggestions, err)
	if err != nil {
		return "", err
	}
	if usernameSuggestions.SuggestionsWithMetadata.Suggestions != nil {
		for sugNameIdx := range usernameSuggestions.SuggestionsWithMetadata.Suggestions {
			sugName := usernameSuggestions.SuggestionsWithMetadata.Suggestions[sugNameIdx].Username
			checkUsername, err := this.checkUsername(sugName)
			if err != nil {
				return "", err
			}
			if checkUsername.Available {
				return sugName, nil
			}
		}

	}

	return "", &ApiError{"not find available username!"}
}

func (this *Register) checkPhoneNumber() error {
	params := map[string]string{
		"phone_id":        this.inst.pid,
		"login_nonce_map": "{}",
		"phone_number":    this.number,
		"guid":            this.inst.uuid,
		"device_id":       this.inst.dID,
		"prefill_shown":   "False",
	}

	err := this.inst.SendRequest(&reqOptions{
		Endpoint: urlCheckPhoneNumber,
		IsPost:   true,
		Signed:   true,
		Query:    params,
	}, nil)
	return err
}

type RespSendSignupSmsCode struct {
	BaseApiResp
	TosVersion  string `json:"tos_version"`
	AgeRequired bool   `json:"age_required"`
}

func (this *Register) sendSignupSmsCode() (*RespSendSignupSmsCode, error) {
	params := map[string]string{
		"phone_id":           this.inst.pid,
		"phone_number":       this.area + this.number,
		"guid":               this.inst.uuid,
		"device_id":          this.inst.dID,
		"android_build_type": "release",
		"waterfall_id":       this.inst.wid,
	}
	resp := &RespSendSignupSmsCode{}
	err := this.inst.SendRequest(
		&reqOptions{
			Endpoint: urlZrToken,
			IsPost:   true,
			Signed:   true,
			Query:    params,
		}, resp)

	return resp, err
}

type RespValidateSignupSmsCodeError struct {
	ErrorType string `json:"error_type"`
	Errors    struct {
		Nonce []string `json:"nonce"`
	} `json:"errors"`
}

type RespValidateSignupSmsCode struct {
	BaseApiResp
	RespValidateSignupSmsCodeError
	NonceValid bool `json:"nonce_valid"`
	Verified   bool `json:"verified"`
}

func (this *Register) validateSignupSmsCode(code string) (*RespValidateSignupSmsCode, error) {
	params := map[string]string{
		"verification_code": code,
		"phone_number":      this.area + this.number,
		"guid":              this.inst.uuid,
		"device_id":         this.inst.dID,
		"waterfall_id":      this.inst.wid,
	}
	resp := &RespValidateSignupSmsCode{}

	err := this.inst.SendRequest(
		&reqOptions{
			Endpoint: urlZrToken,
			IsPost:   true,
			Signed:   true,
			Query:    params,
		}, resp)

	return resp, err
}

type RespUsernameSuggestions struct {
	BaseApiResp
	SuggestionsWithMetadata struct {
		Suggestions []struct {
			Prototype string `json:"prototype"`
			Username  string `json:"username"`
		} `json:"suggestions"`
	} `json:"suggestions_with_metadata"`
}

func (this *Register) usernameSuggestions(username string) (*RespUsernameSuggestions, error) {
	params := map[string]string{
		"phone_id":     this.inst.pid,
		"guid":         this.inst.uuid,
		"name":         username,
		"device_id":    this.inst.dID,
		"email":        "",
		"waterfall_id": this.inst.wid,
	}
	resp := &RespUsernameSuggestions{}

	err := this.inst.SendRequest(
		&reqOptions{
			Endpoint: urlUsernameSuggestions,
			IsPost:   true,
			Signed:   true,
			Query:    params,
		}, resp)

	return resp, err
}

//{
//	"username": "zha",
//	"available": false,
//	"existing_user_password": false,
//	"error": "帐号 zha 不可用",
//	"status": "ok",
//	"error_type": "username_is_taken"
//}
//{
//	"username": "zhanghao7549",
//	"available": true,
//	"existing_user_password": false,
//	"status": "ok"
//}
type RespCheckUsernameError struct {
	Error string `json:"error"`
}

type RespCheckUsername struct {
	BaseApiResp
	RespCheckUsernameError
	Username             string `json:"username"`
	Available            bool   `json:"available"`
	ExistingUserPassword bool   `json:"existing_user_password"`
}

func (this *Register) checkUsername(username string) (*RespCheckUsername, error) {
	params := map[string]string{
		"_uuid":    this.inst.uuid,
		"username": username,
	}
	resp := &RespCheckUsername{}

	err := this.inst.SendRequest(
		&reqOptions{
			Endpoint: urlUsernameSuggestions,
			IsPost:   true,
			Signed:   true,
			Query:    params,
		}, resp)

	return resp, err
}

//{
//    "account_created": false,
//    "allow_contacts_sync": true,
//    "error_type": "username_is_taken",
//    "errors": {
//        "username": [
//            "\u8fd9\u4e2a\u5e10\u53f7\u7528\u4e0d\u4e86\uff0c\u6362\u4e00\u4e2a\u8bd5\u8bd5\u5457\u3002"
//        ]
//    },
//    "existing_user": false,
//    "status": "ok"
//}

type RespCreateValidatedError struct {
	AllowContactsSync bool   `json:"allow_contacts_sync"`
	ErrorType         string `json:"error_type"`
	Errors            struct {
		Username []string `json:"username"`
	} `json:"errors"`
}

type RespCreateValidated struct {
	BaseApiResp
	RespCreateValidatedError
	AccountCreated        string      `json:"account_created"`
	CreatedUser           Account     `json:"created_user"`
	ExistingUser          bool        `json:"existing_user"`
	MultipleUsersOnDevice bool        `json:"multiple_users_on_device"`
	SessionFlushNonce     interface{} `json:"session_flush_nonce"`
}

func (this *Register) createValidated(
	username string,
	firstname string,
	password string,
	code string,
	tosVersion string) (*RespCreateValidated, error) {

	rand.Seed(time.Now().UnixNano())
	params := map[string]string{
		"is_secondary_account_creation":          "false",
		"jazoest":                                genJazoest(this.inst.pid),
		"tos_version":                            tosVersion,
		"suggestedUsername":                      "",
		"verification_code":                      code,
		"sn_result":                              "API_ERROR: class X.9ob:7: ",
		"do_not_auto_login_if_credentials_match": "true",
		"phone_id":                               this.inst.pid,
		"enc_password":                           password,
		"phone_number":                           this.area + this.number,
		"username":                               username,
		"first_name":                             firstname,
		"day":                                    strconv.Itoa(rand.Intn(27) + 1),
		"adid":                                   this.inst.adid,
		"guid":                                   this.inst.uuid,
		"year":                                   "2000",
		"device_id":                              "android-79c028b2e54c371e",
		"_uuid":                                  "df00cccf-3663-412d-9145-585a4a833ce3",
		"month":                                  strconv.Itoa(rand.Intn(12) + 1),
		"sn_nonce":                               genSnNonce(this.area + this.number),
		"force_sign_up_code":                     "",
		"waterfall_id":                           this.inst.wid,
		"qs_stamp":                               "",
		"has_sms_consent":                        "true",
		"one_tap_opt_in":                         "true",
	}
	resp := &RespCreateValidated{}

	err := this.inst.SendRequest(
		&reqOptions{
			Endpoint: urlUsernameSuggestions,
			IsPost:   true,
			Signed:   true,
			Query:    params,
		}, resp)

	return resp, err
}
