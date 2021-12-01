package routine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/goinsta"
)

var CrawlingDB *mongo.Database

var SearchCollection *mongo.Collection
var TagCollection *mongo.Collection
var MediaCollection *mongo.Collection
var UserCollection *mongo.Collection

func InitRoutineCrawDB(taskName string) {
	CrawlingDB = common.GetDB(taskName)
	TagCollection = CrawlingDB.Collection("tags")
	MediaCollection = CrawlingDB.Collection("media")
	UserCollection = CrawlingDB.Collection("users")
	SearchCollection = CrawlingDB.Collection("search")
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

func LoadTags() ([]goinsta.Tags, error) {
	cursor, err := TagCollection.Find(context.TODO(), bson.M{"moreavailable": true}, nil)
	if err != nil {
		return nil, err
	}
	var tags []goinsta.Tags
	err = cursor.All(context.TODO(), &tags)
	return tags, err
}

func SaveSearch(search *goinsta.Search) error {
	_, err := SearchCollection.UpdateOne(context.TODO(),
		bson.D{
			{"q", search.Q},
		}, bson.D{{"$set", search}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadSearch() (*goinsta.Search, error) {
	cursor, err := SearchCollection.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var search *goinsta.Search
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(&search)
	} else {
		return nil, nil
	}
	return search, err
}

type MediaComb struct {
	Media    *goinsta.Item     `json:"media"`
	Comments *goinsta.Comments `json:"comments"`
	Tag      string            `json:"tag"`
}

func SaveMedia(mediaComb *MediaComb) error {
	_, err := MediaCollection.UpdateOne(context.TODO(),
		bson.D{
			{"q", mediaComb.Media.ID},
		}, bson.D{
			{"$set", mediaComb}},
		options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

func LoadMedia() ([]MediaComb, error) {
	cursor, err := MediaCollection.Find(context.TODO(),
		bson.D{{"$or",
			bson.D{{"comments", nil},
				{"comments", bson.M{"hasmore": true}}}}},
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

type UserComb struct {
	User   *goinsta.User `json:"user"`
	Source string        `json:"source"`
}

func SaveUser(userComb *UserComb) error {
	_, err := UserCollection.UpdateOne(context.TODO(),
		bson.D{
			{"user", bson.M{"pk": userComb.User.ID}},
		}, bson.D{{"$set", userComb}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadUser(tag string, sendTaskName string, limit int) ([]UserComb, error) {
	cursor, err := UserCollection.Find(context.TODO(),
		bson.D{{"$and",
			bson.D{{"tag", tag},
				{sendTaskName,
					bson.M{"$exists": false}}}}},
		nil)

	if err != nil {
		return nil, err
	}

	var userCombs = make([]UserComb, limit)
	var index = 0
	for index = range userCombs {
		if cursor.Next(context.TODO()) {
			err = cursor.Decode(&userCombs[index])
			if err != nil {
				break
			}
		} else {
			break
		}
	}

	return userCombs[:index], err
}

//func MarkUser(userComb *UserComb, markName string) error {
//
//}
