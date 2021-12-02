package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
)

type tmpAccount struct {
	username string
	passwd   string
}

var ProxyPath = flag.String("proxy", "", "")
var ResIcoPath = flag.String("ico", "", "")
var TestIsLogin = flag.Bool("test_login", false, "")

func initParams() {
	flag.Parse()
	log.InitDefaultLog("refresh_cookies", true, true)
	if *ProxyPath == "" {
		log.Error("proxy path is null")
		os.Exit(0)
	}
	if *ResIcoPath == "" {
		log.Error("ResourcePath path is null")
		os.Exit(0)
	}
}

type TestLoginResult struct {
	inst    *goinsta.Instagram
	IsLogin bool
	err     error
	str     string
}

func InitAccount(inst *goinsta.Instagram) bool {
	if inst.Proxy.ID != "" {
		_proxy := common.ProxyPool.Get(inst.Proxy.ID)
		if _proxy == nil {
			log.Error("find insta proxy error!")
			return false
		}
		inst.SetProxy(_proxy)
	} else {
		_proxy := common.ProxyPool.GetOne()
		if _proxy == nil {
			log.Error("find insta proxy error!")
			return false
		}
		inst.SetProxy(_proxy)
	}
	return true
}

func main() {
	config.UseCharles = false

	initParams()
	common.InitMogoDB()
	routine.InitRoutine(*ProxyPath)
	err := common.InitResource(*ResIcoPath, "")
	if err != nil {
		log.Error("load res error: %v", err)
		os.Exit(0)
	}

	insts := goinsta.LoadAllAccount()
	if len(insts) == 0 {
		log.Error("there have no account!")
		os.Exit(0)
	}
	log.Info("load account count: %d", len(insts))
	result := make([]TestLoginResult, len(insts))

	if *TestIsLogin {
		for index := range insts {
			inst := insts[index]
			result[index].inst = inst

			if !InitAccount(inst) {
				result[index].str = "no proxy"
				result[index].IsLogin = false
				continue
			}

			if inst.ID == 0 {
				result[index].str = "id is 0,not login"
				result[index].IsLogin = false
				continue
			}
			acc := inst.GetAccount()
			err := acc.Sync()
			if err != nil {
				result[index].str = "account sync error"
				result[index].IsLogin = false
				result[index].err = err
			} else {
				result[index].IsLogin = true
			}
		}
	}

	log.Info("test finish")

	log.Info("---------------login account---------------")
	for index := range result {
		if result[index].IsLogin {
			log.Info("username: %s", result[index].inst.User)
		}
	}
	log.Info("--------------- proxy error --------------")
	for index := range result {
		if result[index].str == "no proxy" {
			log.Warn("username: %s, %s", result[index].inst.User, result[index].str)
		}
	}

	log.Info("-------------not login account--------------")
	for index := range result {
		if !result[index].IsLogin {
			log.Error("username: %s, %s, err: %v", result[index].inst.User, result[index].str, result[index].err)
		}
	}

	for index := range result {
		if result[index].str == "no proxy" {
			result[index].inst.IsLogin = result[index].IsLogin
			goinsta.SaveInstToDB(result[index].inst)
		}
	}
	//for item := accounts.Front(); item != nil; item = item.Next() {
	//	acc := item.Value.(*tmpAccount)
	//	inst := goinsta.New(acc.username, acc.passwd, common.ProxyPool.GetOne())
	//	inst.PrepareNewClient()
	//	err := inst.Login()
	//	if err != nil {
	//		log.Warn("username: %s, login error: %v", acc.username, err.Error())
	//	} else {
	//		log.Info("username: %s, login success", acc.username)
	//	}
	//	_ = goinsta.SaveInstToDB(inst)
	//}
}
