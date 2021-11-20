package goinsta

import (
	"makemoney/goinsta/dbhelper"
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
		RankToken:           inst.rankToken,
		Token:               inst.token,
		FamilyID:            inst.familyID,
		Cookies:             inst.c.Jar.Cookies(url),
		CookiesB:            inst.c.Jar.Cookies(urlb),
		Adid:                inst.adid,
		Wid:                 inst.wid,
		HttpHeader:          inst.httpHeader,
		ProxyID:             inst.proxy.Id,
		IsLogin:             inst.IsLogin,
		RegisterPhoneNumber: inst.registerPhoneNumber,
		RegisterPhoneArea:   inst.registerPhoneArea,
		RegisterIpCountry:   inst.registerIpCountry,
	}
	return dbhelper.SaveNewAccount(Cookies)
}

//
//func ImportReader(r io.Reader) (*Instagram, error) {
//	bytes, err := ioutil.ReadAll(r)
//	if err != nil {
//		return nil, err
//	}
//
//	config := AccountCookies{}
//	err = json.Unmarshal(bytes, &config)
//	if err != nil {
//		return nil, err
//	}
//	return ImportConfig(config)
//}
//
//// ImportConfig imports instagram configuration from a configuration object.
////
//// This function does not set proxy automatically. Use SetProxy after this call.
//func ImportConfig(config AccountCookies) (*Instagram, error) {
//	url, err := neturl.Parse(goInstaAPIUrl)
//	if err != nil {
//		return nil, err
//	}
//
//	inst := &Instagram{
//		User:      config.User,
//		androidID: config.AndroidID,
//		uuid:      config.UUID,
//		rankToken: config.RankToken,
//		token:     config.Token,
//		familyID:  config.FamilyID,
//		c: &http.Client{
//			Transport: &http.Transport{
//				Proxy: http.ProxyFromEnvironment,
//			},
//		},
//	}
//	inst.c.Jar, err = cookiejar.New(nil)
//	if err != nil {
//		return inst, err
//	}
//	inst.c.Jar.SetCookies(url, config.Cookies)
//
//	inst.init()
//	inst.Account = &Account{inst: inst, ID: config.ID}
//	inst.Account.Sync()
//
//	return inst, nil
//}
//
//// Import imports instagram configuration
////
//// This function does not set proxy automatically. Use SetProxy after this call.
//func Import(path string) (*Instagram, error) {
//	f, err := os.Open(path)
//	if err != nil {
//		return nil, err
//	}
//	defer f.Close()
//	return ImportReader(f)
//}
