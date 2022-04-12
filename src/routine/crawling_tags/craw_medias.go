package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
	"time"
)

func CrawMedias(tag string, mediaChan chan *MediaComb, waitCraw *sync.WaitGroup, StopTime time.Time) (retErr error) {
	defer waitCraw.Done()
	inst := goinsta.AccountPool.GetOneBlock(goinsta.OperNameCrawMedia, config.AccountTag)
	err := crawMedias(inst, tag, mediaChan, StopTime)
	if err != nil {
		log.Error("CrawMedias:tag %s account %s error %v", tag, inst.User, err)
		return err
	}
	return nil
}

type MediaComb struct {
	Media *goinsta.Media `bson:"media"`
	Tag   string         `bson:"tag"`
}

func crawMedias(inst *goinsta.Instagram, tag string, mediaChan chan *MediaComb, StopTime time.Time) (retErr error) {
	if err := recover(); err != nil {
		log.Error("account: %s doDevelopMeta panic error: %v", inst.User, err)
		retErr = err.(error)
	}

	feed := inst.NewTagsFeed(tag, goinsta.TabRecent)

	for true {
		if inst.IsSpeedLimit(goinsta.OperNameCrawMedia) {
			if inst.IsSpeedLimitInDay(goinsta.OperNameCrawMedia) {
				return &common.MakeMoneyError{ErrStr: "speed limit in day"}
			}
			cool := inst.GetCoolTime(goinsta.OperNameCrawMedia)
			if cool > 0 {
				time.Sleep(cool)
			}
		}

		next, err := feed.Next()
		if err != nil {
			return err
		}

		var OncePrint = false
		var stop = false
		medias := next.GetAllMedias()

		for index := range medias {
			if medias[index].CommentingDisabledForViewer {
				continue
			}
			mediaComb := &MediaComb{}
			mediaComb.Media = medias[index]
			mediaComb.Tag = feed.Name

			if !OncePrint {
				mediaTime := time.Unix(mediaComb.Media.Caption.CreatedAt, 0)
				if mediaComb.Media.Caption.CreatedAt < StopTime.Unix() {
					stop = true
					log.Info("account %s craw media %s stop! media time is %s", feed.Inst.User, feed.Name, mediaTime.Format("2006-01-02 15:04:05"))
				} else {
					log.Info("account %s craw media %s  current time is %s", feed.Inst.User, feed.Name, mediaTime.Format("2006-01-02 15:04:05"))
				}
				OncePrint = true
			}
			if !routine.CheckMedia(mediaComb.Media.Pk) {
				continue
			}

			mediaChan <- mediaComb

			var crawData routine.CrawData
			crawData.MediaId = medias[index].Id
			crawData.MediaPk = medias[index].Pk
			crawData.UserName = medias[index].User.Username
			crawData.UserPk = medias[index].User.Pk
			crawData.Tag = feed.Name
			err = routine.SaveUser(routine.CrawTagsUserColl, &crawData)
			if err != nil {
				return err
			}

			err = mediaRedis.PutJson(tag, &crawData)
			if err != nil {
				log.Error("mediaRedis PutJson error %v", err)
			}
		}

		if stop {
			log.Info("coro %d break")
			break
		}
	}

	return nil
}
