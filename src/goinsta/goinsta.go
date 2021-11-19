package goinsta

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"makemoney/proxy"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Instagram represent the main API handler
//
// Profiles: Represents instragram's user profile.
// Account:  Represents instagram's personal account.
// Search:   Represents instagram's search.
// Timeline: Represents instagram's timeline.
// Activity: Represents instagram's user activity.
// Inbox:    Represents instagram's messages.
// Location: Represents instagram's locations.
//
// See Scheme section in README.md for more information.
//
// We recommend to use Export and Import functions after first Login.
//
// Also you can use SetProxy and UnsetProxy to set and unset proxy.
// Golang also provides the option to set a proxy using HTTP_PROXY env var.
type Instagram struct {
	User string
	Pass string
	// device id: android-1923fjnma8123
	androidID string
	// uuid: 8493-1233-4312312-5123
	uuid string
	// rankToken
	rankToken string
	// token
	token string

	familyID string
	// ads id
	adid string
	//waterfall id
	wid string
	// challenge URL
	challengeURL string

	id         string
	httpHeader map[string]string

	IsLogin bool

	ReqSuccessCount  int
	ReqErrorCount    int
	ReqApiErrorCount int

	proxy *proxy.Proxy

	// Instagram objects
	// Challenge controls security side of account (Like sms verify / It was me)
	Challenge *Challenge
	// Profiles is the user interaction
	Profiles *Profiles
	// Account stores all personal data of the user and his/her options.
	Account *Account
	// Search performs searching of multiple things (users, locations...)
	Search *Search
	// Timeline allows to receive timeline media.
	Timeline *Timeline
	// Activity are instagram notifications.
	Activity *Activity
	// Inbox are instagram message/chat system.
	Inbox *Inbox
	// Feed for search over feeds
	Feed *Feed
	// User contacts from mobile address book
	Contacts *Contacts
	// Location instance
	Locations *LocationInstance

	c *http.Client
}

// SetHTTPTransport sets http transport. This further allows users to tweak the underlying
// low level transport for adding additional fucntionalities.
func (inst *Instagram) SetHTTPTransport(transport http.RoundTripper) {
	inst.c.Transport = transport
}

// SetCookieJar sets the Cookie Jar. This further allows to use a custom implementation
// of a cookie jar which may be backed by a different data store such as redis.
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

// New creates Instagram structure
func New(username, password string, _proxy *proxy.Proxy) *Instagram {
	// this call never returns error
	jar, _ := cookiejar.New(nil)
	inst := &Instagram{
		User:      username,
		Pass:      password,
		androidID: generateDeviceID(),
		uuid:      generateUUID(), // both uuid must be differents
		familyID:  generateUUID(),
		wid:       generateUUID(),
		c: &http.Client{
			Jar:       jar,
			Transport: _proxy.GetProxy(),
		},
	}
	inst.proxy = _proxy
	//tools.DebugHttpClient(inst.c)

	inst.init()

	return inst
}

func (inst *Instagram) init() {
	inst.Challenge = newChallenge(inst)
	inst.Profiles = newProfiles(inst)
	inst.Activity = newActivity(inst)
	inst.Timeline = newTimeline(inst)
	inst.Search = newSearch(inst)
	inst.Inbox = newInbox(inst)
	inst.Feed = newFeed(inst)
	inst.Contacts = newContacts(inst)
	inst.Locations = newLocation(inst)
}

// SetProxy sets proxy for connection.
func (inst *Instagram) SetProxy(_proxy *proxy.Proxy) {
	inst.proxy = _proxy
	inst.c.Transport = _proxy.GetProxy()
}

// Save exports config to ~/.goinsta
func (inst *Instagram) Save() error {
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("home") // for plan9
	}
	return inst.Export(filepath.Join(home, ".goinsta"))
}

// Export exports *Instagram object options
func (inst *Instagram) Export(path string) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:        inst.Account.ID,
		User:      inst.User,
		AndroidID: inst.androidID,
		UUID:      inst.uuid,
		RankToken: inst.rankToken,
		Token:     inst.token,
		FamilyID:  inst.familyID,
		Cookies:   inst.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, bytes, 0644)
}

