package main

import (
	"container/list"
	"flag"
	mapset "github.com/deckarep/golang-set"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"os"
	"runtime"
	"time"
)

type CrawConfig struct {
	TaskName  string    `json:"task_name"`
	SearchTag string    `json:"search_tag"`
	LastTime  time.Time `json:"last_time"`
	CoroCount int       `json:"coro_count"`
	PorxyPath string    `json:"porxy_path"`
}

var config CrawConfig
var WorkPath string
var PathSeparator = string(os.PathSeparator)
var IsFirstRun bool

func ReqAccount() *goinsta.Instagram {
	inst := goinsta.AccountPool.GetOne()
	_proxy := common.ProxyPool.Get(inst.Proxy.ID)
	if _proxy == nil {
		log.Error("find insta proxy error!")
		os.Exit(0)
	}
	inst.SetProxy(_proxy)
	return inst
}

var TagList = list.New()
var TagIDSet = mapset.NewSet()

func preCrawTags() {
	var err error
	TagList, err = LoadTags()
	if err != nil {
		log.Error("preCrawTags LoadTags error:%v", err)
		os.Exit(0)
	}
	for item := TagList.Front(); item != nil; item = item.Next() {
		TagIDSet.Add(item.Value.(*goinsta.Tags).Id)
	}
}

func doCrawTags(search *goinsta.Search) {
	for true {
		searchResult, err := search.NextTags()
		if err != nil {
			if err.Error() == common.MakeMoneyError_NoMore.Error() {
				log.Info("tags has craw finish!")
			} else {
				log.Error("search next error: %v", err)
			}
			break
		}

		log.Info("%v", searchResult)
		tags := searchResult.GetTags()
		for index := range tags {
			if TagIDSet.Contains(tags[index].Id) {
				log.Info("")
				continue
			}
			TagIDSet.Add(tags[index].Id)
			TagList.PushBack(tags[index])
			err = SaveTags(&tags[index])
			if err != nil {
				log.Error("SaveTags error:%v", err)
			}
		}
		_ = SaveSearch(search)
	}
}

func CrawTags() {
	var search *goinsta.Search
	var err error
	inst := ReqAccount()

	if !IsFirstRun {
		search, err = LoadSearch()
		if err != nil {
			log.Error("preCrawTags LoadTags error:%v", err)
			os.Exit(0)
		}
		search.Inst = inst
		preCrawTags()
	} else {
		search = inst.GetSearch("china")
		_ = SaveSearch(search)
	}
	doCrawTags(search)
	goinsta.AccountPool.ReleaseOne(inst)
}

func CrawMedias(tag *goinsta.Tags) {
	inst := ReqAccount()
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
			if err.Error() == common.MakeMoneyError_NoMore.Error() {
				log.Info("tags has craw finish!")
			} else {
				log.Error("next media error: %v", err)
			}
			break
		}
		medias := tagResult.GetAllMedias()
		for index := range medias {
			err = SaveMedia(medias[index], nil)
			if err != nil {
				log.Error("SaveMedia error:%v", err)
			}
		}
		err = SaveTags(tag)
		if err != nil {
			log.Error("SaveTags error:%v", err)
		}
	}

	goinsta.AccountPool.ReleaseOne(inst)
}

func CrawCommonUser(mediaComb *MediaComb) {
	if mediaComb.Media.CommentCount == 0 {
		return
	}
	inst := ReqAccount()
	mediaComb.Media.SetAccount(inst)

	if mediaComb.Comments != nil {
		mediaComb.Comments.SetAccount(inst)
	} else {
		mediaComb.Comments = mediaComb.Media.GetComments()
	}

	for true {
		respComm, err := mediaComb.Comments.NextComments()
		if err != nil {

		}

		comments := respComm.GetAllComments()
		for index := range comments {

		}
	}

}

func do() {
	for true {

		for index := range tags {
			tag := tags[index]
			tag.Sync(goinsta.TabRecent)
			tag.Stories()
			tagResult, err := tag.Next()
			if err != nil {
				log.Error("next media error: %v", err)
				break
			}
			medias := tagResult.GetAllMedias()
			for mindex := range medias {
				media := medias[mindex]
				comments := media.GetComments()
				commResp, err := comments.NextComments()
				if err != nil {
					log.Error("next comm error: %v", err)
					break
				}
				for cindex := range commResp.Comments {
					log.Info("comment user id: %v", commResp.Comments[cindex].User.ID)
				}
			}
		}
	}
}

func initParams() {
	var err error
	var TaskName = flag.String("tn", "", "task name")
	var SearchTag = flag.String("tag", "", "search tag")
	var LastTime = flag.String("lt", "2021-05-01", "stop last time")
	var ConfigPath = flag.String("cp", "", "config path")
	var CoroCount = flag.Int("coro", 0, "coro count")
	var PorxyPath = flag.String("pp", 0, "proxy path")

	flag.Parse()
	if *ConfigPath != "" {
		return
	}
	if *TaskName == "" {
		*TaskName = "taks_" + time.Now().String()
		log.Info("task name is %s", *TaskName)
	}
	WorkPath, _ = os.Getwd()
	err = os.MkdirAll(WorkPath+PathSeparator, 777)
	if err != nil {
		log.Error("make dir: %s error: %v", WorkPath+PathSeparator, err)
		os.Exit(0)
	}
	//log.InitLogger()

	if *PorxyPath == "" {
		log.Error("PorxyPath is null")
		os.Exit(0)
	}

	if *SearchTag == "" {
		log.Error("SearchTag is null")
		os.Exit(0)
	}
	if *LastTime != "" {
		config.LastTime, err = time.Parse("2006-01-02", *LastTime)
		if err != nil {
			log.Error("parse LastTime: %s, error: %v", *LastTime, err)
			os.Exit(0)
		}
	} else {
		config.LastTime = time.Now().Add(0 - time.Hour*24*30)
		log.Info("last time is last month! time:%v", config.LastTime)
	}
	if *CoroCount == 0 {
		config.CoroCount = runtime.NumCPU()*2 + 1
	}

	config.TaskName = *TaskName
	config.SearchTag = *SearchTag
	config.PorxyPath = *PorxyPath
	log.Info("init config success!")
}

func main() {
	initParams()

	InitTest()
	do()
}
