package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
)

type MakeMoneyConfig struct {
	TaskName             string               `json:"task_name"`
	Coro                 int                  `json:"coro"`
	ProxyPath            string               `json:"proxy_path"`
	TargetUserDB         string               `json:"target_user_db"`
	TargetUserCollection string               `json:"target_user_collection"`
	AccountTag           string               `json:"account_tag"`
	Msgs                 [][]*Message         `json:"msgs"`
	Develop              DevelopAccountConfig `json:"develop"`
}

var config MakeMoneyConfig

func initParams() {
	var err error
	var configPath = flag.String("config", "./config/make_money.json", "task")
	log.InitDefaultLog("send_msg", true, true)
	flag.Parse()
	if *configPath == "" {
		log.Error("config path is null!")
		os.Exit(0)
	}
	err = common.LoadJsonFile(*configPath, &config)
	if err != nil {
		log.Error("load task config error: %v", err)
		os.Exit(0)
	}
	if config.TaskName == "" {
		log.Error("task name is null")
		os.Exit(0)
	}
	if config.Coro == 0 {
		config.Coro = 1
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
	common.UseCharles = false
	goinsta.UsePanic = true
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitSendMsgDB(config.TargetUserDB, config.TargetUserCollection, config.Develop.LogCollName)

	DevelopAccount()
}