// Export exports selected *Instagram object options to an io.Writer
func Export(inst *Instagram, writer io.Writer) error {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return err
	}

	config := ConfigFile{
		ID:        inst.Account.ID,
		User:      inst.User,
		AndroidID: inst.androidID,
		UUID:      inst.uuid,
		RankToken: inst.rankToken,
		Token:     inst.token,
		FamilyID:  inst.familyID,
		Cookies:   inst.c.Jar.Cookies(url),
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

// ImportReader imports instagram configuration from io.Reader
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportReader(r io.Reader) (*Instagram, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := ConfigFile{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return ImportConfig(config)
}

// ImportConfig imports instagram configuration from a configuration object.
//
// This function does not set proxy automatically. Use SetProxy after this call.
func ImportConfig(config ConfigFile) (*Instagram, error) {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return nil, err
	}

	inst := &Instagram{
		User:      config.User,
		androidID: config.AndroidID,
		uuid:      config.UUID,
		rankToken: config.RankToken,
		token:     config.Token,
		familyID:  config.FamilyID,
		c: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		},
	}
	inst.c.Jar, err = cookiejar.New(nil)
	if err != nil {
		return inst, err
	}
	inst.c.Jar.SetCookies(url, config.Cookies)

	inst.init()
	inst.Account = &Account{inst: inst, ID: config.ID}
	inst.Account.Sync()

	return inst, nil
}

// Import imports instagram configuration
//
// This function does not set proxy automatically. Use SetProxy after this call.
func Import(path string) (*Instagram, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ImportReader(f)
}

func (inst *Instagram) ReadHeader(key string) string {
	return inst.httpHeader[key]
}

func (inst *Instagram) readMsisdnHeader() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlMsisdnHeader,
			IsPost:   true,
			Query: map[string]string{
				"device_id": inst.uuid,
			},
		},
	)
	return err
}

//注册成功后触发
func (inst *Instagram) contactPrefill() error {
	var query map[string]string

	if inst.IsLogin {
		query = map[string]string{
			"_uid":      inst.id,
			"device_id": inst.uuid,
			"_uuid":     inst.uuid,
			"usage":     "auto_confirmation",
		}
	} else {
		query = map[string]string{
			"phone_id": inst.familyID,
			"usage":    "prefill",
		}
	}

	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlContactPrefill,
			IsPost:   true,
			IsApiB:   true,
			Query:    query,
		},
	)
	return err
}

func (inst *Instagram) launcherSync() error {
	var query map[string]string

	if inst.IsLogin {
		query = map[string]string{
			"id":                      inst.id,
			"_uid":                    inst.id,
			"_uuid":                   inst.uuid,
			"server_config_retrieval": "1",
		}
	} else {
		query = map[string]string{
			"id":                      inst.uuid,
			"server_config_retrieval": "1",
		}
	}

	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlLauncherSync,
			IsPost:   true,
			IsApiB:   true,
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
			Query: map[string]string{
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
			Query: map[string]string{
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
//			Query:    generateSignature(b2s(result)),
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
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlQeSync,
			Query: map[string]string{
				"id":          inst.id,
				"_uid":        inst.id,
				"_uuid":       inst.uuid,
				"experiments": goInstaExperiments,
			},
			IsPost: true,
			Login:  true,
		},
	)
	return err
}

func (inst *Instagram) megaphoneLog() error {
	_, err := inst.HttpRequest(
		&reqOptions{
			Endpoint: urlMegaphoneLog,
			Query: map[string]string{
				"id":        strconv.FormatInt(inst.Account.ID, 10),
				"type":      "feed_aysf",
				"action":    "seen",
				"reason":    "",
				"device_id": inst.androidID,
				"uuid":      generateMD5Hash(string(time.Now().Unix())),
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
func (inst *Instagram) GetMedia(o interface{}) (*FeedMedia, error) {
	media := &FeedMedia{
		inst:   inst,
		NextID: o,
	}
	return media, media.Sync()
}
