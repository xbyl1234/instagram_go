package main

import (
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
)

func CrawCommentUser(combChan chan *MediaComb, waitCraw *sync.WaitGroup) {
	defer waitCraw.Done()
	for item := range combChan {
		inst := goinsta.AccountPool.GetOneBlock(goinsta.OperNameCrawComment, config.AccountTag)
		err := crawCommentUser(inst, item.Media, item.Tag)
		if err != nil {
			log.Error("CrawCommentUser: account %s error %v", inst.User, err)
		}
		goinsta.AccountPool.ReleaseOne(inst)
	}
}

func crawCommentUser(inst *goinsta.Instagram, media *goinsta.Media, tag string) (errRet error) {
	if err := recover(); err != nil {
		log.Error("account: %s crawCommentUser panic error: %v", inst.User, err)
		errRet = err.(error)
	}
	putRedisCount := 0
	for true {
		comment := inst.NewComments(media.Id)
		result, err := comment.NextComments()
		if err != nil {
			return err
		}
		comments := result.GetAllComments()
		var commentData routine.CrawData
		for index := range comments {
			commentData.UserPk = comments[index].User.Pk
			commentData.UserName = comments[index].User.Username
			commentData.MediaId = media.Id
			commentData.MediaPk = media.Pk
			commentData.ParentCommentId = comments[index].Pk
			commentData.LoggingInfoToken = media.LoggingInfoToken
			commentData.Tag = tag

			err = routine.SaveUser(routine.CrawTagsUserColl, &commentData)
			if err != nil {
				return err
			}

			if putRedisCount < config.AddCommentCount {
				putRedisCount++
				err = mediaRedis.PutJson(tag, &commentData)
				if err != nil {
					log.Error("mediaRedis PutJson error %v", err)
				}
			}
		}
	}

	return nil
}
