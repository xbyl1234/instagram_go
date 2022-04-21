package main

import (
	"encoding/json"
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxys"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

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
var WaitExit sync.WaitGroup
var SuccessCount int32 = 0
var ErrorCount int32 = 0
var (
	TaskLogin       = "relogin"
	TaskRefreshInfo = "refresh_info"
	TaskTestAccount = "test"
	TaskSetEmail    = "email"
	TaskSetBio      = "bio"
)
var (
	SendAll          = "all"
	SendGood         = "good"
	SendBad          = "bad"
	SendNoLogin      = "nologin"
	SendStatusErr    = "badstat"
	SendReqErr       = "badreq"
	SendNoDevice     = "nodevice"
	SendOldStruct    = "odlstruct"
	SendTagsCrawTags = "craw_tags"
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
		//if err != nil {
		//	log.Warn("username: %s, init error: %v", inst.User, err.Error())
		//	return inst, err
		//}

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

func InstRelogin(inst *goinsta.Instagram) error {
	inst, err := Login(inst.User, inst.Pass)
	if err != nil {
		return err
	}

	err = inst.GetAccount().Sync()
	return err
}

func InstRefreshAccountInfo(inst *goinsta.Instagram) error {
	uploadID, _, err := inst.GetUpload().UploadPhotoFromPath(common.Resource.ChoiceIco(), nil)
	if err != nil {
		log.Error("account %s upload ico error: %v", inst.User, err)
		return err
	}

	err = inst.GetAccount().ChangeProfilePicture(uploadID)
	if err != nil {
		log.Error("account %s change ico error: %v", inst.User, err)
	}
	return err
}

func InstSetBio(inst *goinsta.Instagram) error {
	err := inst.GetAccount().Sync()
	if err != nil {
		return err
	}
	err = inst.GetAccount().EditProfile(&goinsta.UserProfile{
		//ExternalUrl: "http://sexy37.com/" + fmt.Sprintf("%d", inst.ID),
		Biography: common.GenString(common.CharSet_abc, 5) + "I have the pictures you want on my blog~ ",
	})
	if err != nil {
		return err
	}
	return nil
}

func InstTestAccount(inst *goinsta.Instagram) error {
	var err error
	if inst.ID == 0 || inst.IsLogin == false {
		inst, err = Login(inst.User, inst.Pass)
		if err != nil {
			return err
		}
	}

	//The password you entered is incorrect
	//invalid character 'O' looking for beginning of value
	//The username you entered doesn't appear to belong to an account
	//invalid character '<' looking

	err = inst.GetAccount().Sync()
	if err != nil {
		log.Error("account: %s, error: %v", inst.User, err)
	} else {
		//log.Info("account %s website %s bio %s", inst.User, inst.GetAccount().Detail.)
	}
	if err.Error() == "login_required" {
		err = InstRelogin(inst)
		//inst.IsLogin = false
		//err = inst.Login()
		if err != nil {
			log.Error("account: %s, error: %v", inst.User, err)
		}
	}

	return err
}

//root:Hty741852..@tcp(127.0.0.1:7707)/email?readTimeout=10s&writeTimeout=10
func InstEmail(inst *goinsta.Instagram) error {
	//err := inst.GetAccount().Sync()
	email := inst.User + "@insemail.work"
	err := inst.GetAccount().SendConfirmEmail(email)
	if err != nil {
		log.Error("send email %s error: %v", email, err)
		return err
	}
	return nil
}

func DispatchAccount() {
	defer WaitExit.Done()
	var Consumer func(inst *goinsta.Instagram) error
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
	case TaskSetEmail:
		Consumer = InstEmail
		break
	case TaskSetBio:
		Consumer = InstSetBio
		break
	default:
		return
	}

	for inst := range TestAccount {
		err := Consumer(inst)
		if err != nil {
			atomic.AddInt32(&ErrorCount, 1)
		} else {
			atomic.AddInt32(&SuccessCount, 1)
		}
		goinsta.SaveInstToDB(inst)
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
		break
	case SendNoDevice:
		if inst.IsLogin && inst.ID != 0 && inst.Status == "" && inst.AccountInfo.Device.DeviceID == "" {
			TestAccount <- inst
		}
		break
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
	case SendOldStruct:
		if inst.ID == 0 && inst.Status == "" {
			inst.CleanCookiesAndHeader()
			inst.AccountInfo = goinsta.GenInstDeviceInfo()
			TestAccount <- inst
		}
		break
	case SendReqErr:
		if strings.Index(inst.Status, "invalid character") != -1 {
			inst.CleanCookiesAndHeader()
			TestAccount <- inst
		}
		break
		//case SendTagsCrawTags:
		//
		//	if  inst. == "" {
		//
		//	}
		//	"craw_tags"
		//	break
	}
}

func SendAccount(insts []*goinsta.Instagram) {
	for index := range insts {
		send(insts[index])
	}

	close(TestAccount)
	WaitExit.Done()
}

func TestDevice() {
	for true {
		_proxy := proxys.ProxyPool.Get("us", "")
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			continue
		}

		inst := goinsta.New("", "", _proxy)
		accInfo, _ := json.Marshal(inst.AccountInfo)
		err := inst.QeSync()
		if err != nil {
			log.Error("err: %s", accInfo)
		} else {
			log.Info("suc: %s", accInfo)
		}
	}
}

func main() {
	goinsta.UsePanic = false
	common.UseCharles = false
	common.UseTruncation = false

	initParams()
	routine.InitRoutine(config.ProxyPath)
	var err error

	//for i := 0; i < 10; i++ {
	//	go TestDevice()
	//}
	//select {}
	//login, err := Login("impatient2017116", "KJVEkjve8752")
	//if err != nil {
	//	return
	//}
	//err = login.GetAccount().Sync()
	//if err != nil {
	//	log.Error("account: %s, error: %v", login.User, err)
	//}
	//goinsta.SaveInstToDB(login)
	//goinsta.CleanStatus()
	//goinsta.ReStruct()
	//return
	err = common.InitResource(config.ResIcoPath, "")
	if err != nil {
		log.Error("load res error: %v", err)
		os.Exit(0)
	}

	//insts := goinsta.LoadAllAccount()
	insts := goinsta.LoadAccountByTags([]string{"comment"})
	if len(insts) == 0 {
		log.Error("there have no account!")
		os.Exit(0)
	}
	log.Info("load account count: %d", len(insts))

	if *TestOne != "" {
		for _, item := range insts {
			if item.User == *TestOne {
				send(item)
				DispatchAccount()
			}
		}
		os.Exit(0)
	}

	WaitExit.Add(1 + config.Coro)
	go SendAccount(insts)

	for i := 0; i < config.Coro; i++ {
		go DispatchAccount()
	}

	WaitExit.Wait()

	log.Info("test finish")
}
