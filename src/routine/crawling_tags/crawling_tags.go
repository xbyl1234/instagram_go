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
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type CrawConfig struct {
	TaskName              string `json:"task_name"`
	SearchTag             string `json:"search_tag"`
	StartTime             string `json:"start_time"`
	StopTime              string `json:"last_time"`
	MediaCoroCount        int    `json:"media_coro_count"`
	CommonCoroCount       int    `json:"common_coro_count"`
	ProxyPath             string `json:"proxy_path"`
	WorkPath              string `json:"config_path"`
	CrawTagsStatus        bool   `json:"craw_tags_status"`
	CrawMediasStatus      bool   `json:"craw_medias_status"`
	CrawCommentUserStatus bool   `json:"craw_comment_user_status"`
}

var config CrawConfig
var WorkPath string
var PathSeparator = string(os.PathSeparator)
var StopTime time.Time
var WaitAll sync.WaitGroup

var TagList = list.New()
var TagIDSet = mapset.NewSet()
var MediaChan = make(chan *routine.MediaComb, 1000)
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
	var search *goinsta.Search
	var err error
	inst := routine.ReqAccount()
	if inst == nil {
		log.Error("CrawTags no more account!")
		return
	}

	defer goinsta.AccountPool.ReleaseOne(inst)

	search, err = routine.LoadSearch()
	if err != nil {
		log.Error("preCrawTags LoadTags error:%v", err)
		os.Exit(0)
	}
	if search != nil {
		search.SetAccount(inst)
	} else {
		search = inst.GetSearch(config.SearchTag)
		_ = routine.SaveSearch(search)
	}
	if !search.HasMore {
		log.Info("pass search tag...")
		config.CrawTagsStatus = true
		return
	}

	for true {
		searchResult, err := search.NextTags()
		if err != nil {
			if common.IsNoMoreError(err) {
				config.CrawTagsStatus = true
				log.Info("tags has craw finish!")
				break
			} else if inst.NeedReplace() || common.IsError(err, common.RequestError) {
				if inst.NeedReplace() {
					goinsta.AccountPool.BlackOne(inst)
					_inst := routine.ReqAccount()
					if _inst == nil {
						log.Error("CrawTags no more account!")
						break
					}
					log.Warn("CrawTags replace account %s->%s", inst.User, _inst.User)
					inst = _inst
					search.SetAccount(_inst)
				} else {
					log.Warn("CrawMedias retrying...user: %s, err: %v", inst.User, err)
				}
				continue
			} else {
				log.Error("search next error: %v", err)
				break
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

func CrawMedias() {
	defer WaitAll.Done()
	inst := routine.ReqAccount()
	if inst == nil {
		log.Error("CrawMedias no more account!")
		return
	}

	for tag := range TagsChan {
		tag.SetAccount(inst)
		err := tag.Sync(goinsta.TabRecent)
		if err != nil {
			log.Error("tag sync error: %v", err)
		}
		_, err = tag.Stories()
		if err != nil {
			log.Error("tag stories error: %v", err)
		}

		for true {
			tagResult, err := tag.Next()
			if err != nil {
				if common.IsNoMoreError(err) {
					num := atomic.AddInt32(&NotFinishTags, -1)
					if num == 0 {
						config.CrawMediasStatus = true
					}
					log.Info("tags %s medias has craw finish!", tag.Name)
					break
				} else if inst.NeedReplace() || common.IsError(err, common.RequestError) {
					if inst.NeedReplace() {
						goinsta.AccountPool.BlackOne(inst)
						_inst := routine.ReqAccount()
						if _inst == nil {
							log.Error("CrawMedias no more account!")
							break
						}
						log.Warn("CrawMedias replace account %s->%s", inst.User, _inst.User)
						inst = _inst
						tag.SetAccount(_inst)
					} else {
						log.Warn("CrawMedias retrying...user: %s, err: %v", inst.User, err)
					}
					continue
				} else {
					log.Error("next media error: %v", err)
					break
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
	goinsta.AccountPool.ReleaseOne(inst)
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

func CrawCommonUser() {
	defer WaitAll.Done()
	inst := routine.ReqAccount()
	if inst == nil {
		log.Error("CrawCommonUser no more account!")
		return
	}

	unknowErrorCount := 0

	for mediaComb := range MediaChan {
		if mediaComb.Media.CommentCount == 0 {
			continue
		}
		mediaComb.Media.SetAccount(inst)
		if mediaComb.Comments != nil {
			mediaComb.Comments.SetAccount(inst)
		} else {
			mediaComb.Comments = mediaComb.Media.GetComments()
		}

		for true {
			respComm, err := mediaComb.Comments.NextComments()
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("media %s comments has craw finish!", mediaComb.Media.ID)
					break
				} else if inst.NeedReplace() || common.IsError(err, common.RequestError) {
					if inst.NeedReplace() {
						goinsta.AccountPool.BlackOne(inst)
						_inst := routine.ReqAccount()
						if _inst == nil {
							log.Error("CrawCommonUser no more account!")
							break
						}
						log.Warn("CrawCommonUser replace account %s->%s", inst.User, _inst.User)
						inst = _inst
						mediaComb.Media.SetAccount(_inst)
						mediaComb.Comments.SetAccount(_inst)
					} else {
						log.Warn("CrawCommonUser retrying...user: %s, err: %v", inst.User, err)
					}
					continue
				} else {
					unknowErrorCount++
					log.Error("NextComments error:%v", err)
					if strings.Index(err.Error(), "Media is unavailable") >= 0 {
						mediaComb.Comments.HasMore = false
						routine.SaveComments(mediaComb)
						break
					}
					if unknowErrorCount > 3 {
						return
					} else if unknowErrorCount != 0 {
						break
					}
				}
			} else {
				unknowErrorCount = 0
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

	goinsta.AccountPool.ReleaseOne(inst)
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

	if config.MediaCoroCount == 0 {
		config.MediaCoroCount = runtime.NumCPU()
	}
	if config.CommonCoroCount == 0 {
		config.CommonCoroCount = runtime.NumCPU() * 2
	}

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
			//go CrawMedias()
		}
	}

	WaitAll.Add(1 + config.CommonCoroCount)
	go SendMedias()
	for index := 0; index < config.CommonCoroCount; index++ {
		go CrawCommonUser()
	}

	WaitAll.Wait()
}
