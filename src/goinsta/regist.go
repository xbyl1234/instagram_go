package goinsta

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/phone"
	proxy2 "makemoney/common/proxy"
	"math/rand"
	"time"
)

type Register struct {
	inst       *Instagram
	number     string
	phone      phone.PhoneVerificationCode
	proxy      *proxy2.Proxy
	HadSendSMS bool
	HadRecvSMS bool
}

func NewRegister(_proxy *proxy2.Proxy, _phone phone.PhoneVerificationCode) *Register {
	register := &Register{}
	register.phone = _phone
	register.proxy = _proxy
	register.HadSendSMS = false
	register.HadRecvSMS = false
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
	//err = this.checkPhoneNumber()
	//if err != nil {
	//	return nil, err
	//}
	respSendSignupSmsCode, err := this.sendSignupSmsCode()
	if err != nil {
		errRelease := this.phone.ReleasePhone(number)
		if errRelease != nil {
			log.Error("release phone: %s, error: %v", number, errRelease)
		}
		return nil, err
	}
	this.HadSendSMS = true

	_ = UpdatePhoneSendOnce(this.phone.GetProvider(), this.phone.GetArea(), this.number)
	var flag = false
	defer func() {
		if flag {
			_ = UpdatePhoneRegisterOnce(this.phone.GetArea(), this.number)
		}
	}()

	code, err := this.phone.RequirePhoneCode(number)
	if err != nil {
		errRelease := this.phone.ReleasePhone(number)
		if errRelease != nil {
			log.Error("release phone: %s, error: %v", number, errRelease)
		}
		return nil, err
	}
	this.HadRecvSMS = true

	_, err = this.validateSignupSmsCode(code)
	if err != nil {
		return nil, err
	}

	realUsername := this.genUsername(username)
	this.inst.User = realUsername

	year := fmt.Sprintf("%d", common.GenNumber(1995, 2000))
	month := fmt.Sprintf("%d", common.GenNumber(1, 11))
	day := fmt.Sprintf("%d", common.GenNumber(1, 27))
	_, err = this.checkAgeEligibility(year, month, day)
	if err != nil {
		return nil, err
	}

	_, err = this.NewUserFlowBegins()
	_, err = this.checkUsername(realUsername, password)

	createValidated, err := this.createValidated(realUsername, firstname, password, code, respSendSignupSmsCode.TosVersion, year, month, day)
	if err != nil {
		return nil, err
	}
	this.inst.IsLogin = true
	this.inst.ID = createValidated.CreatedUser.ID

	_, err = this.NewAccountNuxSeen()
	_, err = this.GetSteps()
	flag = true
	return this.inst, err
}

func (this *Register) genUsername(username string) string {
	usernameSuggestions, err := this.usernameSuggestions(username)
	if err == nil || len(usernameSuggestions.Suggestions) == 0 {
		return username + fmt.Sprintf("%d%d%d",
			common.GenNumber(1990, 2020),
			common.GenNumber(1, 12),
			common.GenNumber(1, 27))
	} else {
		return usernameSuggestions.Suggestions[0]
	}
}

//func (this *Register) checkPhoneNumber() error {
//	params := map[string]interface{}{
//		"phone_id":        this.inst.familyID,
//		"login_nonce_map": "{}",
//		"phone_number":    this.number,
//		"guid":            this.inst.deviceID,
//		"device_id":       this.inst.androidID,
//		"prefill_shown":   "False",
//	}
//
//	_, err := this.inst.HttpRequest(&reqOptions{
//		ApiPath: urlCheckPhoneNumber,
//		IsPost:  true,
//		Signed:  true,
//		Query:   params,
//	})
//
//	return err
//}

type RespSendSignupSmsCode struct {
	BaseApiResp
	TosVersion  string `json:"tos_version"`
	AgeRequired bool   `json:"age_required"`
}

