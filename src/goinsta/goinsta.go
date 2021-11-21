package goinsta

import (
	"makemoney/proxy"
	"makemoney/tools"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strconv"
	"time"
)

type Instagram struct {
	User                string
	Pass                string
	androidID           string
	uuid                string
	rankToken           string
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

	proxy *proxy.Proxy

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

func New(username, password string, _proxy *proxy.Proxy) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		User:      username,
		Pass:      password,
		androidID: generateDeviceID(),
		uuid:      tools.GenerateUUID(), // both uuid must be differents
		familyID:  tools.GenerateUUID(),
		wid:       tools.GenerateUUID(),
		adid:      tools.GenerateUUID(),
		c: &http.Client{
			Jar:       jar,
			Transport: _proxy.GetProxy(),
		},
	}
	inst.proxy = _proxy
	inst.httpHeader = make(map[string]string)
	tools.DebugHttpClient(inst.c)

	inst.init()

	return inst
}

func (inst *Instagram) init() {
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
func (inst *Instagram) SetProxy(_proxy *proxy.Proxy) {
	inst.proxy = _proxy
	inst.c.Transport = _proxy.GetProxy()
}

func (inst *Instagram) ReadHeader(key string) string {
	return inst.httpHeader[key]
}

func (inst *Instagram) readMsisdnHeader() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlMsisdnHeader,
			IsPost:   true,
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
			Endpoint: urlContactPrefill,
			IsPost:   true,
			IsApiB:   true,
			Signed:   true,
			Query:    query,
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
			Endpoint: urlLauncherSync,
			IsPost:   true,
			IsApiB:   true,
			Signed:   true,
			Query:    query,
		},
	)
	return err
}

func (inst *Instagram) zrToken() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlZrToken,
			IsPost:   false,
			IsApiB:   true,
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
			Endpoint: urlLogAttribution,
			IsPost:   true,
			IsApiB:   true,
			Signed:   true,
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
//func (inst *Instagram) Login() error {
//	err := inst.Prepare()
//	if err != nil {
//		return err
//	}
//
//	result, err := json.Marshal(
//		map[string]interface{}{
//			"guid":                inst.uuid,
//			"login_attempt_count": 0,
//			"_csrftoken":          inst.token,
//			"device_id":           inst.androidID,
//			"adid":                inst.adid,
//			"phone_id":            inst.pid,
//			"username":            inst.User,
//			"password":            inst.Pass,
//			"google_tokens":       "[]",
//		},
//	)
//	if err != nil {
//		return err
//	}
//	body, err := inst.sendRequest(
//		&reqOptions{
//			Endpoint: urlLogin,
//			Query:    generateSignature(tools.B2s(result)),
//			IsPost:   true,
//			Login:    true,
//		},
//	)
//	if err != nil {
//		return err
//	}
//	inst.Pass = ""
//
//	// getting account data
//	res := accountResp{}
//	err = json.Unmarshal(body, &res)
//	if err != nil {
//		return err
//	}
//
//	inst.Account = &res.Account
//	inst.Account.inst = inst
//	inst.rankToken = strconv.FormatInt(inst.Account.ID, 10) + "_" + inst.uuid
//	inst.zrToken()
//
//	return err
//}

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
			Endpoint: urlQeSync,
			Query:    params,
			IsPost:   true,
			Login:    true,
			Signed:   true,
		},
	)
	return err
}

func (inst *Instagram) megaphoneLog() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlMegaphoneLog,
			Query: map[string]interface{}{
				"id":        strconv.FormatInt(inst.Account.ID, 10),
				"type":      "feed_aysf",
				"action":    "seen",
				"reason":    "",
				"device_id": inst.androidID,
				"uuid":      tools.GenerateMD5Hash(string(time.Now().Unix())),
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
//			Endpoint: urlExpose,
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
