package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
)

type MakeMoneyConfig struct {
	TaskName             string       `json:"task_name"`
	CoroCount            int          `json:"coro_count"`
	ProxyPath            string       `json:"proxy_path"`
	TargetUserDB         string       `json:"target_user_db"`
	TargetUserCollection string       `json:"target_user_collection"`
	AccountTag           string       `json:"account_tag"`
	Msgs                 [][]*Message `json:"msgs"`
}

var config MakeMoneyConfig

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
	if config.Msgs == nil || len(config.Msgs) == 0 {
		log.Error("Msgs is null")
		os.Exit(0)
	}
	err = LoadScreenplay()
	if err != nil {
		log.Error("LoadScreenplay error: %v", err)
		os.Exit(0)
	}
}

func main() {
	config2.UseCharles = true
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitSendMsgDB(config.TargetUserDB, config.TargetUserCollection)

	intas := goinsta.LoadAccountByTags([]string{config.AccountTag, config.AccountTag})
	if len(intas) == 0 {
		log.Warn("there have no account!")
	} else {
		goinsta.InitAccountPool(intas)
	}

	ShortVideoTask()
}
