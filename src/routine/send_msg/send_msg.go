package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type CrawConfig struct {
	TaskName             string   `json:"task_name"`
	CoroCount            int      `json:"coro_count"`
	ProxyPath            string   `json:"proxy_path"`
	TargetUserDB         string   `json:"target_user_db"`
	TargetUserCollection string   `json:"target_user_collection"`
	TargetUserSource     string   `json:"target_user_source"`
	TextMsg              []string `json:"text_msg"`
	ImageMsg             []string `json:"image_msg"`
	ImageMsgUploadID     []string `json:"image_msg_upload_id"`
	IntervalTimeAccount  int      `json:"interval_time_account"`
	IntervalTimeMsg      int      `json:"interval_time_msg"`
}

var config CrawConfig
var ImageData [][]byte
var WaitAll sync.WaitGroup
var UserChan = make(chan *routine.UserComb, 100)

var SendSuccessCount int32
var SendErrorCount int32

func SendMsg(inst *goinsta.Instagram, user *routine.UserComb) error {
	imageIDs := make([]string, len(config.ImageMsg))
	for index := range config.ImageMsg {
		id, err := inst.GetUpload().RuploadIgPhoto(config.ImageMsg[0])
		if err != nil {
			imageIDs[index] = ""
		} else {
			imageIDs[index] = id
		}
	}
}

func SendTask() {
	defer WaitAll.Done()

	for user := range UserChan {
		inst, err := routine.ReqAccount(true)
		if err != nil {
			log.Error("req account unknow error!")
			return
		}

		err = SendMsg(inst, user)
		if err != nil {
			atomic.AddInt32(&SendErrorCount, 1)
		} else {
			atomic.AddInt32(&SendSuccessCount, 1)
		}
		goinsta.AccountPool.CoolingOne(inst)
	}
}

func SendUser() {
	defer WaitAll.Done()

	for true {
		users, err := routine.LoadUser(config.TargetUserSource, config.TaskName, 100)
		if err != nil {
			log.Error("load user error: %v", err)
			break
		}
		if len(users) == 0 {
			log.Info("no more user to load")
			break
		}

		for index := range users {
			UserChan <- &users[index]
		}
	}

	close(UserChan)
}

func initParams() {
	var err error
	var TaskConfigPath = flag.String("task", "", "task")
	log.InitDefaultLog("send_msg", true, true)
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

	if config.CoroCount == 0 {
		config.CoroCount = 1
	}
	if config.TargetUserDB == "" {
		log.Error("TargetUserDB name is null")
		os.Exit(0)
	}
	if config.TargetUserCollection == "" {
		log.Error("TargetUserCollection is null")
		os.Exit(0)
	}
	if config.TargetUserSource == "" {
		log.Error("TargetUserSource is null")
		os.Exit(0)
	}
	if len(config.TextMsg) == 0 {
		log.Error("TextMsg is null")
		os.Exit(0)
	}
	if len(config.ImageMsg) == 0 {
		log.Error("ImageMsg is null")
		os.Exit(0)
	}
	if len(config.ImageMsgUploadID) == 0 {
		id := time.Now().Unix()
		config.ImageMsgUploadID = make([]string, len(config.ImageMsg))
		for index := range config.ImageMsgUploadID {
			config.ImageMsgUploadID[index] = strconv.FormatInt(id+int64(index), 10)
		}
		common.Dumps(*TaskConfigPath, config)
	}
}

func main() {
	config2.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitSendMsgDB(config.TargetUserDB, config.TargetUserCollection)
	WaitAll.Add(config.CoroCount + 1)
	go SendUser()
	for index := 0; index < config.CoroCount; index++ {
		go SendTask()
	}
	WaitAll.Wait()
}
