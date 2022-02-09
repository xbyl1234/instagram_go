package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"strings"
	"time"
)

//816
//1320
func CrawCommentUser() {
	var currAccount *goinsta.Instagram

	var SetNewAccount = func(mediaComb *routine.MediaComb, inst *goinsta.Instagram) {
		mediaComb.Media.SetAccount(inst)
		if mediaComb.Comments == nil {
			mediaComb.Comments = mediaComb.Media.GetComments()
		} else {
			mediaComb.Comments.SetAccount(inst)
		}
	}

	var RequireAccont = func(mediaComb *routine.MediaComb) {
		if currAccount == nil || currAccount.IsSpeedLimit("craw_comment_user") || currAccount.IsBad() {
			var oldUser string
			if currAccount != nil {
				oldUser = currAccount.User
				goinsta.AccountPool.ReleaseOne(currAccount)
			}

			inst := routine.ReqAccount("craw_comment_user")
			if inst == nil {
				log.Error("CrawCommentUser req account error!")
			}
			currAccount = inst

			SetNewAccount(mediaComb, currAccount)
			log.Warn("CrawCommentUser replace account to %s->%s", oldUser, currAccount.User)
			return
		} else {
			SetNewAccount(mediaComb, currAccount)
			return
		}
	}

	defer func() {
		if currAccount != nil {
			goinsta.AccountPool.ReleaseOne(currAccount)
		}
	}()

	for mediaComb := range MediaChan {
		if mediaComb.Media.CommentCount == 0 {
			mediaComb.Comments.HasMore = false
			mediaComb.Flag = "no comment"
			routine.SaveMedia(mediaComb)
			continue
		}
		log.Info("craw_comment_user %s", mediaComb.Media.ID)

		for true {
			RequireAccont(mediaComb)
			respComm, err := mediaComb.Comments.NextComments()

			_, min, hour, day := currAccount.GetSpeed("craw_comment_user")
			log.Info("account %s craw_comment_user count %d,%d,%d status %s", currAccount.User, min, hour, day, currAccount.Status)
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("media %s comments has craw finish!", mediaComb.Media.ID)
					mediaComb.Comments.HasMore = false
					mediaComb.Flag = "finish"
					routine.SaveMedia(mediaComb)
					break
				} else if common.IsError(err, common.ChallengeRequiredError) ||
					common.IsError(err, common.FeedbackError) ||
					common.IsError(err, common.LoginRequiredError) {
					log.Error("user %s status is %s from CrawCommentUser task, err: %v", currAccount.User,
						currAccount.Status, err)
					RequireAccont(mediaComb)
					continue
				} else if common.IsError(err, common.RequestError) {
					log.Warn("CrawCommentUser retrying...user: %s, err: %v", currAccount.User, err)
					continue
				} else if strings.Index(err.Error(), "Media is unavailable") >= 0 {
					log.Warn("Media %d is unavailable", mediaComb.Media.ID)
					mediaComb.Comments.HasMore = false
					mediaComb.Flag = "unavailable"
					routine.SaveMedia(mediaComb)
					break
				} else {
					log.Error("NextComments unknow error:%v", err)
					continue
				}
			}

			comments := respComm.GetAllComments()
			var userComb routine.UserComb
			for index := range comments {
				userComb.User = &comments[index].User
				userComb.Source = "comments"
				err = routine.SaveUser(routine.CrawTagsUserColl, &userComb)
				if err != nil {
					log.Error("SaveUser error:%v", err)
					break
				}
			}
		}
	}
}

func SendMedias() {
	for true {
		medias, err := routine.LoadMedia(100)
		if err != nil {
			log.Error("load media error: %v", err)
			continue
		}
		if len(medias) == 0 {
			time.Sleep(time.Second * 60)
			continue
		}
		for index := range medias {
			MediaChan <- &medias[index]
		}
	}
	close(MediaChan)
}
