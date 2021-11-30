package main

import (
	"container/list"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/goinsta"
)

var CrawlingDB *mongo.Database
var TagCollection *mongo.Collection
var MediaCollection *mongo.Collection

func InitCrawDB(taskName string) {
	CrawlingDB = common.GetDB(taskName)
	TagCollection = CrawlingDB.Collection("tags")
	MediaCollection = CrawlingDB.Collection("media")
}

func SaveTags(tags *goinsta.Tags) error {
	_, err := TagCollection.UpdateOne(context.TODO(),
		bson.D{
			{"id", tags.Id},
		}, bson.D{{"$set", tags}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadTags() (*list.List, error) {
	var ret = list.New()
	cursor, err := TagCollection.Find(context.TODO(), bson.M{"more_available": true}, nil)
	if err != nil {
		return ret, err
	}
	var tags []goinsta.Tags
	err = cursor.All(context.TODO(), &tags)
	for index := range tags {
		ret.PushBack(&tags[index])
	}
	return ret, err
}

func SaveSearch(search *goinsta.Search) error {
	_, err := TagCollection.UpdateOne(context.TODO(),
		bson.D{
			{"q", search.Q},
		}, bson.D{{"$set", search}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadSearch() (*goinsta.Search, error) {
	cursor, err := TagCollection.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var search *goinsta.Search
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(&search)
	} else {
		err = common.MakeMoneyError_NoMore
	}
	return search, err
}

func SaveMedia(media *goinsta.Item, comments *goinsta.Comments) error {
	_, err := MediaCollection.UpdateOne(context.TODO(),
		bson.D{
			{"q", media.ID},
		}, bson.D{
			{"$set", bson.D{{"media", media},
				{"comments", comments}}}},
		options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

type MediaComb struct {
	Media    *goinsta.Item     `json:"media"`
	Comments *goinsta.Comments `json:"comments"`
}

func LoadMedia() ([]MediaComb, error) {
	cursor, err := MediaCollection.Find(context.TODO(),
		bson.D{{"$or",
			bson.D{{"comments", nil},
				{"comments", bson.M{"has_more": true}}}}},
		nil)

	if err != nil {
		return nil, err
	}

	var result []MediaComb
	err = cursor.All(context.TODO(), &result)
	if err != nil {
		return nil, err
	}

	return result[:], err
}

func SaveUser(user *goinsta.User) error {
	_, err := TagCollection.UpdateOne(context.TODO(),
		bson.D{
			{"q", search.Q},
		}, bson.D{{"$set", search}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadUser() (*goinsta.Search, error) {
	cursor, err := TagCollection.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var search *goinsta.Search
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(&search)
	} else {
		err = common.MakeMoneyError_NoMore
	}
	return search, err
}
