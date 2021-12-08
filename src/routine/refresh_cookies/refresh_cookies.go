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
	inst   *goinsta.Instagram
	status bool
	err    error
	str    string
}

var ProxyPath = flag.String("proxy", "", "")
var ResIcoPath = flag.String("ico", "", "")
var Coro = flag.Int("coro", 16, "")
var TaskType = flag.String("type", "", "")

var TestAccount = make(chan *goinsta.Instagram, 1000)
var TestResult = make(chan *TestLoginResult, 1000)
var TestResultList []*TestLoginResult
var WaitTask sync.WaitGroup
var WaitExit sync.WaitGroup
var (
	TaskLogin       = "login"
	TaskRefreshInfo = "refresh_info"
)

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

func Login(inst *goinsta.Instagram) error {
	err := inst.Login()
	if err != nil {
		log.Warn("username: %s, login error: %v", inst.User, err.Error())
		return err
	}
	log.Info("username: %s, login success", inst.User)
	return nil
}

func InstCleanAndLogin(inst *goinsta.Instagram) *TestLoginResult {
	result := &TestLoginResult{}
	result.inst = inst
	if !inst.IsLogin && inst.Status == "" {
		if routine.SetProxy(inst) {
			//inst.CleanCookiesAndHeader()
			//inst.PrepareNewClient()

			if result.status == false {
				err := Login(inst)
				if err != nil {
					result.str = "login error"
					result.status = false
					result.err = err
				} else {
					result.status = true
				}
			}
		} else {
			result.str = "no proxy"
			result.status = false
		}
		return result
	}

	return nil
}

func RecvCleanAndLogin() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++

		if result.str != "no proxy" {
			if common.IsError(result.err, common.ChallengeRequiredError) {
				result.inst.Status = goinsta.InsAccountError_ChallengeRequired
			}
			//result.inst.IsLogin = result.status
			goinsta.SaveInstToDB(result.inst)
		}
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
}

func InstRefreshAccountInfo(inst *goinsta.Instagram) *TestLoginResult {
	if !inst.IsLogin && inst.Status == "challenge_required" {
		var result = &TestLoginResult{}
		inst.IsLogin = true
		if routine.SetProxy(inst) {
			err := inst.GetAccount().ChangeProfilePicture(common.Resource.ChoiceIco())
			result.inst = inst
			if err == nil {
				result.status = true
			} else {
				result.status = false
				result.err = err
			}
		} else {
			result.str = "no proxy"
			result.status = false
		}
		return result
	} else {
		return nil
	}
}

func RecvRefreshAccountInfo() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++
		if !result.status {
			if result.str == "no proxy" {
				continue
			}

			if common.IsError(result.err, common.ChallengeRequiredError) {
				result.inst.Status = "challenge_required"
				result.inst.IsLogin = false
				common.ProxyPool.Black(result.inst.Proxy, common.BlackType_RegisterRisk)
			}
		} else {
			result.inst.Status = ""
			result.inst.IsLogin = true
		}
		goinsta.SaveInstToDB(result.inst)
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
	common.ProxyPool.Dumps()
}

func DispatchAccount() {
	defer WaitTask.Done()
	var Consumer func(inst *goinsta.Instagram) *TestLoginResult
	switch *TaskType {
	case TaskLogin:
		Consumer = InstCleanAndLogin
		break
	case TaskRefreshInfo:
		Consumer = InstRefreshAccountInfo
		break
	default:
		return
	}

	for inst := range TestAccount {
		result := Consumer(inst)
		if result != nil {
			TestResult <- result
		}
	}
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

func PrintResult(result []*TestLoginResult) {
	log.Info("---------------  success   ---------------")
	for index := range result {
		if result[index].status {
			log.Info("username: %s", result[index].inst.User)
		}
	}
	log.Info("-------------    failed   --------------")
	for index := range result {
		if !result[index].status {
			log.Error("username: %s, %s, err: %v", result[index].inst.User, result[index].str, result[index].err)
		}
	}
	log.Info("--------------- proxy error --------------")
	for index := range result {
		if result[index].str == "no proxy" {
			log.Warn("username: %s, %s", result[index].inst.User, result[index].str)
		}
	}
}

func main() {
	config.UseCharles = false
	config.UseTruncation = false

	initParams()
	common.InitMogoDB()
	routine.InitRoutine(*ProxyPath)

	//goinsta.CleanStatus()
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
	switch *TaskType {
	case TaskLogin:
		go RecvCleanAndLogin()
		break
	case TaskRefreshInfo:
		go RecvRefreshAccountInfo()
		break
	default:
		log.Error("task type error")
		os.Exit(0)
	}

	for i := 0; i < *Coro; i++ {
		go DispatchAccount()
	}

	WaitExit.Wait()

	log.Info("test finish")
}
