package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

type Instagram struct {
	User                string
	Pass                string
	androidID           string
	uuid                string
	token               string
	familyID            string
	adid                string
	wid                 string
	challengeURL        string
	id                  string
	httpHeader          map[string]string
	registerPhoneNumber string
	registerPhoneArea   string
	registerIpCountry   string
	IsLogin             bool

	ReqSuccessCount  int
	ReqErrorCount    int
	ReqApiErrorCount int

	Proxy *common.Proxy

	//Challenge *Challenge
	//Profiles *Profiles
	Account *Account
	//Timeline *Timeline
	//Activity *Activity
	//Inbox *Inbox
	//Feed *Feed
	//Locations *LocationInstance
	Upload *Upload

	c *http.Client
}

func (inst *Instagram) SetCookieJar(jar http.CookieJar) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}
	// First grab the cookies from the existing jar and we'll put it in the new jar.
	cookies := inst.c.Jar.Cookies(url)
	inst.c.Jar = jar
	inst.c.Jar.SetCookies(url, cookies)
	return nil
}

func New(username, password string, _proxy *common.Proxy) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		User:      username,
		Pass:      password,
		androidID: generateDeviceID(),
		uuid:      common.GenerateUUID(), // both uuid must be differents
		familyID:  common.GenerateUUID(),
		wid:       common.GenerateUUID(),
		adid:      common.GenerateUUID(),
		c: &http.Client{
			Jar:       jar,
			Transport: _proxy.GetProxy(),
		},
	}
	inst.Proxy = _proxy
	inst.httpHeader = make(map[string]string)
	common.DebugHttpClient(inst.c)

	inst.init()

	return inst
}

func (inst *Instagram) init() {
	id, err := strconv.ParseInt(inst.id, 10, 64)
	if err != nil {
		log.Warn("account id is null!")
	}
	inst.Account = &Account{inst: inst, ID: id}
	inst.Upload = NewUpload(inst)

	//inst.Challenge = newChallenge(inst)
	//inst.Profiles = newProfiles(inst)
	//inst.Activity = newActivity(inst)
	//inst.Timeline = newTimeline(inst)
	//inst.Inbox = newInbox(inst)
	//inst.Feed = newFeed(inst)
	//inst.Contacts = newContacts(inst)
	//inst.Locations = newLocation(inst)
}

func (inst *Instagram) GetSearch(q string) *Search {
	return newSearch(inst, q)
}

// SetProxy sets proxy for connection.
func (inst *Instagram) SetProxy(_proxy *common.Proxy) {
	inst.Proxy = _proxy
	inst.c.Transport = _proxy.GetProxy()
	common.DebugHttpClient(inst.c)
}

func (inst *Instagram) ReadHeader(key string) string {
	return inst.httpHeader[key]
}

func (inst *Instagram) readMsisdnHeader() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlMsisdnHeader,
			IsPost:  true,
			Query: map[string]interface{}{
				"device_id": inst.uuid,
			},
		},
	)
	return err
}

//注册成功后触发
func (inst *Instagram) contactPrefill() error {
	var query map[string]interface{}

	if inst.IsLogin {
		query = map[string]interface{}{
			"_uid":      inst.id,
			"device_id": inst.uuid,
			"_uuid":     inst.uuid,
			"usage":     "auto_confirmation",
		}
	} else {
		query = map[string]interface{}{
			"phone_id": inst.familyID,
			"usage":    "prefill",
		}
	}

	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlContactPrefill,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (inst *Instagram) launcherSync() error {
	var query map[string]interface{}

	if inst.IsLogin {
		query = map[string]interface{}{
			"id":                      inst.id,
			"_uid":                    inst.id,
			"_uuid":                   inst.uuid,
			"server_config_retrieval": "1",
		}
	} else {
		query = map[string]interface{}{
			"id":                      inst.uuid,
			"server_config_retrieval": "1",
		}
	}

	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlLauncherSync,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query:   query,
		},
	)
	return err
}

