package main

import (
	"container/list"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"sync"
	"time"
)

type DevelopAccountConfig struct {
	MaxSubComment int      `json:"max_sub_comment"`
	Comments      []string `json:"comments"`
	NeedPushImg   bool     `json:"need_push_img"`
	NeedPushVideo bool     `json:"need_push_video"`
	NeedLike      bool     `json:"need_like"`
	FeedBackSleep int      `json:"feed_back_sleep"`
	LogCollName   string   `json:"log_coll_name"`
	//Spec          string   `json:"spec"`
}

func LogMedia(inst *goinsta.Instagram, link string, media *goinsta.Media) {
	routine.SaveShareMediaLog(&routine.ShareMediaLog{
		Username: inst.User,
		Link:     link,
		Media:    media,
	})
}

var emoji = "❤️🙌🔥👏😢😍😮😂"
var metaList *list.List
var badInstList *list.List
var metaListLock sync.Mutex
var waitAll sync.WaitGroup
var laterError = &common.MakeMoneyError{ErrStr: "later"}

func getComments() string {
	comment := config.Develop.Comments[common.GenNumber(0, len(config.Develop.Comments))]
	str := ""
	for i := 0; i < 5; i++ {
		str += string(emoji[common.GenNumber(0, len(emoji))])
	}
	return fmt.Sprintf(comment, str)
}

type DevelopMeta struct {
	inst      *goinsta.Instagram
	feed      *goinsta.VideoFeed
	comments  *goinsta.Comments
	followSet map[int64]bool

	curVideoList    *goinsta.VideosFeedResp
	curComments     *goinsta.RespComments
	nextVideoIdx    int
	nextCommentIdx  int
	subCommentCount int

	addSubCommentFinish bool
	hadShareMedia       bool
	hadCheckMedia       bool

	lastFeedBackTime time.Time
	isRunning        bool
}

func feedVideo(meta *DevelopMeta) error {
	meta.comments = nil
	meta.curVideoList = nil
	meta.curComments = nil
	meta.nextVideoIdx = 0
	meta.nextCommentIdx = 0
	meta.addSubCommentFinish = false
	meta.hadCheckMedia = false

	if meta.inst.IsSpeedLimit(goinsta.OperNameFeedVideo) {
		return laterError
	}
	var err error
	meta.curVideoList, err = meta.feed.Next()
	if err != nil {
		meta.inst.ResetProxy()
		log.Error("account: %s feedVideo.Next error: %v", meta.inst.User, err)
		return err
	}
	if len(meta.curVideoList.Items) == 0 {
		meta.inst.ResetProxy()
		log.Error("account: %s not feedVideo any", meta.inst.User)
		return laterError
	}
	return nil
}

