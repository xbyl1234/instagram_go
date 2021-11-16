package dbhelper

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
		}, bson.D{{"$inc", bson.M{"register_count": 1}},
		}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func InsertNewAccount() {

}