func (inst *Instagram) zrToken() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlZrToken,
			IsPost:  false,
			IsApiB:  true,
			Query: map[string]interface{}{
				"device_id":        inst.androidID,
				"token_hash":       "",
				"custom_device_id": inst.uuid,
				"fetch_reason":     "token_expired",
			},
			HeaderKey: []string{IGHeader_Authorization},
		},
	)
	return err
}

//早于注册登录?
func (inst *Instagram) sendAdID() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlLogAttribution,
			IsPost:  true,
			IsApiB:  true,
			Signed:  true,
			Query: map[string]interface{}{
				"adid": inst.adid,
			},
		},
	)
	return err
}

//is_register bool
func (inst *Instagram) Prepare() error {
	err := inst.readMsisdnHeader()
	if err != nil {
		return err
	}

	err = inst.syncFeatures()
	if err != nil {
		return err
	}

	err = inst.zrToken()
	if err != nil {
		return err
	}

	err = inst.sendAdID()
	if err != nil {
		return err
	}

	err = inst.contactPrefill()
	//if err != nil {
	//	return err
	//}
	return nil
}

// Login performs instagram login.
//
// Password will be deleted after login

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
		FbidV2                         int64         `json:"fbid_v_2"`
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
		NationalNumber                             string `json:"national_number"`
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

func (inst *Instagram) Login() error {
	encodePasswd, _ := encryptPassword(inst.Pass, inst.ReadHeader(IGHeader_EncryptionId), inst.ReadHeader(IGHeader_EncryptionKey))
	params := map[string]interface{}{
		"jazoest":             genJazoest(inst.familyID),
		"country_codes":       "[{\"country_code\":\"" + strings.ReplaceAll(inst.registerPhoneArea, "+", "") + "\",\"source\":[\"default\"]}]",
		"phone_id":            inst.familyID,
		"enc_password":        encodePasswd,
		"username":            inst.User,
		"adid":                inst.adid,
		"guid":                inst.uuid,
		"device_id":           inst.androidID,
		"google_tokens":       "[]",
		"login_attempt_count": "0",
	}
	resp := &RespLogin{}
	err := inst.HttpRequestJson(&reqOptions{
		Login:   false,
		ApiPath: urlLogin,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

// Logout closes current session
func (inst *Instagram) Logout() error {
	_, err := inst.sendSimpleRequest(urlLogout)
	inst.c.Jar = nil
	inst.c = nil
	return err
}

func (inst *Instagram) syncFeatures() error {
	var params map[string]interface{}
	if inst.IsLogin {
		params = map[string]interface{}{
			"id":          inst.id,
			"_uid":        inst.id,
			"_uuid":       inst.uuid,
			"experiments": goInstaExperiments,
		}
	} else {
		params = map[string]interface{}{
			"id":          inst.uuid,
			"experiments": goInstaExperiments,
		}
	}
	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlQeSync,
			Query:   params,
			IsPost:  true,
			Login:   true,
			Signed:  true,
		},
	)
	return err
}

func (inst *Instagram) megaphoneLog() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			ApiPath: urlMegaphoneLog,
			Query: map[string]interface{}{
				"id":        strconv.FormatInt(inst.Account.ID, 10),
				"type":      "feed_aysf",
				"action":    "seen",
				"reason":    "",
				"device_id": inst.androidID,
				"uuid":      common.GenerateMD5Hash(string(time.Now().Unix())),
			},
			IsPost: true,
			Login:  true,
		},
	)
	return err
}

//func (inst *Instagram) expose() error {
//	data, err := inst.prepareData(
//		map[string]interface{}{
//			"id":         inst.Account.ID,
//			"experiment": "ig_android_profile_contextual_feed",
//		},
//	)
//	if err != nil {
//		return err
//	}
//
//	_, err = inst.sendRequest(
//		&reqOptions{
//			ApiPath: urlExpose,
//			Query:    generateSignature(data),
//			IsPost:   true,
//		},
//	)
//
//	return err
//}

// GetMedia returns media specified by id.
//
// The argument can be int64 or string
//
// See example: examples/media/like.go
//func (inst *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
//	media := &FeedMedia{
//		inst:   inst,
//		NextID: o,
//	}
//	return media, media.Sync()
//}
