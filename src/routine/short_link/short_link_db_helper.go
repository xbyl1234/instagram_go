package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common/log"
	"time"
)

var Client *mongo.Client
var ShortLinkLogColl *mongo.Collection
var BlackHistoryColl *mongo.Collection

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

	ShortLinkLogColl = Client.Database("make_money").Collection("short_link_log")
	BlackHistoryColl = Client.Database("make_money").Collection("black_ip")
}

type ShortLinkLogDB struct {
	TimeTick  int64  `bson:"time_tick"`
	Time      string `bson:"time"`
	UserID    string `bson:"user_id"`
	ShortLink string `bson:"short_link"`
	Url       string `bson:"url"`
	UA        string `bson:"ua"`
	IP        string `bson:"ip"`
	IsFb      bool   `bson:"is_fb"`
	IsBlack   bool   `bson:"is_black"`
}

func ShortLinkLog2DB(log *ShortLinkLogDB) error {
	log.TimeTick = time.Now().Unix()
	log.Time = time.Now().String()
	_, err := ShortLinkLogColl.InsertOne(context.TODO(), log)
	return err
}

func LoadBlackHistory() ([]*BlackHistory, error) {
	cursor, err := BlackHistoryColl.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	var item []*BlackHistory
	err = cursor.All(context.TODO(), &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func SaveBlackHistory(item *BlackHistory) error {
	_, err := BlackHistoryColl.UpdateOne(context.TODO(), bson.M{"ip": item.IP},
		bson.M{"$set": item},
		options.Update().SetUpsert(true))
	return err
}
