package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type CrawConfig struct {
	TaskName             string `json:"task_name"`
	CoroCount            int    `json:"coro_count"`
	ProxyPath            string `json:"proxy_path"`
	WorkPath             string `json:"config_path"`
	TargetUserDB         string `json:"target_user_db"`
	TargetUserCollection string `json:"target_user_collection"`
}

var TargetUserChan = make(chan *routine.UserComb, 1000)
var WaitAll sync.WaitGroup
var config CrawConfig
var WorkPath string
var PathSeparator = string(os.PathSeparator)
var CrawCount int32 = 0

func CrawlingFans() {
	defer WaitAll.Done()

	unknowErrorCount := 0
	inst, err := routine.ReqAccount(false)
	if err != nil {
		log.Error("CrawlingFans req account error: %v", err)
		return
	}

	for item := range TargetUserChan {
		if item.User.ID != 0 {
			var followes *goinsta.Followers
			if item.Followes == nil {
				followes = inst.GetFollowers(strconv.FormatInt(item.User.ID, 10))
				item.Followes = followes
			} else {
				followes = item.Followes
				followes.SetAccount(inst)
			}

			targetUser, err := followes.Next()
			if err != nil {
				if common.IsNoMoreError(err) {
					log.Info("user %v fans comments has craw finish!", followes.User)
					continue
				} else if inst.NeedReplace() || common.IsError(err, common.RequestError) {
					if inst.NeedReplace() {
						goinsta.AccountPool.BlackOne(inst)
						_inst, errAcc := routine.ReqAccount(false)
						if errAcc != nil {
							inst = nil
							log.Error("CrawlingFans req account error: %v!", errAcc)
							return
						}
						log.Warn("CrawlingFans replace account %s->%s", inst.User, _inst.User)
						inst = _inst
						followes.SetAccount(_inst)
					} else {
						log.Warn("CrawlingFans retrying...user: %s, err: %v", inst.User, err)
					}
					continue
				} else {
					unknowErrorCount++
					log.Error("Next Fans error:%v", err)
					continue
				}
			} else {
				unknowErrorCount = 0
			}

			for userIndex := range targetUser {
				var userComb routine.UserComb
				userComb.User = &targetUser[userIndex]
				userComb.Source = strconv.FormatInt(item.User.ID, 10)
				err = routine.SaveUser(routine.CrawFansUserColl, &userComb)
				if err != nil {
					log.Error("save target user error! %v", err)
				}
			}

			err = routine.SaveUser(routine.CrawFansTargetUserColl, item)
			if err != nil {
				log.Error("save target user error! %v", err)
			}
		}
	}
}

func LoadTargetUser() {
	defer WaitAll.Done()
	for true {
		result, err := routine.LoadFansTargetUser(1000)
		if err != nil {
			log.Error("load target user error: %v", err)
			break
		}
		if len(result) == 0 {
			log.Info("craw target user finish!")
			break
		}

		for index := range result {
			TargetUserChan <- &result[index]
		}
	}

	close(TargetUserChan)
}

func initParams() {
	var err error
	log.InitDefaultLog("craw_fans", true, true)
	var TaskConfigPath = flag.String("task", "", "task")
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
		config.CoroCount = runtime.NumCPU()*2 + 1
	}

	WorkPath, _ = os.Getwd()
	if config.WorkPath == "" {
		config.WorkPath = WorkPath + PathSeparator + config.TaskName + PathSeparator
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
	routine.InitCrawFansDB(config.TaskName, config.TargetUserDB, config.TargetUserCollection)
	WaitAll.Add(config.CoroCount + 1)
	go LoadTargetUser()
	for index := 0; index < config.CoroCount; index++ {
		go CrawlingFans()
	}
	log.Info("craw finish! count %d", CrawCount)
}
