package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"makemoney/common"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

var ProxyCallBack func(country string, id string) (*common.Proxy, error)

type Operation struct {
	OperName string    `json:"oper_name"`
	NextTime time.Time `json:"next_time"`
}

type OperationLog struct {
	OperName  string    `json:"oper_name"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Count     int       `json:"count"`
}

type InstagramOperate struct {
	Graph       *Graph
	Upload      *Upload
	Account     *Account
	Message     *Message
	UserOperate *UserOperate
	VideoFeed   *VideoFeed
}

type Instagram struct {
	User         string
	Pass         string
	token        string
	challengeURL string
	ID           int64
	httpHeader   map[string]string
	IsLogin      bool
	AccountInfo  *InstAccountInfo
	sessionID    string

	Status       string
	Tags         string
	Proxy        *common.Proxy
	c            *http.Client
	SpeedControl map[string]*SpeedControl
	Operate      InstagramOperate
	MatePoint    interface{}
}

func (this *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(InstagramHost)
	if err != nil {
		return err
	}

	cookies := this.c.Jar.Cookies(url)
	this.c.Jar = jar
	this.c.Jar.SetCookies(url, cookies)
	return nil
}

func New(username, password string, _proxy *common.Proxy) *Instagram {
	c := common.CreateGoHttpClient(common.DefaultHttpTimeout(), common.NeedJar(), _proxy.GetProxy())
	inst := &Instagram{
		User: username,
		Pass: password,

		sessionID: strings.ToUpper(common.GenUUID()),
		Proxy:     _proxy,
		c:         c,
	}

	inst.AccountInfo = GenInstDeviceInfo()
	inst.Operate.Graph = &Graph{inst: inst}
	inst.httpHeader = make(map[string]string)
	inst.SpeedControl = make(map[string]*SpeedControl)
	return inst
}

func (this *Instagram) GetGraph() *Graph {
	return this.Operate.Graph
}

func (this *Instagram) NewSearch(q string) *Search {
	return newSearch(this, q)
}

func (this *Instagram) NewTagsFeed(tagName string, tab string) *TagsFeed {
	return NewTagsFeed(this, tagName, tab)
}

func (this *Instagram) GetUpload() *Upload {
	if this.Operate.Upload == nil {
		this.Operate.Upload = newUpload(this)
	}
	return this.Operate.Upload
}

func (this *Instagram) GetUserOperate() *UserOperate {
	if this.Operate.UserOperate == nil {
		this.Operate.UserOperate = newUserOperate(this)
	}
	return this.Operate.UserOperate
}

func (this *Instagram) GetAccount() *Account {
	if this.Operate.Account == nil {
		this.Operate.Account = &Account{ID: this.ID, inst: this}
	}
	return this.Operate.Account
}

func (this *Instagram) GetVideoFeed() *VideoFeed {
	if this.Operate.VideoFeed == nil {
		this.Operate.VideoFeed = newVideoFeed(this)
	}
	return this.Operate.VideoFeed
}

func (this *Instagram) NewComments(id string) *Comments {
	return newComments(this, id)
}

func (this *Instagram) NewUser(id string) *User {
	pk, _ := strconv.ParseInt(id, 10, 64)
	return &User{ID: pk, inst: this}
}

func (this *Instagram) NewFollowers(id string) *Followers {
	pk, _ := strconv.ParseInt(id, 10, 64)
	return &Followers{User: pk, inst: this, HasMore: true}
}

func (this *Instagram) GetMessage() *Message {
	if this.Operate.Message == nil {
		this.Operate.Message = newMessage(this)
	}
	return this.Operate.Message
}

func (this *Instagram) SetProxy(_proxy *common.Proxy) {
	this.Proxy = _proxy
	_proxy.GetProxy()(this.c)
}

func (this *Instagram) ResetProxy() {
	this.Proxy, _ = ProxyCallBack(this.AccountInfo.Register.RegisterIpCountry, "")
	this.Proxy.GetProxy()(this.c)
}

func (this *Instagram) IsBad() bool {
	if this.Status == InsAccountError_ChallengeRequired ||
		this.Status == InsAccountError_Feedback ||
		this.Status == InsAccountError_LoginRequired {
		return true
	}
	return false
}

func (this *Instagram) InitSpeedControl(OperName string) *SpeedControl {
	sc := this.SpeedControl[OperName]
	//var isCtrl = true
	if sc == nil {
		sc, _ = GetSpeedControl(OperName)
		this.SpeedControl[OperName] = sc
	}
	return sc
}

func (this *Instagram) IsSpeedLimit(OperName string) bool {
	sc := this.InitSpeedControl(OperName)
	return sc.IsSpeedLimit()
}

func (this *Instagram) IsSpeedLimitInDay(OperName string) bool {
	sc := this.InitSpeedControl(OperName)
	return sc.IsSpeedLimitInDay()
}

func (this *Instagram) GetCoolTime(OperName string) time.Duration {
	sc := this.InitSpeedControl(OperName)
	return sc.GetCoolTime()
}

func (this *Instagram) Increase(OperName string) (int, int, int, int) {
	sc := this.InitSpeedControl(OperName)
	return sc.Increase()
}

func (this *Instagram) IncreaseSuccess(OperName string) {
	sc := this.InitSpeedControl(OperName)
	sc.IncreaseSuccess()
}

func (this *Instagram) GetSpeed(OperName string) (int, int, int, int) {
	sc := this.InitSpeedControl(OperName)
	return sc.GetSpeed()
}

func (this *Instagram) CleanCookiesAndHeader() {
	this.httpHeader = make(map[string]string)
	this.c.Jar, _ = cookiejar.New(nil)
	this.IsLogin = false
}

func (this *Instagram) GetHeader(key string) string {
	return this.httpHeader[key]
}

type PrepareResult struct {
	LauncherSync   bool `json:"launcher_sync"`
	GetNamePrefill bool `json:"get_name_prefill"`
	ContactPrefill bool `json:"contact_prefill"`
	QeSync         bool `json:"qe_sync"`
	LogAttribution bool `json:"log_attribution"`
}

func (this *Instagram) PrepareNewClient() PrepareResult {
	var ret PrepareResult
	ret.LauncherSync = this.launcherSync() == nil
	ret.GetNamePrefill = this.getNamePrefill() == nil
	ret.ContactPrefill = this.contactPrefill() == nil
	ret.QeSync = this.QeSync() == nil
	ret.LogAttribution = this.logAttribution() == nil
	return ret
}

func (this *Instagram) AfterLogin() {
	_ = this.QeSync()
	_ = this.launcherSync()
}

func (this *Instagram) getNamePrefill() error {
	var query = map[string]interface{}{
		"phone_id":  this.AccountInfo.Device.DeviceID,
		"device_id": this.AccountInfo.Device.DeviceID,
	}

	resp := &BaseApiResp{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlGetNamePrefill,
			Query:   query,
			Header: map[string]string{
				"X-Ig-Connection-Speed": "-1kbps",
			},
			IsPost: true,
			Signed: true,
		}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *Instagram) contactPrefill() error {
	var query = map[string]interface{}{
		"phone_id": this.AccountInfo.Device.DeviceID,
	}
	resp := &BaseApiResp{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlContactPrefill,
			Query:   query,
			Header: map[string]string{
				"X-Ig-Connection-Speed": "-1kbps",
			},
			IsPost: true,
			Signed: true,
		}, resp)

	err = resp.CheckError(err)
	return err
}

type RespQeSync struct {
	BaseApiResp
	Experiments []struct {
		Name             string        `json:"name"`
		Group            string        `json:"group"`
		AdditionalParams []interface{} `json:"additional_params"`
		Params           []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"params"`
		LoggingId string `json:"logging_id,omitempty"`
	} `json:"experiments"`
	NoClose bool `json:"no_close"`
}