func doDevelopMeta(meta *DevelopMeta) (retErr error) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("account: %s doDevelopMeta panic error: %v", meta.inst.User, err)
			retErr = err.(error)
		}
	}()

	opt := meta.inst.GetUserOperate()
	var err error

	for true {
		if meta.curVideoList == nil || meta.nextVideoIdx >= len(meta.curVideoList.Items) {
			err = feedVideo(meta)
			if err != nil {
				return err
			}
		}

		for ; meta.nextVideoIdx < len(meta.curVideoList.Items); meta.nextVideoIdx++ {
			media := meta.curVideoList.Items[meta.nextVideoIdx]
			if media.Media.CommentingDisabledForViewer {
				continue
			}
			if !meta.hadCheckMedia {
				meta.hadCheckMedia = true
				err = routine.SaveShareMediaPk(media.Media.Pk)
				if err != nil {
					log.Warn("repeated pk: %d", media.Media.Pk)
					continue
				}
			}

			if !meta.hadShareMedia {
				var shareMedia string
				shareMedia, err = opt.ShareMedia(media.Media.Id)
				if err != nil {
					log.Error("account: %s ShareMedia error: %v", meta.inst.User, err)
				} else {
					LogMedia(meta.inst, shareMedia, media.Media)
				}
				meta.hadShareMedia = true
			}

			if media.Media.HasMoreComments {
				if meta.comments == nil {
					meta.comments = meta.inst.NewComments(media.Media.Id)
				}

				for !meta.addSubCommentFinish {
					if meta.curComments == nil {
						if meta.inst.IsSpeedLimit(goinsta.OperNameCrawComment) {
							return laterError
						}

						meta.curComments, err = meta.comments.NextComments()
						if err != nil {
							log.Error("account: %s NextComments error: %v", meta.inst.User, err)
							meta.addSubCommentFinish = true
							break
						}
					}

					for ; meta.nextCommentIdx < len(meta.curComments.Comments); meta.nextCommentIdx++ {
						if meta.inst.IsSpeedLimit(goinsta.OperNameComment) {
							return laterError
						}
						comment := meta.curComments.Comments[meta.nextCommentIdx]
						err = opt.AddComment(&goinsta.AddCommentParams{
							ParentCommentId:  fmt.Sprintf("%d", comment.Pk),
							UserName:         comment.User.Username,
							LoggingInfoToken: media.Media.LoggingInfoToken,
							MediaId:          media.Media.Id,
							CommentText:      getComments(),
						})
						if err != nil {
							log.Error("account: %s AddComment for sub error: %v", meta.inst.User, err)
						}

						if config.Develop.NeedLike {
							if meta.followSet[comment.User.Pk] == false {
								if meta.inst.IsSpeedLimit(goinsta.OperNameLikeUser) {
									return laterError
								}
								meta.followSet[comment.User.Pk] = true
								err = opt.LikeUser(comment.User.Pk)
								if err != nil {
									return err
								}
							}
						}

						meta.subCommentCount++
						if meta.subCommentCount > config.Develop.MaxSubComment {
							goto finishSubComment
						}

					}

					meta.curComments = nil
					meta.nextCommentIdx = 0
				}

			finishSubComment:
				meta.comments = nil
				meta.addSubCommentFinish = true
			}

			if meta.inst.IsSpeedLimit(goinsta.OperNameComment) {
				return laterError
			}
			err = opt.AddComment(&goinsta.AddCommentParams{
				ParentCommentId:  "",
				UserName:         "",
				LoggingInfoToken: media.Media.LoggingInfoToken,
				MediaId:          media.Media.Id,
				CommentText:      getComments(),
			})
			if err != nil {
				log.Error("account: %s AddComment error: %v", meta.inst.User, err)
			}

			if config.Develop.NeedLike {
				if meta.inst.IsSpeedLimit(goinsta.OperNameLikeUser) {
					return laterError
				}
				if !media.Media.User.IsFavorite {
					err = opt.LikeUser(media.Media.User.Pk)
					if err != nil {
						return err
					}
				}
			}

			meta.hadCheckMedia = false
			meta.subCommentCount = 0
			meta.hadShareMedia = false
			meta.addSubCommentFinish = false
			meta.curComments = nil
			meta.nextCommentIdx = 0
			meta.comments = nil
		}
	}

	return nil
}

func getMeta() *DevelopMeta {
	var retMeta *DevelopMeta
	for true {
		metaListLock.Lock()
		if metaList.Len() == 0 {
			log.Warn("no more meta in list!")
		}

		for item := metaList.Front(); item != nil; item = item.Next() {
			meta := item.Value.(*DevelopMeta)
			if meta.isRunning {
				continue
			}

			if !meta.lastFeedBackTime.IsZero() &&
				time.Since(meta.lastFeedBackTime) < time.Second*time.Duration(config.Develop.FeedBackSleep) {
				continue
			}
			if meta.inst.IsSpeedLimit(goinsta.OperNameCrawComment) ||
				meta.inst.IsSpeedLimit(goinsta.OperNameComment) ||
				meta.inst.IsSpeedLimit(goinsta.OperNameLikeUser) {
				continue
			}
			retMeta = meta
			retMeta.isRunning = true
			metaList.Remove(item)
			break
		}
		metaListLock.Unlock()

		if retMeta != nil {
			return retMeta
		}
		log.Warn("all account can not run")
		//time.Sleep(5 * time.Second)
		time.Sleep(5 * time.Minute)
	}

	return nil
}

