package goinsta

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/phone"
	"math/rand"
	"strconv"
	"time"
)

type Register struct {
	inst   *Instagram
	number string
	phone  phone.PhoneVerificationCode
	proxy  *common.Proxy
}

func NewRegister(_proxy *common.Proxy, _phone phone.PhoneVerificationCode) *Register {
	register := &Register{}
	register.phone = _phone
	register.proxy = _proxy
	return register
}

func (this *Register) Do(username string, firstname string, password string) (*Instagram, error) {
	inst, err := this.do(username, firstname, password)
	if err != nil {
		return nil, err
	}

	inst.RegisterIpCountry = this.proxy.Country
	inst.RegisterPhoneArea = this.phone.GetArea()
	inst.RegisterPhoneNumber = this.number
	return inst, err
}

func (this *Register) do(username string, firstname string, password string) (*Instagram, error) {
	this.inst = New(username, password, this.proxy)
	this.inst.PrepareNewClient()

	number, err := this.phone.RequirePhoneNumber()
	if err != nil {
		return nil, err
	}

	log.Info("get phone number: %s", number)
	this.number = number
	err = this.checkPhoneNumber()
	//if err != nil {
	//	return nil, err
	//}
	respSendSignupSmsCode, err := this.sendSignupSmsCode()
	if err != nil {
		return nil, err
	}

	UpdatePhoneSendOnce(this.phone.GetProvider(), this.phone.GetArea(), this.number)
	var flag = false
	defer func() {
		if flag {
			UpdatePhoneRegisterOnce(this.phone.GetArea(), this.number)
		}
	}()

	code, err := this.phone.RequirePhoneCode(number)
	if err != nil {
		return nil, err
	}

	_, err = this.validateSignupSmsCode(code)
	if err != nil {
		return nil, err
	}

	realUsername, err := this.genUsername(username)
	if err != nil {
		return nil, err
	}
	this.inst.User = realUsername

	createValidated, err := this.createValidated(realUsername, firstname, password, code, respSendSignupSmsCode.TosVersion)
	if err != nil {
		return nil, err
	}

	this.inst.IsLogin = true
	this.inst.ID = createValidated.CreatedUser.ID
	flag = true
	return this.inst, err
}

func (this *Register) genUsername(username string) (string, error) {
	usernameSuggestions, err := this.usernameSuggestions(username)
	err = usernameSuggestions.CheckError(err)
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
	return username + fmt.Sprintf("%d%d%d",
		common.GenNumber(1990, 2020),
		common.GenNumber(1, 12),
		common.GenNumber(1, 27)), nil
}

func (this *Register) checkPhoneNumber() error {
	params := map[string]interface{}{
		"phone_id":        this.inst.familyID,
		"login_nonce_map": "{}",
		"phone_number":    this.number,
		"guid":            this.inst.uuid,
		"device_id":       this.inst.androidID,
		"prefill_shown":   "False",
	}

	_, err := this.inst.HttpRequest(&reqOptions{
		ApiPath: urlCheckPhoneNumber,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	})
	return err
}

type RespSendSignupSmsCode struct {
	BaseApiResp
	TosVersion  string `json:"tos_version"`
	AgeRequired bool   `json:"age_required"`
}

func (this *Register) sendSignupSmsCode() (*RespSendSignupSmsCode, error) {
	params := map[string]interface{}{
		"phone_id":           this.inst.familyID,
		"phone_number":       this.phone.GetArea() + this.number,
		"guid":               this.inst.uuid,
		"device_id":          this.inst.androidID,
		"android_build_type": "release",
		"waterfall_id":       this.inst.wid,
	}
	resp := &RespSendSignupSmsCode{}
	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSendSignupSmsCode,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
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
	params := map[string]interface{}{
		"verification_code": code,
		"phone_number":      this.phone.GetArea() + this.number,
		"guid":              this.inst.uuid,
		"device_id":         this.inst.androidID,
		"waterfall_id":      this.inst.wid,
	}
	resp := &RespValidateSignupSmsCode{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlValidateSignupSmsCode,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
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
	params := map[string]interface{}{
		"phone_id":     this.inst.familyID,
		"guid":         this.inst.uuid,
		"name":         username,
		"device_id":    this.inst.androidID,
		"email":        "",
		"waterfall_id": this.inst.wid,
	}
	resp := &RespUsernameSuggestions{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlUsernameSuggestions,
			IsPost:  true,
			Signed:  true,
			Query:   params,
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
	params := map[string]interface{}{
		"_uuid":    this.inst.uuid,
		"username": username,
	}
	resp := &RespCheckUsername{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlCheckUsername,
			IsPost:  true,
			Signed:  true,
			Query:   params,
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
	AccountCreated        bool        `json:"account_created"`
	CreatedUser           UserDetail  `json:"created_user"`
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
	encodePasswd, err := encryptPassword(password, this.inst.ReadHeader(IGHeader_EncryptionId), this.inst.ReadHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	params := map[string]interface{}{
		"is_secondary_account_creation": "false",
		"jazoest":                       genJazoest(this.inst.familyID),
		"tos_version":                   tosVersion,
		"suggestedUsername":             "",
		"verification_code":             code,
		//"sn_result":                              "VERIFICATION_PENDING: request time is " + strconv.FormatInt(time.Now().Unix(), 10),
		"sn_result":                              "API_ERROR: class X.9ob:7: ",
		"do_not_auto_login_if_credentials_match": "true",
		"phone_id":                               this.inst.familyID,
		"enc_password":                           encodePasswd,
		"phone_number":                           this.phone.GetArea() + this.number,
		"username":                               username,
		"first_name":                             firstname,
		"day":                                    strconv.Itoa(rand.Intn(27) + 1),
		"adid":                                   this.inst.adid,
		"guid":                                   this.inst.uuid,
		"year":                                   "2000",
		"device_id":                              this.inst.androidID,
		"_uuid":                                  this.inst.uuid,
		"month":                                  strconv.Itoa(rand.Intn(12) + 1),
		"sn_nonce":                               genSnNonce(this.phone.GetArea() + this.number),
		"force_sign_up_code":                     "",
		"waterfall_id":                           this.inst.wid,
		"qs_stamp":                               "",
		"has_sms_consent":                        "true",
		"one_tap_opt_in":                         "true",
	}
	resp := &RespCreateValidated{}

	err = this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlCreateValidated,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}
