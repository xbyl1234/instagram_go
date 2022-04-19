package main

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
)

func RunAddComment(coro int, tag []string) {
	queue, _ := common.CreateQueue(&routine.DBConfig.Redis)
	recvChan := make(chan *routine.CrawData, 1)
	for _, item := range tag {
		go func(tag string) {
			for true {
				get, err := queue.BLGet(tag)
				if err != nil {
					continue
				}
				data := &routine.CrawData{}
				json.Unmarshal([]byte(get), data)
				recvChan <- data
			}
		}(item)
	}
	wait := &sync.WaitGroup{}
	wait.Add(coro)

	for i := 0; i < coro; i++ {
		go func() {
			defer wait.Done()
			for item := range recvChan {
				inst := goinsta.AccountPool.GetOneBlock(goinsta.OperNameComment, config.AccountTag)
				err := addComment(inst, item)
				if err == nil {
					routine.UpdateCrawData(item, "had_comment", inst.User)
					if !inst.IsSpeedLimit(goinsta.OperNameLikeUser) {
						err = followUser(inst, item)
						if err != nil {
							log.Error("account: %s followUser error: %v", inst.User, err)
						}
						routine.UpdateCrawData(item, "had_follow", inst.User)
					}
				} else {
					routine.UpdateCrawData(item, "had_comment", err.Error())
					log.Error("account: %s AddComment error: %v", inst.User, err)
				}
				goinsta.AccountPool.ReleaseOne(inst)
			}
		}()
	}

	wait.Wait()
}

func addComment(inst *goinsta.Instagram, data *routine.CrawData) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("account: %s addComment panic error: %v", inst.User, err)
			retErr = err.(error)
		}
	}()
	ParentCommentId := ""
	if data.ParentCommentId != 0 {
		ParentCommentId = fmt.Sprintf("%d", data.ParentCommentId)
	}
	err := inst.GetUserOperate().AddComment(&goinsta.AddCommentParams{
		ParentCommentId:  ParentCommentId,
		UserName:         data.UserName,
		LoggingInfoToken: data.LoggingInfoToken,
		MediaId:          data.MediaId,
		CommentText:      getComments(),
	})

	return err
}

func RunFollow(coro int) {
	recvChan := make(chan *routine.CrawData, 10)
	go routine.LoadCrawData([]string{"had_follow"}, recvChan)

	wait := &sync.WaitGroup{}
	wait.Add(coro)

	go func() {
		for item := range recvChan {
			inst := goinsta.AccountPool.GetOneBlock(goinsta.OperNameLikeUser, config.AccountTag)
			err := followUser(inst, item)
			if err != nil {
				log.Error("account: %s RunFollow error: %v", inst.User, err)
				routine.UpdateCrawData(item, "had_follow", err.Error())
			} else {
				routine.UpdateCrawData(item, "had_follow", inst.User)
			}
			goinsta.AccountPool.ReleaseOne(inst)
		}
	}()
}

func followUser(inst *goinsta.Instagram, data *routine.CrawData) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("account: %s followUser panic error: %v", inst.User, err)
			retErr = err.(error)
		}
	}()
	err := inst.GetUserOperate().LikeUser(data.UserPk)
	return err
}

func RunPost() {

}

func RunMsg() {

}
