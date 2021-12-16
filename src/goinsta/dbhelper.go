package goinsta

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"time"
)

type MogoDBHelper struct {
	Client         *mongo.Client
	Phone          *mongo.Collection
	Account        *mongo.Collection
	UploadIDRecord *mongo.Collection
}

var MogoHelper *MogoDBHelper = nil

func InitMogoDB(mogoUri string) {
	//"mongodb://xbyl:XBYLxbyl1234@62.216.92.183:27017"
	clientOptions := options.Client().ApplyURI(mogoUri)
	//clientOptions := options.Client().ApplyURI("mongodb://xbyl:xbyl741852JHK@192.168.187.1:27017")

	var err error
	MogoHelper = &MogoDBHelper{}

	MogoHelper.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("mongo %v", err)
	}

	err = MogoHelper.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("mongo %v", err)
	}

	MogoHelper.Phone = MogoHelper.Client.Database("inst").Collection("phone")
	MogoHelper.Account = MogoHelper.Client.Database("inst").Collection("account")
	MogoHelper.UploadIDRecord = MogoHelper.Client.Database("inst").Collection("upload_id")
}

func GetDB(name string) *mongo.Database {
	return MogoHelper.Client.Database(name)
}

type PhoneStorage struct {
	Area          string        `bson:"area"`
	Phone         string        `bson:"phone"`
	SendCount     string        `bson:"send_count"`
	RegisterCount int           `bson:"register_count"`
	Provider      string        `bson:"provider"`
	LastUseTime   time.Duration `bson:"last_use_time"`
}

func UpdatePhoneSendOnce(provider string, area string, number string) error {
	_, err := MogoHelper.Phone.UpdateOne(context.TODO(),
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
	_, err := MogoHelper.Phone.UpdateOne(context.TODO(),
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
	Status              string            `json:"status"`
	UserAgent           string            `json:"user_agent"`
	LastSendMsgTime     int               `json:"last_send_msg_time"`
}

func SaveNewAccount(account AccountCookies) error {
	_, err := MogoHelper.Account.UpdateOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account},
		options.Update().SetUpsert(true))
	return err
}

func LoadDBAllAccount() ([]AccountCookies, error) {
	cursor, err := MogoHelper.Account.Find(context.TODO(), bson.M{}, nil)
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

func CleanStatus() error {
	_, err := MogoHelper.Account.UpdateMany(
		context.TODO(),
		bson.M{},
		bson.M{"$set": bson.M{"status": ""}},
		options.Update().SetUpsert(true))
	return err
}

func SaveInstToDB(inst *Instagram) error {
	url, _ := neturl.Parse(goInstaHost)
	urlb, _ := neturl.Parse(goInstaHost_B)

	Cookies := AccountCookies{
		ID:                  inst.ID,
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
		RegisterPhoneNumber: inst.RegisterPhoneNumber,
		RegisterPhoneArea:   inst.RegisterPhoneArea,
		RegisterIpCountry:   inst.RegisterIpCountry,
		Status:              inst.Status,
		UserAgent:           inst.UserAgent,
		LastSendMsgTime:     inst.LastSendMsgTime,
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
	url, err := neturl.Parse(goInstaHost)
	if err != nil {
		return nil, err
	}
	urlb, err := neturl.Parse(goInstaHost_B)
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
		ID:                  config.ID,
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
		RegisterPhoneNumber: config.RegisterPhoneNumber,
		RegisterPhoneArea:   config.RegisterPhoneArea,
		RegisterIpCountry:   config.RegisterIpCountry,
		Status:              config.Status,
		UserAgent:           config.UserAgent,
		sessionID:           common.GenUUID(),
		LastSendMsgTime:     config.LastSendMsgTime,
		c: &http.Client{
			Jar: jar,
		},
	}
	if inst.UserAgent == "" {
		inst.UserAgent = GenUserAgent()
	}

	inst.Proxy = &proxy.Proxy{ID: config.ProxyID}
	common.DebugHttpClient(inst.c)

	return inst, nil
}

type UploadIDRecord struct {
	FileMd5  string `bson:"file_md5"`
	Username string `bson:"username"`
	FileType string `bson:"file_type"`
	FileName string `bson:"file_name"`
	UploadID string `bson:"upload_id"`
}

func SaveUploadID(record *UploadIDRecord) error {
	_, err := MogoHelper.UploadIDRecord.UpdateOne(
		context.TODO(),
		bson.M{"file_md5": record.FileMd5},
		bson.M{"$set": record},
		options.Update().SetUpsert(true))
	return err
}

func FindUploadID(username string, fileMd5 string) (*UploadIDRecord, error) {
	cursor, err := MogoHelper.UploadIDRecord.Find(context.TODO(),
		bson.D{{"$and",
			bson.D{
				{"username", username},
				{"file_md5", fileMd5},
			},
		}}, nil)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var result = &UploadIDRecord{}
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(result)
		return result, err
	}

	return nil, &common.MakeMoneyError{ErrType: common.NoMoreError}
}
