package main

import (
	"flag"
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

func do() {
	inst := goinsta.AccountPool.GetOne()
	_proxy := common.ProxyPool.Get(inst.Proxy.ID)
	if _proxy == nil {
		log.Error("find insta proxy error!")
		os.Exit(0)
	}
	inst.SetProxy(_proxy)

	search := inst.GetSearch("china")
	for true {
		searchResult, err := search.NextTags()
		if err != nil {
			log.Error("search next error: %v", err)
			break
		}
		log.Info("%v", searchResult)

		tags := searchResult.GetTags()
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