func (this *Instagram) QeSync() error {
	query := &struct {
		Id                    string `json:"id,omitempty"`
		Uuid                  string `json:"_uuid,omitempty"`
		Uid                   string `json:"_uid,omitempty"`
		ServerConfigRetrieval string `json:"server_config_retrieval,omitempty"`
	}{}

	if this.IsLogin {
		query.Id = fmt.Sprintf("%d", this.ID)
		query.Uuid = this.AccountInfo.Device.DeviceID
		query.Uid = fmt.Sprintf("%d", this.ID)
		query.ServerConfigRetrieval = "1"
	} else {
		query.Id = this.AccountInfo.Device.DeviceID
		query.ServerConfigRetrieval = "1"
	}
	resp := &RespQeSync{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlQeSync,
			Json:    query,
			Header: map[string]string{
				"X-Ig-Connection-Speed": "-1kbps",
			},
			IsPost: true,
			Signed: true,
		}, resp)
	err = resp.CheckError(err)
	if err != nil {
		return err
	}

	//if len(resp.Experiments) == 0 {
	//	return &common.MakeMoneyError{ErrStr: "device info error!"}
	//}
	return err
}

func (this *Instagram) launcherSync() error {
	var query map[string]interface{}
	if this.IsLogin {
		query = map[string]interface{}{
			"id":                      this.ID,
			"_uuid":                   this.AccountInfo.Device.DeviceID,
			"_uid":                    this.ID,
			"server_config_retrieval": "1",
		}
	} else {
		query = map[string]interface{}{
			"id":                      this.AccountInfo.Device.DeviceID,
			"server_config_retrieval": "1",
		}
	}
	resp := &BaseApiResp{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlLauncherSync,
			Query:   query,
			Header: map[string]string{
				"X-Ig-Connection-Speed": "-1kbps",
			},
			IsPost: true,
			Signed: true,
		}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *Instagram) logAttribution() error {
	query := map[string]interface{}{
		"type": "app_first_launch",
		"adid": this.AccountInfo.Device.IDFA,
	}

	resp := &BaseApiResp{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlLogAttribution,
			Query:   query,
			Header: map[string]string{
				"X-Ig-Connection-Speed": "-1kbps",
			},
			IsPost: true,
			Signed: true,
		}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *Instagram) DeviceRegister() error {
	query := map[string]interface{}{
		"_uuid":                    this.AccountInfo.Device.DeviceID,
		"device_id":                this.AccountInfo.Device.DeviceID,
		"device_token":             this.AccountInfo.Device.DeviceToken,
		"family_device_id":         this.AccountInfo.Device.DeviceID,
		"device_app_installations": "{\"threads\":false,\"igtv\":false,\"instagram\":true}",
		"users":                    fmt.Sprintf("%d", this.ID),
		"device_type":              "ios",
	}

	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath:        urlDeviceRegister,
			HeaderSequence: LoginHeaderMap[urlDeviceRegister],
			Query:          query,
			IsPost:         true,
			Signed:         false,
		},
	)
	return err
}

