package goinsta

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
)

type InstDeviceInfoV1 struct {
	Version        string  `json:"version" bson:"version"`
	VersionCode    string  `json:"version_code" bson:"version_code"`
	BloksVersionID string  `json:"bloks_version_id" bson:"bloks_version_id"`
	UserAgent      string  `json:"user_agent" bson:"user_agent"`
	IDFA           string  `json:"idfa" bson:"idfa"`
	AppLocale      string  `json:"app_locale" bson:"app_locale"`
	TimezoneOffset string  `json:"timezone_offset" bson:"timezone_offset"`
	StartupCountry string  `json:"startup_country" bson:"startup_country"`
	AcceptLanguage string  `json:"accept_language" bson:"accept_language"`
	NetWorkType    string  `json:"net_work_type" bson:"net_work_type"`
	DeviceID       string  `json:"device_id" bson:"device_id"`
	FamilyID       string  `json:"family_id" bson:"family_id"`
	WaterID        string  `json:"water_id" bson:"water_id"`
	DeviceToken    string  `json:"device_token" bson:"device_token"`
	SystemVersion  string  `json:"system_version" bson:"system_version"`
	LensModel      string  `json:"lens_model" bson:"lens_model"`
	FocalLength    float64 `json:"focal_length" bson:"focal_length"`
	Aperture       float64 `json:"aperture" bson:"aperture"`
	Longitude      float64 `json:"longitude" bson:"longitude"`
	Latitude       float64 `json:"latitude" bson:"latitude"`
}
type AccountCookiesV1 struct {
	ID                  int64                    `json:"id" bson:"id"`
	Username            string                   `json:"username" bson:"username"`
	Passwd              string                   `json:"passwd" bson:"passwd"`
	HttpHeader          map[string]string        `json:"http_header" bson:"http_header"`
	ProxyID             string                   `json:"proxy_id" bson:"proxy_id"`
	IsLogin             bool                     `json:"is_login" bson:"is_login"`
	Token               string                   `json:"token" bson:"token"`
	Cookies             []*http.Cookie           `json:"cookies" bson:"cookies"`
	CookiesB            []*http.Cookie           `json:"cookies_b" bson:"cookies_b"`
	Device              *InstDeviceInfoV1        `json:"device" bson:"device"`
	RegisterEmail       string                   `json:"register_email" bson:"register_email"`
	RegisterPhoneNumber string                   `json:"register_phone_number" bson:"register_phone_number"`
	RegisterPhoneArea   string                   `json:"register_phone_area" bson:"register_phone_area"`
	RegisterIpCountry   string                   `json:"register_ip_country" bson:"register_ip_country"`
	RegisterTime        int64                    `json:"register_time" bson:"register_time"`
	Status              string                   `json:"status" bson:"status"`
	LastSendMsgTime     int                      `json:"last_send_msg_time" bson:"last_send_msg_time"`
	Tags                string                   `json:"tags" bson:"tags"`
	SpeedControl        map[string]*SpeedControl `json:"speed_control" bson:"speed_control"`
}

func ReplaceInstToDB(inst *Instagram) error {
	url, _ := neturl.Parse(InstagramHost)
	urlb, _ := neturl.Parse(InstagramHost_B)

	Cookies := AccountCookies{
		ID:           inst.ID,
		Username:     inst.User,
		Passwd:       inst.Pass,
		Token:        inst.token,
		AccountInfo:  inst.AccountInfo,
		Cookies:      inst.c.Jar.Cookies(url),
		CookiesB:     inst.c.Jar.Cookies(urlb),
		HttpHeader:   inst.httpHeader,
		ProxyID:      inst.Proxy.ID,
		IsLogin:      inst.IsLogin,
		Status:       inst.Status,
		SpeedControl: inst.SpeedControl,
		Tags:         inst.Tags,
	}
	return ReplaceAccount(Cookies)
}

func ReStruct() {
	cursor, err := MogoHelper.Account.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
	}
	var ret []AccountCookiesV1
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
	}

	for idx, item := range ret {
		log.Info("%d", idx)
		config1, err := ConvConfig1(&item)
		if err != nil {
			log.Error("%v", err)
		}
		err = ReplaceInstToDB(config1)
		if err != nil {
			log.Error("%v", err)
		}
	}
}

func ConvConfig1(config *AccountCookiesV1) (*Instagram, error) {
	url, err := neturl.Parse(InstagramHost)
	if err != nil {
		return nil, err
	}
	urlb, err := neturl.Parse(InstagramHost_B)
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	jar.SetCookies(url, config.Cookies)
	jar.SetCookies(urlb, config.CookiesB)

	inst := &Instagram{
		ID:           config.ID,
		User:         config.Username,
		Pass:         config.Passwd,
		token:        config.Token,
		httpHeader:   config.HttpHeader,
		IsLogin:      config.IsLogin,
		Status:       config.Status,
		SpeedControl: config.SpeedControl,

		AccountInfo: &InstAccountInfo{
			Device: InstDeviceInfo{
				Version:        config.Device.Version,
				VersionCode:    config.Device.VersionCode,
				BloksVersionID: config.Device.BloksVersionID,
				UserAgent:      config.Device.UserAgent,
				IDFA:           config.Device.IDFA,
				DeviceID:       config.Device.DeviceID,
				FamilyID:       config.Device.FamilyID,
				WaterID:        config.Device.WaterID,
				DeviceToken:    config.Device.DeviceToken,
				SystemVersion:  config.Device.SystemVersion,
				LensModel:      config.Device.LensModel,
				FocalLength:    config.Device.FocalLength,
				Aperture:       config.Device.Aperture,
				NetWorkType:    config.Device.NetWorkType,
			},
			Location: InstLocationInfo{
				Lon:            0,
				Lat:            0,
				AppLocale:      config.Device.AppLocale,
				StartupCountry: config.Device.StartupCountry,
				AcceptLanguage: config.Device.AcceptLanguage,
				Timezone:       config.Device.TimezoneOffset,
			},
			Register: InstRegisterInfo{
				RegisterEmail:       config.RegisterEmail,
				RegisterPhoneNumber: config.RegisterPhoneNumber,
				RegisterPhoneArea:   config.RegisterPhoneArea,
				RegisterIpCountry:   config.RegisterIpCountry,
				RegisterTime:        config.RegisterTime,
			},
		},
		sessionID: strings.ToUpper(common.GenUUID()),
		c: &http.Client{
			Jar: jar,
		},
		Tags: config.Tags,
	}

	if inst.SpeedControl == nil {
		inst.SpeedControl = make(map[string]*SpeedControl)
	}
	for key, value := range inst.SpeedControl {
		ReSetRate(value, key)
	}

	inst.Operate.Graph = &Graph{inst: inst}

	return inst, nil
}
