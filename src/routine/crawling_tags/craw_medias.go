package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
	"time"
)

//830
//847
func CrawMedias(TagsChan chan *goinsta.Tags, waitCraw *sync.WaitGroup, StopTime time.Time) {
	defer waitCraw.Done()
	var currAccount *goinsta.Instagram
	var SetNewAccount = func(tag *goinsta.Tags) {
		inst := routine.ReqAccount("craw_medias")
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
			if currAccount.IsSpeedLimit("craw_medias") || tag.Inst.IsBad() {
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
		}
	}()

	for tag := range TagsChan {
		for true {
			RequireAccount(tag)
			tagResult, err := tag.Next()
			_, min, hour, day := currAccount.GetSpeed("craw_medias")
			log.Info("account %s craw_medias count %d,%d,%d status %s", currAccount.User, min, hour, day, currAccount.Status)
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("tags %s medias has craw finish!", tag.Name)
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
						log.Info("%d %d", mediaComb.Media.Caption.CreatedAt, StopTime.Unix())
						log.Info("craw media stop! media time is %s %s", mediaTime.Format("2006-01-02 15:04:05"), StopTime.Format("2006-01-02 15:04:05"))
					} else {
						log.Info("craw media current time is %s", mediaTime.Format("2006-01-02 15:04:05"))
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
					break
				}
			}

			//err = routine.SaveTags(tag)
			//if err != nil {
			//	log.Error("SaveTags error:%v", err)
			//}
			if stop {
				break
			}
		}

		goinsta.AccountPool.ReleaseOne(currAccount)
	}
}

func SendTags(TagsChan chan *goinsta.Tags) {
	for true {
		for item := TagList.Front(); item != nil; item = item.Next() {
			tags := item.Value.(*goinsta.Tags)
			if tags.MoreAvailable {
				TagsChan <- tags
			}
		}
	}
}