type RespLogin struct {
	BaseApiResp
	LoggedInUser struct {
		AccountBadges                  []interface{} `json:"account_badges"`
		AccountType                    int           `json:"account_type"`
		AllowContactsSync              bool          `json:"allow_contacts_sync"`
		AllowedCommenterType           string        `json:"allowed_commenter_type"`
		BizUserInboxState              int           `json:"biz_user_inbox_state"`
		CanBoostPost                   bool          `json:"can_boost_post"`
		CanSeeOrganicInsights          bool          `json:"can_see_organic_insights"`
		CanSeePrimaryCountryInSettings bool          `json:"can_see_primary_country_in_settings"`
		CountryCode                    int           `json:"country_code"`
		FbidV2                         int64         `json:"fbid_v2"`
		FollowFrictionType             int           `json:"follow_friction_type"`
		FullName                       string        `json:"full_name"`
		HasAnonymousProfilePicture     bool          `json:"has_anonymous_profile_picture"`
		HasPlacedOrders                bool          `json:"has_placed_orders"`
		InteropMessagingUserFbid       int64         `json:"interop_messaging_user_fbid"`
		IsBusiness                     bool          `json:"is_business"`
		IsCallToActionEnabled          interface{}   `json:"is_call_to_action_enabled"`
		IsPrivate                      bool          `json:"is_private"`
		IsUsingUnifiedInboxForDirect   bool          `json:"is_using_unified_inbox_for_direct"`
		IsVerified                     bool          `json:"is_verified"`
		Nametag                        struct {
			Emoji         string `json:"emoji"`
			Gradient      int    `json:"gradient"`
			Mode          int    `json:"mode"`
			SelfieSticker int    `json:"selfie_sticker"`
		} `json:"nametag"`
		NationalNumber                             int64  `json:"national_number"`
		PhoneNumber                                string `json:"phone_number"`
		Pk                                         int64  `json:"pk"`
		ProfessionalConversionSuggestedAccountType int    `json:"professional_conversion_suggested_account_type"`
		ProfilePicUrl                              string `json:"profile_pic_url"`
		ReelAutoArchive                            string `json:"reel_auto_archive"`
		ShowInsightsTerms                          bool   `json:"show_insights_terms"`
		TotalIgtvVideos                            int    `json:"total_igtv_videos"`
		Username                                   string `json:"username"`
		WaAddressable                              bool   `json:"wa_addressable"`
		WaEligibility                              int    `json:"wa_eligibility"`
	} `json:"logged_in_user"`
	SessionFlushNonce interface{} `json:"session_flush_nonce"`
}

