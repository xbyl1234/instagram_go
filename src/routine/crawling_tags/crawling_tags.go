package main

import (
	"container/list"
	"flag"
	mapset "github.com/deckarep/golang-set"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type CrawConfig struct {
	AccountPoolTags            string `json:"account_pool_tags"`
	TaskName                   string `json:"task_name"`
	SearchTag                  string `json:"search_tag"`
	StartTime                  string `json:"start_time"`
	StopTime                   string `json:"last_time"`
	MediaCoroCount             int    `json:"media_coro_count"`
	CommentCoroCount           int    `json:"comment_coro_count"`
	ProxyPath                  string `json:"proxy_path"`
	WorkPath                   string `json:"config_path"`
	CrawTagsStatus             bool   `json:"craw_tags_status"`
	CrawMediasStatus           bool   `json:"craw_medias_status"`
	CrawCommentUserStatus      bool   `json:"craw_comment_user_status"`
	CrawMediaMaxRequestCount   int    `json:"craw_media_max_request_count"`
	CrawCommentMaxRequestCount int    `json:"craw_comment_max_request_count"`
}

var config CrawConfig
var WorkPath string
var PathSeparator = string(os.PathSeparator)
var StopTime time.Time
var WaitAll sync.WaitGroup

var TagList = list.New()
var TagIDSet = mapset.NewSet()
var MediaChan = make(chan *routine.MediaComb, 20)
var NotFinishTags int32
var TagsChan = make(chan *goinsta.Tags, 10)

func LoadTags() {
	retTagList, err := routine.LoadTags()
	if err != nil {
		log.Error("preCrawTags LoadTags error:%v", err)
		os.Exit(0)
	}

	for index := range retTagList {
		TagIDSet.Add(retTagList[index].Id)
		TagList.PushBack(&retTagList[index])
	}
}

func CrawTags() {
	var currAccount *goinsta.Instagram
	var search *goinsta.Search
	var err error

	var RequireAccount = func(search *goinsta.Search) *goinsta.Search {
		inst, err := routine.ReqAccount(config.AccountPoolTags, true)
		if err != nil {
			log.Error("CrawTags req account error: %v!", err)
			return nil
		}
		if search == nil {
			search = inst.GetSearch(config.SearchTag)
			_ = routine.SaveSearch(search)
			log.Info("CrawTags set account %s", inst.User)
		} else {
			if search.Inst != nil {
				oldUser := search.Inst.User
				goinsta.AccountPool.ReleaseOne(search.Inst)
				search.SetAccount(inst)
				log.Warn("CrawTags replace account %s->%s", oldUser, inst.User)
			} else {
				search.SetAccount(inst)
				log.Info("CrawTags set account %s", inst.User)
			}
		}
		return search
	}

	search, err = routine.LoadSearch()
	if err != nil {
		log.Error("preCrawTags LoadTags error:%v", err)
		os.Exit(0)
	}
	if search != nil && !search.HasMore {
		log.Info("pass search tag...")
		config.CrawTagsStatus = true
		return
	}

	search = RequireAccount(search)
	defer func() {
		if currAccount != nil {
			goinsta.AccountPool.ReleaseOne(currAccount)
		}
	}()

	for true {
		searchResult, err := search.NextTags()
		if err != nil {
			if common.IsNoMoreError(err) {
				config.CrawTagsStatus = true
				log.Info("tags has craw finish!")
				break
			} else if common.IsError(err, common.ChallengeRequiredError) || common.IsError(err, common.LoginRequiredError) {
				search = RequireAccount(search)
				continue
			} else if common.IsError(err, common.RequestError) {
				log.Warn("CrawMedias retrying...user: %s, err: %v", search.Inst.User, err)
				continue
			} else {
				log.Error("search next unknow error: %v", err)
				continue
			}
		}

		log.Info("%v", searchResult)
		tags := searchResult.GetTags()
		for index := range tags {
			if TagIDSet.Contains(tags[index].Id) {
				log.Info("")
				continue
			}
			TagIDSet.Add(tags[index].Id)
			TagList.PushBack(&tags[index])
			err = routine.SaveTags(&tags[index])
			if err != nil {
				log.Error("SaveTags error:%v", err)
			}
		}
		_ = routine.SaveSearch(search)
	}
}

//830
//847
func CrawMedias() {
	defer WaitAll.Done()
	var currAccount *goinsta.Instagram
	var SetNewAccount = func(tag *goinsta.Tags) {
		inst, err := routine.ReqAccount(config.AccountPoolTags, true)
		if err != nil {
			log.Error("CrawMedias req account error: %v", err)
			return
		}
		tag.SetAccount(inst)
		currAccount = inst
		err = tag.Sync(goinsta.TabRecent)
		if err != nil {
			log.Error("tag sync error: %v", err)
		}
		_, err = tag.Stories()
		if err != nil {
			log.Error("tag stories error: %v", err)
		}
	}

	var RequireAccount = func(tag *goinsta.Tags, reqCount int) int {
		var oldUser string
		if tag.Inst == nil {
			SetNewAccount(tag)
			log.Info("CrawMedias set account %s", tag.Inst.User)
			return 0
		} else {
			oldUser = tag.Inst.User
			if reqCount > config.CrawMediaMaxRequestCount || tag.Inst.IsBad() {
				goinsta.AccountPool.ReleaseOne(tag.Inst)
				SetNewAccount(tag)
				log.Warn("CrawMedias replace account %s->%s", oldUser, tag.Inst.User)
				return 0
			} else {
				return reqCount
			}
		}
	}
	defer func() {
		if currAccount != nil {
			goinsta.AccountPool.ReleaseOne(currAccount)
		}
	}()

	var reqCount = 0
	for tag := range TagsChan {
		reqCount = RequireAccount(tag, reqCount)
		for true {
			reqCount = RequireAccount(tag, reqCount)
			reqCount++
			tagResult, err := tag.Next()
			log.Info("account %s req count %d status %s", currAccount.User, reqCount, currAccount.Status)
			if err != nil {
				if common.IsNoMoreError(err) {
					num := atomic.AddInt32(&NotFinishTags, -1)
					if num == 0 {
						config.CrawMediasStatus = true
					}
					log.Info("tags %s medias has craw finish!", tag.Name)
					break
				} else if common.IsError(err, common.ChallengeRequiredError) ||
					common.IsError(err, common.FeedbackError) ||
					common.IsError(err, common.LoginRequiredError) {
					log.Error("user %s status is %s from CrawCommentUser task, err: %v", currAccount.User,
						currAccount.User, err)
					reqCount = RequireAccount(tag, reqCount)
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
					if mediaTime.Sub(StopTime) < 0 {
						stop = true
						tag.MoreAvailable = false
						log.Info("craw media stop! current time is %s", mediaTime.Format("2006-01-02 15:04:05"))
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

			err = routine.SaveTags(tag)
			if err != nil {
				log.Error("SaveTags error:%v", err)
			}
			if stop {
				break
			}
		}
	}
}

func SendTags() {
	defer WaitAll.Done()
	for true {
		for item := TagList.Front(); item != nil; item = item.Next() {
			tags := item.Value.(*goinsta.Tags)
			if tags.MoreAvailable {
				TagsChan <- tags
			}
		}
	}

	close(TagsChan)
}

//816
//1320
func CrawCommentUser() {
	defer WaitAll.Done()
	var currAccount *goinsta.Instagram
	reqCount := 0

	var SetNewAccount = func(mediaComb *routine.MediaComb, inst *goinsta.Instagram) {
		mediaComb.Media.SetAccount(inst)
		if mediaComb.Comments == nil {
			mediaComb.Comments = mediaComb.Media.GetComments()
		} else {
			mediaComb.Comments.SetAccount(inst)
		}
	}

	var RequireAccont = func(mediaComb *routine.MediaComb, reqCount int) int {
		if reqCount > config.CrawCommentMaxRequestCount || currAccount == nil || currAccount.IsBad() {
			var oldUser string
			if currAccount != nil {
				oldUser = currAccount.User
				goinsta.AccountPool.ReleaseOne(currAccount)
			}

			inst, err := routine.ReqAccount(config.AccountPoolTags, true)
			if err != nil {
				log.Error("CrawCommentUser req account error: %v!", err)
			}
			currAccount = inst

			SetNewAccount(mediaComb, currAccount)
			log.Warn("CrawCommentUser replace account to %s->%s", oldUser, currAccount.User)
			return 0
		} else {
			SetNewAccount(mediaComb, currAccount)
			return reqCount
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
			routine.SaveComments(mediaComb)
			continue
		}
		for true {
			reqCount = RequireAccont(mediaComb, reqCount)
			reqCount++
			respComm, err := mediaComb.Comments.NextComments()
			log.Info("account %s req count %d status %s", currAccount.User, reqCount, currAccount.Status)
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("media %s comments has craw finish!", mediaComb.Media.ID)
					break
				} else if common.IsError(err, common.ChallengeRequiredError) ||
					common.IsError(err, common.FeedbackError) ||
					common.IsError(err, common.LoginRequiredError) {
					log.Error("user %s status is %s from CrawCommentUser task, err: %v", currAccount.User,
						currAccount.Status, err)
					reqCount = RequireAccont(mediaComb, reqCount)
					continue
				} else if common.IsError(err, common.RequestError) {
					log.Warn("CrawCommentUser retrying...user: %s, err: %v", currAccount.User, err)
					continue
				} else if strings.Index(err.Error(), "Media is unavailable") >= 0 {
					log.Warn("Media %d is unavailable", mediaComb.Media.ID)
					mediaComb.Comments.HasMore = false
					routine.SaveComments(mediaComb)
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
			err = routine.SaveComments(mediaComb)
			if err != nil {
				log.Error("SaveMedia error:%v", err)
				break
			}
		}
	}
}

func SendMedias() {
	defer WaitAll.Done()
	for true {
		medias, err := routine.LoadMedia(100)
		if err != nil {
			log.Error("load media error: %v", err)
			break
		}
		if len(medias) == 0 && config.CrawTagsStatus && config.CrawMediasStatus {
			config.CrawCommentUserStatus = true
			log.Info("craw common user finish!")
			break
		}

		for index := range medias {
			MediaChan <- &medias[index]
		}
	}

	close(MediaChan)
}

func initParams() {
	var err error
	var TaskConfigPath = flag.String("task", "", "task")
	//-tn test_craw  -tag game -coro 1 -pp C:\Users\Administrator\Desktop\project\github\instagram_project\data\zone2_ips_us.txt
	log.InitDefaultLog("craw_user", true, true)
	flag.Parse()
	if *TaskConfigPath == "" {
		log.Error("task config path is null!")
		os.Exit(0)
	}

	err = common.LoadJsonFile(*TaskConfigPath, &config)
	if err != nil {
		log.Error("load task config error: %v", err)
		os.Exit(0)
	}

	if config.TaskName == "" {
		log.Error("task name is null")
		os.Exit(0)
	}

	if config.StopTime != "" {
		StopTime, err = time.Parse("2006-01-02", config.StopTime)
		if err != nil {
			log.Error("parse StopTime: %s, error: %v", config.StopTime, err)
			os.Exit(0)
		}
	} else {
		StopTime = time.Now().Add(0 - time.Hour*24*30*12)
		config.StopTime = StopTime.Format("2006-01-02")
		log.Info("stop time is last year! time:%v", config.StopTime)
	}
	if config.StartTime == "" {
		config.StartTime = time.Now().Format("2006-01-02")
	}
	if config.AccountPoolTags == "" {
		log.Error("parse AccountPoolTags is null")
		os.Exit(0)
	}
	//if config.MediaCoroCount == 0 {
	//	config.MediaCoroCount = runtime.NumCPU()
	//}
	//if config.CommentCoroCount == 0 {
	//	config.CommentCoroCount = runtime.NumCPU() * 2
	//}

	WorkPath, _ = os.Getwd()
	if config.WorkPath == "" {
		config.WorkPath = WorkPath + PathSeparator + config.TaskName + PathSeparator
	}

	err = os.MkdirAll(config.WorkPath, 777)
	if err != nil {
		log.Error("make dir: %s error: %v", config.WorkPath, err)
		os.Exit(0)
	}

	err = common.Dumps(*TaskConfigPath, &config)
	if err != nil {
		log.Error("Dumps config error: %v", err)
		os.Exit(0)
	}
	log.Info("init config success!")
}

func main() {
	config2.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitCrawTagsDB(config.TaskName)

	LoadTags()
	CrawTags()
	log.Info("tags count: %d", TagList.Len())

	if TagList.Len() == 0 {
		config.CrawMediasStatus = true
		log.Info("pass craw medias...")
	} else {
		NotFinishTags = int32(TagList.Len())
		WaitAll.Add(config.MediaCoroCount + 1)
		go SendTags()
		for index := 0; index < config.MediaCoroCount; index++ {
			go CrawMedias()
		}
	}

	WaitAll.Add(1 + config.CommentCoroCount)
	go SendMedias()
	for index := 0; index < config.CommentCoroCount; index++ {
		go CrawCommentUser()
	}

	WaitAll.Wait()
}
