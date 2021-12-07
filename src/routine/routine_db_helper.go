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

var CrawTagsSearchColl *mongo.Collection
var CrawTagsTagColl *mongo.Collection
var CrawTagsMediaColl *mongo.Collection
var CrawTagsUserColl *mongo.Collection

func InitCrawTagsDB(taskName string) {
	CrawlingDB = common.GetDB(taskName)
	CrawTagsTagColl = CrawlingDB.Collection("tags")
	CrawTagsMediaColl = CrawlingDB.Collection("media")
	CrawTagsUserColl = CrawlingDB.Collection("users")
	CrawTagsSearchColl = CrawlingDB.Collection("search")
}

var CrawFansUserColl *mongo.Collection
var CrawFansTargetUserColl *mongo.Collection

func InitCrawFansDB(taskName string, targetFansDBName string, targetFansCollName string) {
	CrawFansUserColl = common.GetDB("inst_fans").Collection(taskName)
	CrawFansTargetUserColl = common.GetDB(targetFansDBName).Collection(targetFansCollName)
}

func SaveTags(tags *goinsta.Tags) error {
	_, err := CrawTagsTagColl.UpdateOne(context.TODO(),
		bson.D{
			{"id", tags.Id},
		}, bson.D{{"$set", tags}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadTags() ([]goinsta.Tags, error) {
	cursor, err := CrawTagsTagColl.Find(context.TODO(), bson.M{"moreavailable": true}, nil)
	if err != nil {
		return nil, err
	}
	var tags []goinsta.Tags
	err = cursor.All(context.TODO(), &tags)
	_ = cursor.Close(context.TODO())
	return tags, err
}

func SaveSearch(search *goinsta.Search) error {
	_, err := CrawTagsSearchColl.UpdateOne(context.TODO(),
		bson.D{
			{"q", search.Q},
		}, bson.D{{"$set", search}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadSearch() (*goinsta.Search, error) {
	cursor, err := CrawTagsSearchColl.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var search *goinsta.Search
	if cursor.Next(context.TODO()) {
		err = cursor.Decode(&search)
	}
	_ = cursor.Close(context.TODO())
	return search, err
}

type MediaComb struct {
	Media    *goinsta.Item     `json:"media"`
	Comments *goinsta.Comments `json:"comments"`
	Tag      string            `json:"tag"`
}

func SaveMedia(mediaComb *MediaComb) error {
	_, err := CrawTagsMediaColl.UpdateOne(context.TODO(),
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

func SaveComments(mediaComb *MediaComb) error {
	_, err := CrawTagsMediaColl.UpdateOne(context.TODO(),
		bson.D{
			{"q", mediaComb.Media.ID},
		}, bson.D{
			{"$set", bson.M{"comments": mediaComb.Comments}}},
		options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

func LoadMedia(limit int) ([]MediaComb, error) {
	cursor, err := CrawTagsMediaColl.Find(context.TODO(),
		bson.D{{"$or", []bson.M{bson.M{"comments": nil},
			bson.M{"comments": bson.M{"hasmore": true}}}}},
		nil)
	if err != nil {
		return nil, err
	}

	var result []MediaComb = make([]MediaComb, limit)
	index := 0
	for cursor.Next(context.TODO()) && index < limit {
		err = cursor.Decode(&result[index])
		if err != nil {
			break
		}
		index++
	}
	_ = cursor.Close(context.TODO())

	return result[:index], err
}

type UserComb struct {
	User     *goinsta.User      `json:"user"`
	Source   string             `json:"source"`
	Followes *goinsta.Followers `json:"followes"`
}

func SaveUser(Coll *mongo.Collection, userComb *UserComb) error {
	_, err := Coll.UpdateOne(context.TODO(),
		bson.D{
			{"user", bson.M{"pk": userComb.User.ID}},
		}, bson.D{{"$set", userComb}}, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func LoadUser(tag string, sendTaskName string, limit int) ([]UserComb, error) {
	cursor, err := CrawTagsUserColl.Find(context.TODO(),
		bson.D{{"$and",
			bson.D{{"tag", tag},
				{sendTaskName,
					bson.M{"$exists": false}}}}},
		nil)

	if err != nil {
		return nil, err
	}

	var result = make([]UserComb, limit)
	index := 0
	for cursor.Next(context.TODO()) && index < limit {
		err = cursor.Decode(&result[index])
		if err != nil {
			break
		}
		index++
	}
	_ = cursor.Close(context.TODO())

	return result[:index], err
}

func LoadFansTargetUser(limit int) ([]UserComb, error) {
	cursor, err := CrawFansTargetUserColl.Find(context.TODO(),
		bson.D{{"$or", []bson.M{bson.M{"followes": nil},
			bson.M{"followes": bson.M{"hasmore": true}}}}})
	if err != nil {
		return nil, err
	}

	var result = make([]UserComb, limit)
	index := 0
	for cursor.Next(context.TODO()) && index < limit {
		err = cursor.Decode(&result[index])
		if err != nil {
			break
		}
		index++
	}
	_ = cursor.Close(context.TODO())

	return result[:index], err
}