func (this *Instagram) Login() error {
	encodePasswd, _ := EncryptPassword(this.Pass, this.GetHeader(IGHeader_EncryptionId), this.GetHeader(IGHeader_EncryptionKey))
	params := map[string]interface{}{
		"phone_id":            this.AccountInfo.Device.DeviceID,
		"reg_login":           "0",
		"device_id":           this.AccountInfo.Device.DeviceID,
		"has_seen_aart_on":    "0",
		"username":            this.User,
		"login_attempt_count": "0",
		"enc_password":        encodePasswd,
	}
	resp := &RespLogin{}
	err := this.HttpRequestJson(&reqOptions{
		ApiPath: urlLogin,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)

	if err == nil {
		if this.GetHeader(IGHeader_Authorization) != "" {
			this.IsLogin = true
		}
		this.ID = resp.LoggedInUser.Pk
	}

	return err
}

type LookResp struct {
	BaseApiResp
	MultipleUsersFound bool   `json:"multiple_users_found"`
	EmailSent          bool   `json:"email_sent"`
	SmsSent            bool   `json:"sms_sent"`
	LookupSource       string `json:"lookup_source"`
	CorrectedInput     string `json:"corrected_input"`
	ObfuscatedPhone    string `json:"obfuscated_phone"`
	User               User   `json:"user"`
	HasValidPhone      bool   `json:"has_valid_phone"`
	CanEmailReset      bool   `json:"can_email_reset"`
	CanSmsReset        bool   `json:"can_sms_reset"`
	CanWaReset         bool   `json:"can_wa_reset"`
	UserId             int64  `json:"user_id"`
	Email              string `json:"email"`
	PhoneNumber        string `json:"phone_number"`
	FbLoginOption      bool   `json:"fb_login_option"`
	IsAutoconfTestUser bool   `json:"is_autoconf_test_user"`
}

func (this *Instagram) UserLookup() (*LookResp, error) {
	resp := &LookResp{}
	err := this.HttpRequestJson(
		&reqOptions{
			ApiPath: urlLookup,
			IsPost:  true,
			Signed:  true,
			Query: map[string]interface{}{
				"q":             this.AccountInfo.Device.DeviceID,
				"skip_recovery": this.AccountInfo.Device.DeviceID,
				"waterfall_id":  this.AccountInfo.Device.WaterID,
			},
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespAddressBookLink struct {
	BaseApiResp
	Items []interface{} `json:"items"`
	Users []interface{} `json:"users"`
}

type AddressBook struct {
	PhoneNumbers   []string `json:"phone_numbers"`
	EmailAddresses []string `json:"email_addresses"`
	LastName       string   `json:"last_name"`
	FirstName      string   `json:"first_name"`
}

func (this *Instagram) AddressBookLink(addr []AddressBook) (*RespAddressBookLink, error) {
	addrJson, err := json.Marshal(addr)
	if err != nil {
		return nil, err
	}

	body := spew.Sprintf("contacts=%s&_uuid=%s&device_id=%s&phone_id=%s", common.InstagramQueryEscape(common.B2s(addrJson)),
		this.AccountInfo.Device.DeviceID, this.AccountInfo.Device.DeviceID, this.AccountInfo.Device.DeviceID)

	resp := &RespAddressBookLink{}
	err = this.HttpRequestJson(
		&reqOptions{
			ApiPath:        urlAddressBookLink + "?include=extra_display_name,thumbnails",
			IsPost:         true,
			Signed:         false,
			HeaderSequence: LoginHeaderMap[urlAddressBookLink],
			Body:           bytes.NewBuffer([]byte(body)),
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}
