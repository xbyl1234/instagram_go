package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta/dbhelper"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
)

func SaveInstToDB(inst *Instagram) error {
	url, _ := neturl.Parse(goInstaAPIUrl)
	urlb, _ := neturl.Parse(goInstaAPIUrl_B)

	Cookies := dbhelper.AccountCookies{
		ID:                  inst.id,
		Username:            inst.User,
		Passwd:              inst.Pass,
		AndroidID:           inst.androidID,
		UUID:                inst.uuid,
		Token:               inst.token,
		FamilyID:            inst.familyID,
		Cookies:             inst.c.Jar.Cookies(url),
		CookiesB:            inst.c.Jar.Cookies(urlb),
		Adid:                inst.adid,
		Wid:                 inst.wid,
		HttpHeader:          inst.httpHeader,
		ProxyID:             inst.Proxy.ID,
		IsLogin:             inst.IsLogin,
		RegisterPhoneNumber: inst.registerPhoneNumber,
		RegisterPhoneArea:   inst.registerPhoneArea,
		RegisterIpCountry:   inst.registerIpCountry,
	}
	return dbhelper.SaveNewAccount(Cookies)
}

func LoadAllAccount() []*Instagram {
	config, err := dbhelper.LoadDBAllAccount()
	if err != nil {
		return nil
	}
	var ret []*Instagram
	for item := range config {
		inst, err := ConvConfig(&config[item])
		if err != nil {
			log.Warn("conv config to inst error:%v", err)
			continue
		}
		ret = append(ret, inst)
	}
	return ret
}

func ConvConfig(config *dbhelper.AccountCookies) (*Instagram, error) {
	url, err := neturl.Parse(goInstaAPIUrl)
	if err != nil {
		return nil, err
	}
	urlb, err := neturl.Parse(goInstaAPIUrl_B)
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
		id:                  config.ID,
		User:                config.Username,
		Pass:                config.Passwd,
		androidID:           config.AndroidID,
		uuid:                config.UUID,
		token:               config.Token,
		familyID:            config.FamilyID,
		adid:                config.Adid,
		wid:                 config.Wid,
		httpHeader:          config.HttpHeader,
		IsLogin:             config.IsLogin,
		registerPhoneNumber: config.RegisterPhoneNumber,
		registerPhoneArea:   config.RegisterPhoneArea,
		registerIpCountry:   config.RegisterIpCountry,
		c: &http.Client{
			Jar: jar,
		},
	}

	inst.Proxy = &common.Proxy{ID: config.ProxyID}

	inst.init()

	return inst, nil
}