func (this *Register) sendSignupSmsCode() (*RespSendSignupSmsCode, error) {
	params := map[string]interface{}{
		"device_id":    this.inst.deviceID,
		"phone_number": this.phone.GetArea() + this.number,
		"phone_id":     this.inst.deviceID,
		"source":       "regular",
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

type RespValidateSignupSmsCode struct {
	BaseApiResp
	PnTaken  bool `json:"pn_taken"`
	Verified bool `json:"verified"`
}

func (this *Register) validateSignupSmsCode(code string) (*RespValidateSignupSmsCode, error) {
	params := map[string]interface{}{
		"device_id":         this.inst.deviceID,
		"phone_number":      this.phone.GetArea() + this.number,
		"waterfall_id":      this.inst.wid,
		"verification_code": code,
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
	Suggestions []string `json:"suggestions"`
}

func (this *Register) usernameSuggestions(username string) (*RespUsernameSuggestions, error) {
	params := map[string]interface{}{
		"name":         username,
		"device_id":    this.inst.deviceID,
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

	err = resp.CheckError(err)
	return resp, err
}

type RespCheckAge struct {
	BaseApiResp
	EligibleToRegister      bool `json:"eligible_to_register"`
	ParentalConsentRequired bool `json:"parental_consent_required"`
	IsSupervisedUser        bool `json:"is_supervised_user"`
}

func (this *Register) checkAgeEligibility(year string, month string, day string) (*RespCheckAge, error) {
	params := map[string]interface{}{
		"day":   day,
		"year":  year,
		"month": month,
	}
	resp := &RespCheckAge{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCheckAgeEligibility,
		IsPost:  true,
		Signed:  false,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

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

func (this *Register) checkUsername(username string, password string) (*RespCheckUsername, error) {
	encodePasswd, err := encryptPassword(password, this.inst.ReadHeader(IGHeader_EncryptionId), this.inst.ReadHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}
	params := map[string]interface{}{
		"enc_password": encodePasswd,
		"username":     username,
		"device_id":    this.inst.deviceID,
	}
	resp := &RespCheckUsername{}

	err = this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlCheckUsername,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *Register) NewUserFlowBegins() (*BaseApiResp, error) {
	params := map[string]interface{}{
		"device_id": this.inst.deviceID,
	}
	resp := &BaseApiResp{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlNewUserFlowBegins,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

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
	tosVersion string,
	year string,
	month string,
	day string) (*RespCreateValidated, error) {
	encodePasswd, err := encryptPassword(password, this.inst.ReadHeader(IGHeader_EncryptionId), this.inst.ReadHeader(IGHeader_EncryptionKey))
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	params := map[string]interface{}{
		"tos_version":                            tosVersion,
		"verification_code":                      code,
		"do_not_auto_login_if_credentials_match": "0",
		"phone_id":                               this.inst.deviceID,
		"enc_password":                           encodePasswd,
		"phone_number":                           this.phone.GetArea() + this.number,
		"username":                               username,
		"first_name":                             firstname,
		"day":                                    day,
		"year":                                   year,
		"device_id":                              this.inst.deviceID,
		"month":                                  month,
		"has_seen_aart_on":                       "0",
		"force_create_account":                   "0",
		"waterfall_id":                           this.inst.wid,
		"ck_error":                               "CKErrorDomain: 9",
		"has_sms_consent":                        "true",
		"ck_environment":                         "production",
		"ck_container":                           "iCloud.com.burbn.instagram",
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

func (this *Register) NewAccountNuxSeen() (*BaseApiResp, error) {
	params := map[string]interface{}{
		"is_fb_installed": false,
	}
	resp := &BaseApiResp{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlNewAccountNuxSeen,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *Register) GetSteps() (*BaseApiResp, error) {
	params := map[string]interface{}{
		"device_id":                     this.inst.deviceID,
		"is_secondary_account_creation": "0",
		"push_permission_requested":     "0",
		"network_type":                  "wifi-none",
		"is_account_linking_flow":       "0",
	}
	resp := &BaseApiResp{}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlGetSteps,
			IsPost:  true,
			Signed:  true,
			Query:   params,
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}
