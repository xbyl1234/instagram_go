package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
	"sync/atomic"
	"time"
)

var indCoro int32 = 0

var CrawMediaAccountTag = "craw_media"

//830
//847
func CrawMedias(TagsChan chan *goinsta.Tags, waitCraw *sync.WaitGroup, StopTime time.Time) {
	defer waitCraw.Done()
	var currAccount *goinsta.Instagram
	var SetNewAccount = func(tag *goinsta.Tags) {
		inst := routine.ReqAccount(goinsta.OperNameCrawMedia, CrawMediaAccountTag)
		if inst == nil {
			log.Error("CrawMedias req account error")
			return
		}
		tag.SetAccount(inst)
		currAccount = inst
		err := tag.Sync(goinsta.TabRecent)
		if err != nil {
			log.Error("tag sync error: %v", err)
		}
		_, err = tag.Stories()
		if err != nil {
			log.Error("tag stories error: %v", err)
		}
	}

	var RequireAccount = func(tag *goinsta.Tags) {
		var oldUser string
		if tag.Inst == nil {
			SetNewAccount(tag)
			log.Info("CrawMedias set account %s", tag.Inst.User)
		} else {
			oldUser = tag.Inst.User
			if currAccount.IsSpeedLimit(goinsta.OperNameCrawMedia) || tag.Inst.IsBad() {
				goinsta.AccountPool.ReleaseOne(tag.Inst)
				SetNewAccount(tag)
				log.Warn("CrawMedias replace account %s->%s", oldUser, tag.Inst.User)
				return
			} else {
				return
			}
		}
	}
	defer func() {
		if currAccount != nil {
			goinsta.AccountPool.ReleaseOne(currAccount)
			currAccount = nil
		}
	}()

	myIdx := atomic.AddInt32(&indCoro, 1)

	for tag := range TagsChan {
		log.Info("coro %d will craw media %s", myIdx, tag.Name)

		for true {
			RequireAccount(tag)
			tagResult, err := tag.Next()
			_, min, hour, day := currAccount.GetSpeed(goinsta.OperNameCrawMedia)
			log.Info("coro %d account %s craw media %s count %d,%d,%d status %s", myIdx, currAccount.User, tag.Name, min, hour, day, currAccount.Status)
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("coro %d tags %s medias has craw finish!", myIdx, tag.Name)
					break
				} else if common.IsError(err, common.ChallengeRequiredError) ||
					common.IsError(err, common.FeedbackError) ||
					common.IsError(err, common.LoginRequiredError) {
					log.Error("user %s status is %s from CrawCommentUser task, err: %v", currAccount.User,
						currAccount.User, err)
					RequireAccount(tag)
					continue
				} else if common.IsError(err, common.RequestError) {
					log.Warn("CrawMedias retrying...user: %s, err: %v", currAccount.User, err)
					continue
				} else {
					log.Error("next media unknow error: %v", err)
					continue
				}
			}

			var OncePrint = false
			var stop = false
			medias := tagResult.GetAllMedias()
			var mediaComb routine.MediaComb
			for index := range medias {
				mediaComb.Media = medias[index]
				mediaComb.Tag = tag.Name
				if !OncePrint {
					mediaTime := time.Unix(mediaComb.Media.Caption.CreatedAt, 0)
					if mediaComb.Media.Caption.CreatedAt < StopTime.Unix() {
						stop = true
						tag.MoreAvailable = false
						log.Info("coro %d account %s craw media %s stop! media time is %s", myIdx, currAccount.User, tag.Name, mediaTime.Format("2006-01-02 15:04:05"))
					} else {
						log.Info("coro %d account %s craw media %s  current time is %s", myIdx, currAccount.User, tag.Name, mediaTime.Format("2006-01-02 15:04:05"))
					}
					OncePrint = true
				}

				err = routine.SaveMedia(&mediaComb)
				if err != nil {
					log.Error("SaveMedia error:%v", err)
				}

				var userComb routine.UserComb
				userComb.User = &medias[index].User
				userComb.Source = "media"
				err = routine.SaveUser(routine.CrawTagsUserColl, &userComb)
				if err != nil {
					log.Error("SaveUser error:%v", err)
				}
			}

			//err = routine.SaveTags(tag)
			//if err != nil {
			//	log.Error("SaveTags error:%v", err)
			//}
			if stop {
				log.Info("coro %d break", myIdx)
				break
			}
		}

		log.Info("coro %d craw media %s out! account: %s", myIdx, tag.Name, currAccount.User)
		goinsta.AccountPool.ReleaseOne(currAccount)
	}

	log.Info("coro %d exit!", myIdx)
}

func SendTags(TagsChan chan *goinsta.Tags) {
	for item := TagList.Front(); item != nil; item = item.Next() {
		tags := item.Value.(*goinsta.Tags)
		if tags.MoreAvailable {
			TagsChan <- tags
		}
	}
}
