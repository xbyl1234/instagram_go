package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
)

type tmpAccount struct {
	username string
	passwd   string
}
type TestLoginResult struct {
	inst    *goinsta.Instagram
	IsLogin bool
	err     error
	str     string
}

var ProxyPath = flag.String("proxy", "", "")
var ResIcoPath = flag.String("ico", "", "")
var TestIsLogin = flag.Bool("test_login", false, "")
var Coro = flag.Int("coro", 8, "")

var TestAccount = make(chan *goinsta.Instagram, 1000)
var TestResult = make(chan *TestLoginResult, 1000)
var TestResultList []*TestLoginResult
var WaitTask sync.WaitGroup
var WaitExit sync.WaitGroup

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

func SetProxy(inst *goinsta.Instagram) bool {
	var _proxy *common.Proxy
	if inst.Proxy.ID != "" {
		_proxy = common.ProxyPool.Get(inst.Proxy.ID)
		if _proxy == nil {
			log.Warn("find insta proxy %s error!", inst.Proxy.ID)
		}
	}

	if _proxy == nil {
		_proxy = common.ProxyPool.GetNoRisk(false, false)
		if _proxy == nil {
			log.Error("get insta proxy error!")
		}
	}

	if _proxy != nil {
		inst.SetProxy(_proxy)
	}
	return true
}

func Login(inst *goinsta.Instagram) error {
	err := inst.Login()
	if err != nil {
		log.Warn("username: %s, login error: %v", inst.User, err.Error())
		return err
	}
	log.Info("username: %s, login success", inst.User)
	return nil
}

func TestAndLogin() {
	for inst := range TestAccount {
		result := &TestLoginResult{}
		result.inst = inst
		if !inst.IsLogin {
			if SetProxy(inst) {
				inst.CleanCookiesAndHeader()
				inst.PrepareNewClient()
				//acc := inst.GetAccount()
				//err := acc.Sync()
				//if err != nil {
				//	result.str = "account sync error"
				//	result.IsLogin = false
				//	result.err = err
				//} else {
				//	result.IsLogin = true
				//}

				if result.IsLogin == false {
					err := Login(inst)
					if err != nil {
						result.str = "login error"
						result.IsLogin = false
						result.err = err
					} else {
						result.IsLogin = true
					}
				}
			} else {
				result.str = "no proxy"
				result.IsLogin = false
			}
			TestResult <- result
		}
	}
	WaitTask.Done()
}

func SendAccount(insts []*goinsta.Instagram) {
	for index := range insts {
		TestAccount <- insts[index]
	}

	close(TestAccount)
	WaitTask.Wait()
	close(TestResult)
	WaitExit.Done()
}

func RecvAccount() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++

		if result.str != "no proxy" {
			if common.IsError(result.err, common.ChallengeRequiredError) {
				result.inst.Status = goinsta.InsAccountError_ChallengeRequired
			}
			result.inst.IsLogin = result.IsLogin
			goinsta.SaveInstToDB(result.inst)
		}
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
}

func PrintResult(result []*TestLoginResult) {
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
}

func main() {
	config.UseCharles = true

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
	TestResultList = make([]*TestLoginResult, len(insts))

	WaitExit.Add(2)
	WaitTask.Add(*Coro)
	go SendAccount(insts)
	go RecvAccount()

	for i := 0; i < *Coro; i++ {
		go TestAndLogin()
	}

	WaitExit.Wait()

	log.Info("test finish")
}
