package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
	"sync"
)

type TestLoginResult struct {
	inst   *goinsta.Instagram
	status bool
	err    error
	str    string
}
type Config struct {
	ProxyPath  string `json:"proxy_path"`
	ResIcoPath string `json:"res_ico_path"`
	Coro       int    `json:"coro"`
}

var config Config

var TaskType = flag.String("type", "", "")
var TestOne = flag.String("test_one", "", "")
var ConfigPath = flag.String("config", "./config/refresh.json", "")

var TestAccount = make(chan *goinsta.Instagram, 1000)
var TestResult = make(chan *TestLoginResult, 1000)
var TestResultList []*TestLoginResult
var WaitTask sync.WaitGroup
var WaitExit sync.WaitGroup
var (
	TaskLogin       = "login"
	TaskRefreshInfo = "refresh_info"
	TaskTestAccount = "test"
)

func initParams() {
	flag.Parse()
	log.InitDefaultLog("refresh_cookies", true, true)
	err := common.LoadJsonFile(*ConfigPath, &config)
	if err != nil {
		log.Error("load config error: %v", err)
		os.Exit(0)
	}

	if config.ProxyPath == "" {
		log.Error("proxy path is null")
		os.Exit(0)
	}
	if config.ResIcoPath == "" {
		log.Error("ResourcePath path is null")
		os.Exit(0)
	}
	if config.Coro == 0 {
		config.Coro = 1
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
				proxy.ProxyPool.Black(result.inst.Proxy, proxy.BlackType_RegisterRisk)
			}
		} else {
			result.inst.Status = ""
			result.inst.IsLogin = true
		}
		goinsta.SaveInstToDB(result.inst)
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
	proxy.ProxyPool.Dumps()
}

func InstTestAccount(inst *goinsta.Instagram) *TestLoginResult {
	var result = &TestLoginResult{}
	result.status = true
	result.inst = inst
	if inst.Status == "challenge_required" {
		result.status = false
		return result
	}

	if routine.SetProxy(inst) {
		if inst.ID == 0 || inst.IsLogin == false {
			inst.CleanCookiesAndHeader()
			inst.PrepareNewClient()
			err := Login(inst)
			if err != nil {
				result.err = err
				if common.IsError(err, common.ChallengeRequiredError) {
					result.inst.Status = "challenge_required"
					return result
				} else if common.IsError(err, common.RequestError) {
					log.Error("account: %s, request error: %v", inst.User, err)
					result.status = false
					return result
				} else {
					log.Error("account: %s, unknow error: %v", inst.User, err)
					result.inst.Status = err.Error()
					return result
				}
			}
			result.inst.IsLogin = true
		}
		//The password you entered is incorrect
		//invalid character 'O' looking for beginning of value
		//The username you entered doesn't appear to belong to an account
		//invalid character '<' looking

		result.inst.Status = ""
		err := inst.GetAccount().Sync()
		if err != nil {
			result.err = err
			if common.IsError(err, common.ChallengeRequiredError) {
				result.inst.Status = "challenge_required"
				return result
			} else if common.IsError(err, common.RequestError) {
				log.Error("account: %s, request error: %v", inst.User, err)
				result.status = false
				return result
			} else if strings.Index(err.Error(), "invalid character '<' looking") != -1 {
				inst.CleanCookiesAndHeader()
				result.inst.Status = err.Error()
				return result
			} else {
				log.Error("account: %s, unknow error: %v", inst.User, err)
				result.inst.Status = err.Error()
				return result
			}
		} else {
			inst.Status = ""
		}
	} else {
		result.str = "no proxy"
		result.status = false
	}

	return result
}

func RecvTestAccount() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++
		if result.status {
			goinsta.SaveInstToDB(result.inst)
		}
	}
	WaitExit.Done()
	PrintResult(TestResultList[:index])
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
	case TaskTestAccount:
		Consumer = InstTestAccount
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
func testOne(insts []*goinsta.Instagram, username string) {
	var inst *goinsta.Instagram
	for index := range insts {
		if insts[index].User == username {
			inst = insts[index]
			break
		}
	}
	if inst == nil {
		log.Error("not find user: %s", username)
		os.Exit(0)
	}

	result := InstTestAccount(inst)
	if result.status {
		if result.err == nil {
			log.Info("islogin: %v, acc status: %s", result.inst.IsLogin, result.inst.Status)
		} else {
			log.Info("acc: %s, str: %s, err: %v", result.inst.Status, result.str, result.err)
		}
		goinsta.SaveInstToDB(inst)
	} else {
		log.Info("result status is false, str: %s, err: %v", result.str, result.err)
	}
}

func main() {
	config2.UseCharles = false
	config2.UseTruncation = false

	initParams()
	routine.InitRoutine(config.ProxyPath)

	//goinsta.CleanStatus()
	err := common.InitResource(config.ResIcoPath, "")
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

	if *TestOne != "" {
		testOne(insts, *TestOne)
		os.Exit(0)
	}

	WaitExit.Add(2)
	WaitTask.Add(config.Coro)
	go SendAccount(insts)
	switch *TaskType {
	case TaskLogin:
		go RecvCleanAndLogin()
		break
	case TaskRefreshInfo:
		go RecvRefreshAccountInfo()
		break
	case TaskTestAccount:
		go RecvTestAccount()
		break
	default:
		log.Error("task type error")
		os.Exit(0)
	}

	for i := 0; i < config.Coro; i++ {
		go DispatchAccount()
	}

	WaitExit.Wait()

	log.Info("test finish")
}
