package goinsta

import (
	"makemoney/common"
	"makemoney/common/proxy"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"strings"
)

var (
	InsAccountError_ChallengeRequired = "challenge_required"
	InsAccountError_LoginRequired     = "login_required"
	InsAccountError_Feedback          = "feedback_required"
)

var ProxyCallBack func(id string) (*proxy.Proxy, error)

type Instagram struct {
	User                string
	Pass                string
	deviceID            string
	token               string
	familyID            string
	wid                 string
	challengeURL        string
	ID                  int64
	httpHeader          map[string]string
	IsLogin             bool
	UserAgent           string
	Status              string
	sessionID           string
	RegisterPhoneNumber string
	RegisterPhoneArea   string
	RegisterIpCountry   string
	RegisterTime        int64
	ReqSuccessCount     int
	ReqErrorCount       int
	ReqApiErrorCount    int
	ReqContError        int
	LastSendMsgTime     int
	MatePoint           interface{}
	Proxy               *proxy.Proxy
	c                   *http.Client
	graph               *Graph
	account             *Account
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

func New(username, password string, _proxy *proxy.Proxy) *Instagram {
	var tr *http.Transport
	if _proxy != nil {
		tr = _proxy.GetProxy()
	}

	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		User:      username,
		Pass:      password,
		deviceID:  strings.ToUpper(common.GenUUID()),
		wid:       common.GenUUID(),
		Proxy:     _proxy,
		UserAgent: GenUserAgent(),
		sessionID: strings.ToUpper(common.GenUUID()),
		c: &http.Client{
			Jar:       jar,
			Transport: tr,
		},
	}

	inst.familyID = inst.deviceID
	inst.graph = &Graph{inst: inst}
	inst.httpHeader = make(map[string]string)

	common.DebugHttpClient(inst.c)
	return inst
}

func (this *Instagram) GetSearch(q string) *Search {
	return newSearch(this, q)
}

func (this *Instagram) GetUpload() *Upload {
	return newUpload(this)
}

func (this *Instagram) GetAccount() *Account {
	if this.account == nil {
		this.account = &Account{ID: this.ID, inst: this}
	}
	return this.account
}

func (this *Instagram) GetUser(id string) *User {
	pk, _ := strconv.ParseInt(id, 10, 64)
	return &User{ID: pk, inst: this}
}

func (this *Instagram) GetFollowers(id string) *Followers {
	pk, _ := strconv.ParseInt(id, 10, 64)
	return &Followers{User: pk, inst: this, HasMore: true}
}

func (this *Instagram) GetMessage() *Message {
	return &Message{inst: this}
}

// SetProxy sets proxy for connection.
func (this *Instagram) SetProxy(_proxy *proxy.Proxy) {
	this.Proxy = _proxy
	this.c.Transport = _proxy.GetProxy()
	common.DebugHttpClient(this.c)
}

func (this *Instagram) NeedReplace() bool {
	if this.Status == InsAccountError_ChallengeRequired {
		return true
	}

	//if this.ReqContError >= 3 {
	//	return true
	//}
	return false
}

func (this *Instagram) CleanCookiesAndHeader() {
	this.httpHeader = make(map[string]string)
	this.c.Jar, _ = cookiejar.New(nil)
	this.IsLogin = false
}

func (this *Instagram) ReadHeader(key string) string {
	return this.httpHeader[key]
}

func (this *Instagram) PrepareNewClient() {
	_ = this.contactPrefill()
	_ = this.qeSync()
	_ = this.launcherSync()
	_ = this.getNamePrefill()
}

func (this *Instagram) qeSync() error {
	var params = map[string]interface{}{
		"id":                      this.deviceID,
		"server_config_retrieval": "1",
	}
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlQeSync,
			Query:   params,
			IsPost:  true,
			Signed:  true,
		},
	)
	return err
}

func (this *Instagram) launcherSync() error {
	var query = map[string]interface{}{
		"id":                      this.deviceID,
		"server_config_retrieval": "1",
	}
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlLauncherSync,
			IsPost:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (this *Instagram) getNamePrefill() error {
	var query = map[string]interface{}{
		"phone_id":  this.deviceID,
		"device_id": this.deviceID,
	}
	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlGetNamePrefill,
			IsPost:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (this *Instagram) contactPrefill() error {
	var query = map[string]interface{}{
		"phone_id": this.deviceID,
	}

	_, err := this.HttpRequest(
		&reqOptions{
			ApiPath: urlContactPrefill,
			IsPost:  true,
			IsApiB:  false,
			Signed:  true,
			Query:   query,
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
	encodePasswd, _ := encryptPassword(this.Pass, this.ReadHeader(IGHeader_EncryptionId), this.ReadHeader(IGHeader_EncryptionKey))
	params := map[string]interface{}{
		"phone_id":            this.deviceID,
		"reg_login":           "0",
		"device_id":           this.deviceID,
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
		if this.ReadHeader(IGHeader_Authorization) != "" {
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
				"q":             this.deviceID,
				"skip_recovery": this.deviceID,
				"waterfall_id":  this.wid,
			},
		}, resp)

	err = resp.CheckError(err)
	return resp, err
}
