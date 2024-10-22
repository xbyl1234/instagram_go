package routine

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
)

var CrawlingDB *mongo.Database

var CrawTagsSearchColl *mongo.Collection
var CrawTagsTagColl *mongo.Collection
var CrawTagsMediaColl *mongo.Collection
var CrawTagsUserColl *mongo.Collection

func InitCrawTagsDB(taskName string) {
	CrawlingDB = goinsta.GetDB(taskName)
	CrawTagsTagColl = CrawlingDB.Collection("tags")
	CrawTagsMediaColl = CrawlingDB.Collection("media")
	CrawTagsUserColl = CrawlingDB.Collection("users")
	CrawTagsSearchColl = CrawlingDB.Collection("search")

	_, err := CrawTagsMediaColl.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"media_pk": 1,
			},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		log.Error("mongo create index error: %v", err)
	}
}

func CheckMedia(pk int64) bool {
	_, err := CrawTagsMediaColl.InsertOne(context.TODO(), bson.M{"media_pk": pk})
	return err == nil
}

var CrawFasDB *mongo.Database
var CrawFansUserColl *mongo.Collection
var CrawFansTargetUserColl *mongo.Collection

func InitCrawFansDB(taskName string, targetFansDBName string, targetFansCollName string) {
	CrawFasDB = goinsta.GetDB("inst_fans")
	CrawFansUserColl = CrawFasDB.Collection(taskName)
	CrawFansTargetUserColl = goinsta.GetDB(targetFansDBName).Collection(targetFansCollName)
}

var SendTargeUserColl *mongo.Collection
var ShareMediaLogColl *mongo.Collection

func InitSendMsgDB(TargetUserDB string, TargetUserCollection string, logColl string) {
	targetDB := goinsta.GetDB(TargetUserDB)
	SendTargeUserColl = targetDB.Collection(TargetUserCollection)

	logDB := goinsta.GetDB("instagram_log")
	ShareMediaLogColl = logDB.Collection(logColl)
	_, err := ShareMediaLogColl.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"media.pk": 1,
			},
			Options: options.Index().SetSparse(true),
		},
	)
	if err != nil {
		log.Error("mongo create index error: %v", err)
	}
}

type ShareMediaLog struct {
	Username string         `bson:"username"`
	Link     string         `bson:"link"`
	Media    *goinsta.Media `bson:"media"`
	Time     string         `bson:"time"`
}

func SaveShareMediaPk(pk int64) error {
	_, err := ShareMediaLogColl.InsertOne(context.TODO(), bson.M{"media": bson.M{"pk": pk}})
	if err != nil {
		return err
	}
	return nil
}

