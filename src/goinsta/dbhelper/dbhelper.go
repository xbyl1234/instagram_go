package dbhelper

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type MogoDBHelper struct {
	client  *mongo.Client
	phone   *mongo.Collection
	account *mongo.Collection
}

var MogoHelper *MogoDBHelper = nil

func InitMogoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	var err error
	MogoHelper = &MogoDBHelper{}

	MogoHelper.client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = MogoHelper.client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	MogoHelper.phone = MogoHelper.client.Database("inst").Collection("phone")
	MogoHelper.account = MogoHelper.client.Database("inst").Collection("account")
}

func GetMogoHelper() *MogoDBHelper {
	if MogoHelper == nil {
		InitMogoDB()
	}
	return MogoHelper
}

func GetMogoPhoneConn() *mongo.Collection {
	return MogoHelper.phone
}

func GetMogoAccountConn() *mongo.Collection {
	return MogoHelper.account
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
	_, err := MogoHelper.phone.UpdateOne(context.TODO(),
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
	_, err := MogoHelper.phone.UpdateOne(context.TODO(),
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
	_, err := MogoHelper.account.UpdateOne(
		context.TODO(),
		bson.M{"username": account.Username},
		bson.M{"$set": account},
		options.Update().SetUpsert(true))
	return err
}

func LoadDBAllAccount() ([]AccountCookies, error) {
	cursor, err := MogoHelper.account.Find(context.TODO(), bson.M{}, nil)
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
