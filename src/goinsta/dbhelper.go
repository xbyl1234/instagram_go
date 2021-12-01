package goinsta

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"time"
)

type PhoneStorage struct {
	Area          string        `bson:"area"`
	Phone         string        `bson:"phone"`
	SendCount     string        `bson:"send_count"`
	RegisterCount int           `bson:"register_count"`
	Provider      string        `bson:"provider"`
	LastUseTime   time.Duration `bson:"last_use_time"`
}

func UpdatePhoneSendOnce(provider string, area string, number string) error {
	_, err := common.MogoHelper.Phone.UpdateOne(context.TODO(),
		bson.D{
			{"area", area},
			{"phone", number},
		}, bson.D{{"$set", bson.D{{"area", area},
			{"phone", number},
			{"provider", provider},
			{"last_use_time", time.Now()},
		}}, {"$inc", bson.M{"send_count": 1}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func UpdatePhoneRegisterOnce(area string, number string) error {
	_, err := common.MogoHelper.Phone.UpdateOne(context.TODO(),
		bson.D{
			{"area", area},
			{"phone", number},
		}, bson.D{{"$inc", bson.M{"register_count": 1}}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

type AccountCookies struct {
	ID                  int64             `json:"id"`
	Username            string            `json:"username"`
	Passwd              string            `json:"passwd"`
	Adid                string            `json:"adid"`
	Wid                 string            `json:"wid"`
	HttpHeader          map[string]string `json:"http_header"`
	ProxyID             string            `json:"proxy_id"`
	IsLogin             bool              `json:"is_login"`
	AndroidID           string            `json:"android_id"`
	UUID                string            `json:"uuid"`
	Token               string            `json:"token"`
	FamilyID            string            `json:"family_id"`
	Cookies             []*http.Cookie    `json:"cookies"`
	CookiesB            []*http.Cookie    `json:"cookies_b"`
	RegisterPhoneNumber string            `json:"register_phone_number"`
	RegisterPhoneArea   string            `json:"register_phone_area"`
	RegisterIpCountry   string            `json:"register_ip_country"`
}

func SaveNewAccount(account AccountCookies) error {
	_, err := common.MogoHelper.Account.UpdateOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account},
		options.Update().SetUpsert(true))
	return err
}

func LoadDBAllAccount() ([]AccountCookies, error) {
	cursor, err := common.MogoHelper.Account.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var ret []AccountCookies
	err = cursor.All(context.TODO(), &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func SaveInstToDB(inst *Instagram) error {
	url, _ := neturl.Parse(goInstaAPIUrl)
	urlb, _ := neturl.Parse(goInstaAPIUrl_B)

	Cookies := AccountCookies{
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
	return SaveNewAccount(Cookies)
}

func LoadAllAccount() []*Instagram {
	config, err := LoadDBAllAccount()
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

func ConvConfig(config *AccountCookies) (*Instagram, error) {
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

	return inst, nil
}