func SaveShareMediaLog(log *ShareMediaLog) error {
	log.Time = common.GetShanghaiTimeString()
	_, err := ShareMediaLogColl.UpdateOne(context.TODO(),
		bson.D{
			{"media.id", log.Media.Pk},
		}, bson.D{{"$set", log}}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	return nil
}

//func SaveTags(tags *goinsta.TagsFeed) error {
//	_, err := CrawTagsTagColl.UpdateOne(context.TODO(),
//		bson.D{
//			{"id", tags.Id},
//		}, bson.D{{"$set", tags}}, options.Update().SetUpsert(true))
//	if err != nil {
//		return err
//	}
//	return nil
//}

func LoadTags() ([]goinsta.TagsFeed, error) {
	cursor, err := CrawTagsTagColl.Find(context.TODO(), bson.M{"more_available": true}, nil)
	if err != nil {
		return nil, err
	}
	var tags []goinsta.TagsFeed
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

func LoadSearch() ([]*goinsta.Search, error) {
	cursor, err := CrawTagsSearchColl.Find(context.TODO(), bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	var search []*goinsta.Search
	cursor.All(context.TODO(), &search)
	_ = cursor.Close(context.TODO())
	return search, err
}

//
//func SaveMedia(mediaComb *MediaComb) error {
//	if mediaComb.Flag == "" {
//		if mediaComb.Media.CommentCount == 0 {
//			mediaComb.Flag = "no comment"
//		} else if mediaComb.Comments != nil && !mediaComb.Comments.HasMore {
//			mediaComb.Flag = "no comment"
//		}
//	}
//
//	_, err := CrawTagsMediaColl.UpdateOne(context.TODO(),
//		bson.D{
//			{"q", mediaComb.Media.ID},
//		}, bson.D{
//			{"$set", mediaComb}},
//		options.Update().SetUpsert(true))
//	return err
//}
//
//func LoadMedia(limit int) ([]MediaComb, error) {
//	cursor, err := CrawTagsMediaColl.Find(context.TODO(),
//		bson.D{{"media.comment_count", bson.M{"$gt": 0}},
//			{"$or", []bson.M{{"flag": bson.M{"$exists": false}},
//				{"flag": ""}}},
//		},
//		nil)
//	if err != nil {
//		return nil, err
//	}
//
//	var result = make([]MediaComb, limit)
//	index := 0
//	for cursor.Next(context.TODO()) && index < limit {
//		err = cursor.Decode(&result[index])
//		if err != nil {
//			break
//		}
//
//		_, _ = CrawTagsMediaColl.UpdateOne(context.TODO(),
//			bson.D{
//				{"q", result[index].Media.ID},
//			}, bson.D{
//				{"$set", bson.M{"flag": "loaded"}}},
//			options.Update().SetUpsert(true))
//
//		index++
//	}
//
//	_ = cursor.Close(context.TODO())
//
//	return result[:index], err
//}

type SendHistory struct {
	TaskName string `bson:"task_name"`
	HadRead  bool   `bson:"had_read"`
	HadOpen  bool   `bson:"had_open"`
}

type CrawData struct {
	Tag              string `json:"tag"`
	UserPk           int64  `json:"user_pk"`
	UserName         string `bson:"user_name"`
	MediaPk          int64  `json:"media_pk"`
	MediaId          string `bson:"media_id"`
	ParentCommentId  int64  `bson:"parent_comment_id"`
	LoggingInfoToken string `bson:"logging_info_token"`
	HadFollow        string `bson:"had_follow"`
	HadComment       string `bson:"had_comment"`
}

func UpdateCrawData(userComb *CrawData, key string, value string) error {
	_, err := SendTargeUserColl.UpdateOne(context.TODO(), bson.M{"user_pk": userComb.UserPk},
		bson.M{"$set": bson.M{key: value}},
		options.Update().SetUpsert(true))
	return err
}

//func SaveBlackUser(userComb *UserComb) error {
//	_, err := SendTargeUserColl.UpdateOne(context.TODO(), bson.M{"user.pk": userComb.User.ID},
//		bson.M{"$set": bson.M{"black": true}},
//		options.Update().SetUpsert(true))
//	return err
//}

func SaveUser(Coll *mongo.Collection, crawData *CrawData) error {
	_, err := Coll.UpdateOne(context.TODO(),
		bson.D{
			{"user", bson.M{"pk": crawData.UserPk}},
		}, bson.D{{"$set", crawData}}, options.Update().SetUpsert(true))
	return err
}

//func DelDup() {
//	cursor, err := SendTargeUserColl.Find(context.TODO(), bson.M{}, nil)
//	if err != nil {
//		return
//	}
//
//	set := make(map[int64]int)
//	var item UserComb
//	for cursor.Next(context.TODO()) {
//		err = cursor.Decode(&item)
//		if err != nil {
//			break
//		}
//		set[item.User.ID] = set[item.User.ID] + 1
//
//	}
//	_ = cursor.Close(context.TODO())
//
//	for k, v := range set {
//		if v > 1 {
//			for i := 0; i < v-1; i++ {
//				fmt.Println(k)
//				SendTargeUserColl.DeleteOne(context.TODO(), bson.M{"user.pk": k})
//			}
//		}
//	}
//}

func LoadCrawData(key []string, recvChan chan *CrawData) error {
	filter := []bson.M{
		{"$or": []bson.M{{"black": bson.M{"$exists": false}}, {"black": false}}},
	}
	for _, item := range key {
		filter = append(filter,
			bson.M{"$or": []bson.M{
				{
					item: bson.M{"$exists": false},
				},
				{
					item: nil,
				},
			}})
	}
	cursor, err := SendTargeUserColl.Find(context.TODO(),
		bson.M{"$and": filter}, nil)

	if err != nil {
		return err
	}

	for cursor.Next(context.TODO()) {
		var result *CrawData
		err = cursor.Decode(&result)
		if err != nil {
			break
		}
		recvChan <- result
	}
	_ = cursor.Close(context.TODO())
	return err
}

//func LoadFansTargetUser(limit int) ([]UserComb, error) {
//	cursor, err := CrawFansTargetUserColl.Find(context.TODO(),
//		bson.D{{"$or", []bson.M{{"followes": nil},
//			{"followes": bson.M{"has_more": true}}}}})
//	if err != nil {
//		return nil, err
//	}
//
//	var result = make([]UserComb, limit)
//	index := 0
//	for cursor.Next(context.TODO()) && index < limit {
//		err = cursor.Decode(&result[index])
//		if err != nil {
//			break
//		}
//		index++
//	}
//	_ = cursor.Close(context.TODO())
//
//	return result[:index], err
//}
