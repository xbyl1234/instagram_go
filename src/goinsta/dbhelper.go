package goinsta

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"net/http/cookiejar"
	neturl "net/url"
	"strings"
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
	ID           int64                    `json:"id" bson:"id"`
	Username     string                   `json:"username" bson:"username"`
	Passwd       string                   `json:"passwd" bson:"passwd"`
	HttpHeader   map[string]string        `json:"http_header" bson:"http_header"`
	ProxyID      string                   `json:"proxy_id" bson:"proxy_id"`
	IsLogin      bool                     `json:"is_login" bson:"is_login"`
	Token        string                   `json:"token" bson:"token"`
	Cookies      []*http.Cookie           `json:"cookies" bson:"cookies"`
	CookiesB     []*http.Cookie           `json:"cookies_b" bson:"cookies_b"`
	AccountInfo  *InstAccountInfo         `json:"account_info" bson:"account_info"`
	Status       string                   `json:"status" bson:"status"`
	Tags         string                   `json:"tags" bson:"tags"`
	SpeedControl map[string]*SpeedControl `json:"speed_control" bson:"speed_control"`
}

func SaveNewAccount(account AccountCookies) error {
	_, err := MogoHelper.Account.UpdateOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account},
		options.Update().SetUpsert(true))
	return err
}

func ReplaceAccount(account AccountCookies) error {
	_, err := MogoHelper.Account.ReplaceOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account})
	return err
}

func LoadDBAccountByTags(tags []string) ([]AccountCookies, error) {
	filter := make([]bson.M, len(tags))
	for idx := range tags {
		filter[idx] = make(bson.M)
		filter[idx]["tags"] = tags[idx]
	}
	cursor, err := MogoHelper.Account.Find(context.TODO(), bson.M{"$or": filter}, nil)
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
	return SaveNewAccount(Cookies)
}

func ConvConfig(config *AccountCookies) (*Instagram, error) {
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

	c := common.CreateGoHttpClient(common.DefaultHttpTimeout())
	c.Jar = jar

	inst := &Instagram{
		ID:           config.ID,
		User:         config.Username,
		Pass:         config.Passwd,
		token:        config.Token,
		AccountInfo:  config.AccountInfo,
		httpHeader:   config.HttpHeader,
		IsLogin:      config.IsLogin,
		Status:       config.Status,
		SpeedControl: config.SpeedControl,
		sessionID:    strings.ToUpper(common.GenUUID()),
		c:            c,
		Tags:         config.Tags,
	}

	if inst.AccountInfo == nil {
		inst.AccountInfo = GenInstDeviceInfo()
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

func LoadAccountByTags(tags []string) []*Instagram {
	config, err := LoadDBAccountByTags(tags)
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

type UploadIDRecord struct {
	FileMd5  string `bson:"file_md5"`
	UserID   int64  `bson:"user_id"`
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

func LoadUploadID(userID int64) ([]UploadIDRecord, error) {
	cursor, err := MogoHelper.UploadIDRecord.Find(context.TODO(),
		bson.D{{"user_id", userID}}, nil)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var result []UploadIDRecord
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
