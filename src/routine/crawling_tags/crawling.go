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
	"sync/atomic"
	"time"
)

type CrawConfig struct {
	CrawMediaAccountTag   string `json:"craw_media_account_tag"`
	CrawCommentAccountTag string `json:"craw_comment_account_tag"`
	TaskName              string `json:"task_name"`
	MediaCoroCount        int    `json:"media_coro_count"`
	CommentCoroCount      int    `json:"comment_coro_count"`
	ProxyPath             string `json:"proxy_path"`
	KeyWordPath           string `json:"key_word_path"`
	CrawMediasFreq        string `json:"craw_medias_freq"`
}

var config CrawConfig
var MediaChan = make(chan *routine.MediaComb, 20)
var CronTask *cron.Cron
var timedTaskerID = 0
var curRunTimedTaskerCount int32 = 0

func TimedTasker() {
	now := time.Now()
	year, month, day := now.Date()
	scanEnd := time.Date(year, month, day-1, 0, 0, 0, 0, time.Local)

	count := atomic.AddInt32(&curRunTimedTaskerCount, 1)
	if count > 1 {
		atomic.AddInt32(&curRunTimedTaskerCount, -1)
		log.Error("task id %d: last task was running, this time is %s, so return!!!", timedTaskerID, now.String())
		return
	}

	log.Info("task id %d: will running, this time is %s, scan end time is %s", timedTaskerID, now.String(), scanEnd.String())
	var TagsChan = make(chan *goinsta.Tags, 10)
	var waitCraw sync.WaitGroup
	waitCraw.Add(config.MediaCoroCount)

	for index := 0; index < config.MediaCoroCount; index++ {
		go CrawMedias(TagsChan, &waitCraw, scanEnd)
	}

	SendTags(TagsChan)
	log.Info("send tags finish!")
	close(TagsChan)
	waitCraw.Wait()
	log.Info("task id %d: finish! this time is %s, scan end time is %s", timedTaskerID, now.String(), scanEnd.String())
	atomic.AddInt32(&curRunTimedTaskerCount, -1)
}

func initParams() {
	var err error
	var TaskConfigPath = flag.String("task", "", "task")
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

	if config.CrawMediaAccountTag == "" || config.CrawCommentAccountTag == "" {
		log.Error("parse Account Tags is null")
		os.Exit(0)
	}
	//if config.MediaCoroCount == 0 {
	//	config.MediaCoroCount = runtime.NumCPU()
	//}
	//if config.CommentCoroCount == 0 {
	//	config.CommentCoroCount = runtime.NumCPU() * 2
	//}
	log.Info("init config success!")
}

func main() {
	common.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitCrawTagsDB(config.TaskName)

	intas := goinsta.LoadAccountByTags([]string{config.CrawMediaAccountTag, config.CrawCommentAccountTag})
	if len(intas) == 0 {
		log.Warn("there have no account!")
	} else {
		goinsta.InitAccountPool(intas)
	}

	LoadTags()
	CheckTagsAndRunCrawTags()

	CronTask = cron.New(cron.WithSeconds())
	_, _ = CronTask.AddFunc(config.CrawMediasFreq, TimedTasker)
	log.Info("start timer task")
	CronTask.Start()
	go TimedTasker()

	go SendMedias()
	for index := 0; index < config.CommentCoroCount; index++ {
		go CrawCommentUser()
	}

	select {}
}
