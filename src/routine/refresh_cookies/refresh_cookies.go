package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
	"sync"
)

type TestLoginResult struct {
	inst *goinsta.Instagram
	err  error
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
var SendAllAccount = flag.String("acc", "all", "")

var TestAccount = make(chan *goinsta.Instagram, 1000)
var TestResult = make(chan *TestLoginResult, 1000)
var TestResultList []*TestLoginResult
var WaitTask sync.WaitGroup
var WaitExit sync.WaitGroup
var (
	TaskLogin       = "relogin"
	TaskRefreshInfo = "refresh_info"
	TaskTestAccount = "test"
)
var (
	SendAll       = "all"
	SendGood      = "good"
	SendBad       = "bad"
	SendNoLogin   = "nologin"
	SendStatusErr = "badstat"
	SendReqErr    = "badreq"
	SendNoDevice  = "nodevice"
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

func Login(username string, password string) (*goinsta.Instagram, error) {
	var inst = goinsta.New(username, password, nil)
	var err error

	if routine.SetProxy(inst) {
		inst.PrepareNewClient()
		err = inst.Login()
		if err != nil {
			log.Warn("username: %s, login error: %v", inst.User, err.Error())
			return inst, err
		}
		log.Info("username: %s, login success", inst.User)
		return inst, nil
	} else {
		return inst, &common.MakeMoneyError{
			ErrStr:    "no proxy",
			ErrType:   0,
			ExternErr: nil,
		}
	}
}

func InstRelogin(inst *goinsta.Instagram) *TestLoginResult {
	result := &TestLoginResult{}
	var err error
	result.inst, err = Login(inst.User, inst.Pass)

	if err != nil {
		result.err = err
		return result
	}

	err = result.inst.GetAccount().Sync()
	if err != nil {
		result.err = err
	}
	return result
}

func RecvCleanAndLogin() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++
		goinsta.SaveInstToDB(result.inst)
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
}

func InstRefreshAccountInfo(inst *goinsta.Instagram) *TestLoginResult {
	if !inst.IsLogin && inst.Status == "challenge_required" {
		var result = &TestLoginResult{}
		inst.IsLogin = true
		if routine.SetProxy(inst) {
			uploadID, err := inst.GetUpload().RuploadPhotoFromPath(common.Resource.ChoiceIco())
			if err == nil {
				err = inst.GetAccount().EditProfile(&goinsta.UserProfile{
					UploadId: uploadID,
				})
			}
			result.inst = inst
			if err == nil {
			} else {
				result.err = err
			}
		} else {
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
		//if !result.status {
		//	if result.err.Error() == "no proxy" {
		//		continue
		//	}
		//
		//	if common.IsError(result.err, common.ChallengeRequiredError) {
		//		result.inst.Status = "challenge_required"
		//		result.inst.IsLogin = false
		//		proxy.ProxyPool.Black(result.inst.Proxy, proxy.BlackType_RegisterRisk)
		//	}
		//} else {
		//	result.inst.Status = ""
		//	result.inst.IsLogin = true
		//}
		goinsta.SaveInstToDB(result.inst)
	}
	PrintResult(TestResultList[:index])
	WaitExit.Done()
	//proxy.ProxyPool.Dumps()
}

func InstTestAccount(inst *goinsta.Instagram) *TestLoginResult {
	var err error
	var result = &TestLoginResult{}
	result.inst = inst

	if routine.SetProxy(inst) {
		if inst.ID == 0 || inst.IsLogin == false {
			result.inst, err = Login(inst.User, inst.Pass)
			if err != nil {
				result.err = err
				return result
			}
		}

		//The password you entered is incorrect
		//invalid character 'O' looking for beginning of value
		//The username you entered doesn't appear to belong to an account
		//invalid character '<' looking

		result.inst.Status = ""
		err = result.inst.GetAccount().Sync()
		if err != nil {
			result.err = err
			log.Error("account: %s, error: %v", inst.User, err)
			return result
			//if strings.Index(err.Error(), "invalid character") != -1 {
			//	log.Error("account: %s, cookies error: %v", inst.User, err)
			//	inst.CleanCookiesAndHeader()
			//	continue
			//}
		}
	} else {
		result.err = &common.MakeMoneyError{
			ErrStr: "no proxy",
		}
	}

	return result
}

func RecvTestAccount() {
	index := 0
	for result := range TestResult {
		TestResultList[index] = result
		index++
		if result.err.Error() != "no proxy" {
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
		Consumer = InstRelogin
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

func send(inst *goinsta.Instagram) {
	switch *SendAllAccount {
	case SendAll:
		TestAccount <- inst
		break
	case SendGood:
		if inst.IsLogin && inst.ID != 0 && inst.Status == "" {
			TestAccount <- inst
		}
	case SendNoDevice:
		if inst.IsLogin && inst.ID != 0 && inst.Status == "" && inst.DeviceID == "" {
			TestAccount <- inst
		}
	case SendBad:
		if !inst.IsLogin || inst.ID == 0 || inst.Status != "" {
			TestAccount <- inst
		}
		break
	case SendNoLogin:
		if !inst.IsLogin {
			TestAccount <- inst
		}
		break
	case SendStatusErr:
		if inst.Status != "" {
			TestAccount <- inst
		}
		break
	case SendReqErr:
		if strings.Index(inst.Status, "invalid character") != -1 {
			inst.CleanCookiesAndHeader()
			TestAccount <- inst
		}
		break
	}
}

func SendAccount(insts []*goinsta.Instagram) {
	for index := range insts {
		send(insts[index])
	}

	close(TestAccount)
	WaitTask.Wait()
	close(TestResult)
	WaitExit.Done()
}

func PrintResult(result []*TestLoginResult) {
	log.Info("---------------  success   ---------------")
	for index := range result {
		if result[index].err == nil {
			log.Info("username: %s", result[index].inst.User)
		}
	}
	log.Info("-------------    failed   --------------")
	for index := range result {
		if result[index].err != nil && result[index].err.Error() != "no proxy" {
			log.Error("username: %s, err: %v", result[index].inst.User, result[index].err)
		}
	}
	log.Info("--------------- proxy error --------------")
	for index := range result {
		if result[index].err != nil && result[index].err.Error() == "no proxy" {
			log.Warn("username: %s", result[index].inst.User)
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

	if result.err == nil {
		log.Info("islogin: %v, acc status: %s", result.inst.IsLogin, result.inst.Status)
	} else {
		log.Info("account: %s, err: %v", result.inst.Status, result.err)
	}
	goinsta.SaveInstToDB(inst)
}

func main() {
	config2.UseCharles = false
	config2.UseTruncation = true

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
