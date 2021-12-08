package common

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MogoDBHelper struct {
	Client  *mongo.Client
	Phone   *mongo.Collection
	Account *mongo.Collection
}

var MogoHelper *MogoDBHelper = nil

func InitMogoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://xbyl:XBYLxbyl1234@62.216.92.183:27017")
	//clientOptions := options.Client().ApplyURI("mongodb://xbyl:xbyl741852JHK@192.168.187.1:27017")

	var err error
	MogoHelper = &MogoDBHelper{}

	MogoHelper.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = MogoHelper.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	MogoHelper.Phone = MogoHelper.Client.Database("inst").Collection("phone")
	MogoHelper.Account = MogoHelper.Client.Database("inst").Collection("account")
}
func GetDB(name string) *mongo.Database {
	return MogoHelper.Client.Database(name)
}

func GetMogoHelper() *MogoDBHelper {
	if MogoHelper == nil {
		InitMogoDB()
	}
	return MogoHelper
}

func GetMogoPhoneConn() *mongo.Collection {
	return MogoHelper.Phone
}

func GetMogoAccountConn() *mongo.Collection {
	return MogoHelper.Account
}
