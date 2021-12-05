package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/phone"
	"makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

var ProxyPath = flag.String("proxy", "", "")
var ResIcoPath = flag.String("ico", "", "")
var ResUsernamePath = flag.String("name", "", "")
var Coro = flag.Int("coro", runtime.NumCPU(), "")
var RegisterCount = flag.Int("count", 0, "")

var Count int32 = 0
var SuccessCount = 0
var ErrorCreateCount = 0
var ErrorSendSMSCount = 0
var ErrorRecvSMSCount = 0
var ErrorOtherCount = 0

var PhoneProvider phone.PhoneVerificationCode

var WaitAll sync.WaitGroup

func Register() {
	for true {
		curCount := atomic.AddInt32(&Count, 1)
		if curCount > int32(*RegisterCount) {
			break
		}

		_proxy := common.ProxyPool.GetNoRisk(true, true)
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			break
		}
		log.Info("get proxy ip: %s", _proxy.Rip)

		regisert := goinsta.NewRegister(_proxy, PhoneProvider)
		username := common.Resource.ChoiceUsername()
		password := common.GenString(common.CharSet_ABC, 4) +
			common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)

		inst, err := regisert.Do(username, username, password)
		if err == nil {
			log.Info("register success, username %s, passwd %s", inst.User, inst.Pass)
			err = goinsta.SaveInstToDB(inst)
			if err != nil {
				log.Error("save inst: %s %s error: %v", inst.User, inst.Pass, err)
			}
			//err = inst.GetAccount().ChangeProfilePicture(common.Resource.ChoiceIco())
			//if err != nil {
			//	log.Error("user: %s, change ico error: %v", inst.User, err)
			//}
			SuccessCount++
		} else {
			if common.IsError(err, common.ApiError) {
				if strings.Index(err.Error(), "wait a few minutes") != -1 || strings.Index(err.Error(), "请稍等几分钟再试") != -1 {
					common.ProxyPool.Black(_proxy, common.BlackType_RegisterRisk)
				} else if strings.Index(err.Error(), "feedback_required") != -1 {
					common.ProxyPool.Black(_proxy, common.BlackType_RegisterRisk)
				}
			}

			if !regisert.HadSendSMS {
				ErrorSendSMSCount++
			} else if regisert.HadSendSMS && !regisert.HadRecvSMS {
				ErrorRecvSMSCount++
			} else if regisert.HadSendSMS && regisert.HadRecvSMS {
				ErrorCreateCount++
			} else {
				ErrorOtherCount++
			}
			//challenge_required
			//feedback_required
			log.Warn("register error, %v", err)
		}
	}
	WaitAll.Done()
}

func initParams() {
	flag.Parse()
	log.InitDefaultLog("register", true, true)
	if *ProxyPath == "" {
		log.Error("proxy path is null")
		os.Exit(0)
	}
	if *ResIcoPath == "" {
		log.Error("ResourcePath is null")
		os.Exit(0)
	}
	if *ResUsernamePath == "" {
		log.Error("ResUsernamePath is null")
		os.Exit(0)
	}
	if *RegisterCount == 0 {
		log.Error("RegisterCount is 0")
		os.Exit(0)
	}
}

//girlchina001
//a123456789
func main() {
	config.UseCharles = false
	config.UseTruncation = false

	initParams()
	common.InitMogoDB()
	routine.InitRoutine(*ProxyPath)

	var err error
	PhoneProvider, err = phone.NewPhoneVerificationCode("do889")
	if err != nil {
		log.Error("create phone provider error!%v", err)
		os.Exit(0)
	}
	err = common.InitResource(*ResIcoPath, *ResUsernamePath)
	if err != nil {
		log.Error("InitResource error!%v", err)
		os.Exit(0)
	}

	WaitAll.Add(*Coro)
	for i := 0; i < *Coro; i++ {
		go Register()
	}
	WaitAll.Wait()
	log.Info("success: %d, create err: %d, send msg err: %d, recv msg err: %d, other err: %d",
		SuccessCount,
		ErrorCreateCount,
		ErrorSendSMSCount,
		ErrorRecvSMSCount,
		ErrorOtherCount)
}
