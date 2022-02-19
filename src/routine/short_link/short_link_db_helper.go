package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common/log"
	"time"
)

var Client *mongo.Client
var ShortLinkLog *mongo.Collection

func InitShortLinkDB(mogoUri string) {
	//"mongodb://xbyl:XBYLxbyl1234@62.216.92.183:27017"
	clientOptions := options.Client().ApplyURI(mogoUri)
	//clientOptions := options.Client().ApplyURI("mongodb://xbyl:xbyl741852JHK@192.168.187.1:27017")

	var err error

	Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error("mongo %v", err)
	}

	err = Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Error("mongo %v", err)
	}

	ShortLinkLog = Client.Database("make_money").Collection("short_link_log")
	//Account = Client.Database("make_money").Collection("account")
	//UploadIDRecord = Client.Database("make_money").Collection("upload_id")
}

type ShortLinkLogDB struct {
	TimeTick  int64
	Time      string
	UserID    string
	ShortLink string
	Url       string
	UA        string
	IP        string
}

func ShortLinkLog2DB(log *ShortLinkLogDB) error {
	log.TimeTick = time.Now().Unix()
	log.Time = time.Now().String()
	_, err := ShortLinkLog.InsertOne(context.TODO(), log)
	return err
}
