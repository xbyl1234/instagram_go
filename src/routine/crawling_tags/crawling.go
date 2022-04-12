package main

import (
	"flag"
	"github.com/robfig/cron/v3"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"time"
)

type CrawConfig struct {
	AddCommentCount  int    `json:"add_comment_count"`
	AccountTag       string `json:"account_tag"`
	TaskName         string `json:"task_name"`
	CommentCoroCount int    `json:"comment_coro_count"`
	ProxyPath        string `json:"proxy_path"`
	KeyWordPath      string `json:"key_word_path"`
	CrawMediasFreq   string `json:"craw_medias_freq"`
}

var config CrawConfig
var CronTask *cron.Cron
var timedTaskerID = 0
var mediaRedis *common.Queue

func TimedTasker() {
	now := time.Now()
	year, month, day := now.Date()
	scanEnd := time.Date(year, month, day-1, 0, 0, 0, 0, time.Local)

	log.Info("task id %d: will running, this time is %s, scan end time is %s", timedTaskerID, now.String(), scanEnd.String())

	tags, err := LoadKeyWord()
	if err != nil {
		log.Error("%v", err)
		return
	}

	var waitCrawMedia sync.WaitGroup
	var waitCrawComment sync.WaitGroup
	waitCrawMedia.Add(len(tags))
	waitCrawComment.Add(config.CommentCoroCount)

	mediaChan := make(chan *MediaComb, 10)

	for _, item := range tags {
		go CrawMedias(item, mediaChan, &waitCrawMedia, scanEnd)
	}

	for index := 0; index < config.CommentCoroCount; index++ {
		go CrawCommentUser(mediaChan, &waitCrawComment)
	}

	waitCrawMedia.Wait()
	waitCrawComment.Wait()
	log.Info("task id %d: finish! this time is %s, scan end time is %s", timedTaskerID, now.String(), scanEnd.String())
}

func initParams() {
	var err error
	var TaskConfigPath = flag.String("config", "./config/craw.json", "")
	log.InitDefaultLog("craw", true, true)
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

	log.Info("init config success!")
}

func main() {
	common.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitCrawTagsDB(config.TaskName)

	intas := goinsta.LoadAccountByTags([]string{config.AccountTag})
	if len(intas) == 0 {
		log.Warn("there have no account!")
	} else {
		goinsta.InitAccountPool(intas)
	}

	//CronTask = cron.New(cron.WithSeconds())
	//_, _ = CronTask.AddFunc(config.CrawMediasFreq, TimedTasker)
	//log.Info("start timer task")
	//CronTask.Start()
	mediaRedis, _ = common.CreateQueue(&routine.DBConfig.Redis)

	go TimedTasker()

	select {}
}