func putMeta(meta *DevelopMeta) {
	meta.isRunning = false

	goinsta.SaveInstToDB(meta.inst)
	metaListLock.Lock()
	metaList.PushBack(meta)
	metaListLock.Unlock()
}

func pushBadInst(inst *goinsta.Instagram) {
	goinsta.SaveInstToDB(inst)

	metaListLock.Lock()
	defer metaListLock.Unlock()

	badInstList.PushBack(inst)
}

func developIns() {
	defer waitAll.Done()
	for true {
		meta := getMeta()

		err := doDevelopMeta(meta)
		hadPush := false
		if err != nil {
			if err.Error() == goinsta.InsAccountError_ChallengeRequired ||
				err.Error() == goinsta.InsAccountError_LoginRequired {
				log.Error("account: %s error: %v", meta.inst.User, err)
				//pushBadInst(meta.inst)
				//hadPush = true
				meta.lastFeedBackTime = time.Now()
			} else if err.Error() == goinsta.InsAccountError_Feedback {
				log.Warn("account: %s feedback_required error: %v", meta.inst.User, err)
				meta.lastFeedBackTime = time.Now()
			}
		}

		if !hadPush {
			putMeta(meta)
		}
	}
}

func developServer() {
	insts := goinsta.LoadAccountByTags([]string{config.AccountTag})
	if len(insts) == 0 {
		log.Warn("there have no account!")
		return
	}
	log.Info("load account count %d", len(insts))
	badInstList = list.New()
	metaList = list.New()
	for _, inst := range insts {
		metaList.PushBack(
			&DevelopMeta{
				inst:      inst,
				feed:      inst.GetVideoFeed(),
				followSet: map[int64]bool{},
			})
	}

	waitAll.Add(config.Coro + 1)
	for i := 0; i < config.Coro; i++ {
		go developIns()
	}
	waitAll.Wait()

	log.Info("finish task!")
}

func DevelopAccount() {
	developServer()
	//insts := goinsta.LoadAccountByTags([]string{"dev8"})
	//inst := insts[0]
	//feed := inst.GetVideoFeed()
	//inst.GetAccount().Sync()
	//
	//inst.SetProxy(proxys.ProxyPool.Get(inst.AccountInfo.Register.RegisterIpCountry, ""))
	//for true {
	//	log.Info("%s", inst.Proxy.Rip)
	//	curVideoList, err := feed.Next()
	//	if err != nil {
	//		feed.MoreAvailable = true
	//		feed.SessionId = fmt.Sprintf("%d_%s", inst.ID, strings.ToUpper(common.GenUUID()))
	//		log.Error("account: %s feedVideo.Next error: %v", inst.User, err)
	//	}
	//	print(curVideoList)
	//	if len(curVideoList.Items) == 0 {
	//		log.Error("account: %s not feedVideo any", inst.User)
	//	}
	//	inst.SetProxy(proxys.ProxyPool.Get(inst.AccountInfo.Register.RegisterIpCountry, ""))
	//	time.Sleep(5 * time.Second)
	//}
	//
	//dev := &DevelopMeta{
	//	inst:      insts[0],
	//	feed:      insts[0].GetVideoFeed(),
	//	followSet: map[int64]bool{},
	//}
	//	config.Develop = &config.Develop{
	//		MaxSubComment: 5,
	//		Comments: []string{"I'm an overseas student and I need some help. Can you support me live on the this website %s  followmebsix.com",
	//			"%s Can you come to my studio and watch me perform followmebsix.com",
	//			"My ⊙⊙ are big. %s Want to see? at this followmebsix.com",
	//		},
	//	}
	//	for true {
	//		log.Info("---------------- retry -------------------")
	//		doDevelopMeta(dev)
	//		goinsta.SaveInstToDB(insts[0])
	//		log.Info("sleep...")
	//		time.Sleep(60 * time.Second)
	//	}
	//CronTask := cron.New(cron.WithSeconds())
	//_, err := CronTask.AddFunc(config.Develop.Spec, developServer)
	//if err != nil {
	//	log.Error("pusher run time task error: %v", err)
	//	return
	//}
	//CronTask.Start()
	//select {}

	//inst := insts[0]
	////inst.GetAccount().SetGender()
	////inst.GetAccount().SetBiography("you can look me ⊙⊙ at followmebsix.com")
	////inst.GetAccount().UpdatePronouns("her")
	//inst.GetAccount().Sync()
	//inst.GetAccount().EditProfile(&goinsta.UserProfile{
	//	ExternalUrl: "http://followmebsix.com",
	//})
	//
	////print(next)
}

//func doDevelop(inst *goinsta.Instagram) (retErr error) {
//	defer func() {
//		if err := recover(); err != nil {
//			log.Error("account: %s error: %v", inst.User, err)
//			retErr = err.(error)
//			if retErr.Error() == goinsta.InsAccountError_Feedback {
//				retErr = laterError
//				log.Error("account: %s feedback", inst.User)
//			}
//		}
//	}()
//
//	opt := inst.GetUserOperate()
//
//	if !inst.IsSpeedLimit(goinsta.OperNamePostImg) {
//
//	}
//	if !inst.IsSpeedLimit(goinsta.OperNamePostImg) {
//
//	}
//
//	feed := inst.GetVideoFeed()
//	for true {
//		if inst.IsSpeedLimit(goinsta.OperNameFeedVideo) {
//			return laterError
//		}
//		next, err := feed.Next()
//		if err != nil {
//			log.Error("account: %s feedVideo.Next error: %v", inst.User, err)
//			return err
//		}
//		if len(next.Items) == 0 {
//			log.Error("account: %s not feedVideo any", inst.User)
//			return laterError
//		}
//
//		for _, media := range next.Items {
//			if media.Media.CommentingDisabledForViewer {
//				continue
//			}
//
//			if media.Media.HasMoreComments {
//				comments := inst.NewComments(media.Media.Id)
//				for true {
//					if inst.IsSpeedLimit(goinsta.OperNameCrawComment) {
//						return laterError
//					}
//
//					var nextComments *goinsta.RespComments
//					nextComments, err = comments.NextComments()
//					if err != nil {
//						log.Error("account: %s NextComments error: %v", inst.User, err)
//						break
//					}
//
//					for _, cm := range nextComments.GetAllComments() {
//						if inst.IsSpeedLimit(goinsta.OperNameComment) {
//							return laterError
//						}
//
//						err = opt.AddComment(&goinsta.AddCommentParams{
//							ParentCommentId:  fmt.Sprintf("%d", cm.Pk),
//							UserName:         cm.User.Username,
//							LoggingInfoToken: media.Media.LoggingInfoToken,
//							MediaId:          media.Media.Id,
//							CommentText:      "i like u comment," + common.GenString(common.CharSet_123, 5) + " and more sexy img in followmebsix.com",
//						})
//						if err != nil {
//							log.Error("account: %s AddComment for sub error: %v", inst.User, err)
//							break
//						}
//					}
//				}
//			}
//
//			var shareMedia string
//			shareMedia, err = opt.ShareMedia(media.Media.Id)
//			if err != nil {
//				log.Error("account: %s ShareMedia error: %v", inst.User, err)
//			} else {
//				LogMedia(inst, shareMedia, media.Media)
//			}
//
//			err = opt.AddComment(&goinsta.AddCommentParams{
//				ParentCommentId:  "",
//				UserName:         "",
//				LoggingInfoToken: media.Media.LoggingInfoToken,
//				MediaId:          media.Media.Id,
//				CommentText:      "i like u video," + common.GenString(common.CharSet_123, 5) + " and more sexy img in followmebsix.com",
//			})
//			if err != nil {
//				log.Error("account: %s AddComment error: %v", inst.User, err)
//			}
//		}
//	}
//	return nil
//}
