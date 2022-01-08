package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/phone"
	"makemoney/common/proxy"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	ProxyPath       string `json:"proxy_path"`
	ResIcoPath      string `json:"res_ico_path"`
	ResUsernamePath string `json:"res_username_path"`
	Coro            int    `json:"coro"`
	Country         string `json:"country"`
}

var ConfigPath = flag.String("config", "./config/register.json", "")
var RegisterCount = flag.Int("count", 0, "")

var config Config

var Count int32 = 0
var SuccessCount = 0
var ErrorCreateCount = 0
var ErrorSendSMSCount = 0
var ErrorRecvSMSCount = 0
var ErrorOtherCount = 0
var ErrorChallengeRequired = 0

var PhoneProvider phone.PhoneVerificationCode

var WaitAll sync.WaitGroup

var logTicker *time.Ticker

func LogStatus() {
	for range logTicker.C {
		log.Info("success: %d,challenge err: %d ,create err: %d, send msg err: %d, recv msg err: %d, other err: %d",
			SuccessCount,
			ErrorChallengeRequired,
			ErrorCreateCount,
			ErrorSendSMSCount,
			ErrorRecvSMSCount,
			ErrorOtherCount)
	}
}

func Register() {
	for true {
		curCount := atomic.AddInt32(&Count, 1)
		if curCount > int32(*RegisterCount) {
			break
		}
		_proxy := proxy.ProxyPool.GetNoRisk(config.Country, true, true)
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			break
		}
		log.Info("get proxy ip: %s", _proxy.Rip)
		regisert := goinsta.NewRegister(_proxy, PhoneProvider)
		username := common.Resource.ChoiceUsername()
		//username := common.GenString(common.CharSet_abc, 10)
		password := common.GenString(common.CharSet_ABC, 4) +
			common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)

		inst, err := regisert.Do(username, username, password)
		var statErr = err
		if err == nil {
			inst.PrepareNewClient()
			err = inst.GetAccount().Sync()
			if err == nil {
				var uploadID string
				uploadID, err = inst.GetUpload().RuploadPhotoFromPath(common.Resource.ChoiceIco())
				err = inst.GetAccount().EditProfile(&goinsta.UserProfile{
					UploadId: uploadID,
				})
				if err != nil {
					if common.IsError(err, common.ChallengeRequiredError) {
						//proxy.ProxyPool.Black(_proxy, proxy.BlacktypeRegisterrisk)
						ErrorChallengeRequired++
					}
					log.Error("user: %s, change ico error: %v", inst.User, err)
				}
			} else {
				statErr = err
				log.Error("username %s, account sync error: %v", inst.User, inst.Pass, err)
			}

			inst.RegisterTime = time.Now().Unix()
			err = goinsta.SaveInstToDB(inst)
			if err != nil {
				log.Error("save inst: %s %s error: %v", inst.User, inst.Pass, err)
			}
		}

		if statErr != nil {
			if common.IsError(statErr, common.ApiError) {
				if strings.Index(statErr.Error(), "wait a few minutes") != -1 || strings.Index(statErr.Error(), "请稍等几分钟再试") != -1 {
					//proxy.ProxyPool.Black(_proxy, proxy.BlacktypeRegisterrisk)
				} else if strings.Index(statErr.Error(), "feedback_required") != -1 {
					//proxy.ProxyPool.Black(_proxy, proxy.BlacktypeRegisterrisk)
				}
			}

			if common.IsError(statErr, common.ChallengeRequiredError) {
				ErrorChallengeRequired++
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
			log.Warn("register error,username: %s, proxy ip: %s, error: %v", username, _proxy.Rip, statErr)
		} else {
			SuccessCount++
			log.Info("register success, username %s, passwd %s", inst.User, inst.Pass)
		}
	}
	WaitAll.Done()
}

func initParams() {
	flag.Parse()
	log.InitDefaultLog("register", true, true)
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
		log.Error("ResourcePath is null")
		os.Exit(0)
	}
	if config.ResUsernamePath == "" {
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
	config2.UseCharles = true
	config2.UseTruncation = true

	initParams()
	routine.InitRoutine(config.ProxyPath)

	var err error
	PhoneProvider, err = phone.NewPhoneVerificationCode("taxin")
	if err != nil {
		log.Error("create phone provider error!%v", err)
		os.Exit(0)
	}
	err = common.InitResource(config.ResIcoPath, config.ResUsernamePath)
	if err != nil {
		log.Error("InitResource error!%v", err)
		os.Exit(0)
	}

	WaitAll.Add(config.Coro)
	for i := 0; i < config.Coro; i++ {
		go Register()
	}

	logTicker = time.NewTicker(time.Second * 10)
	go LogStatus()
	WaitAll.Wait()
	logTicker.Stop()
	log.Info("success: %d,challenge err: %d ,create err: %d, send msg err: %d, recv msg err: %d, other err: %d",
		SuccessCount,
		ErrorChallengeRequired,
		ErrorCreateCount,
		ErrorSendSMSCount,
		ErrorRecvSMSCount,
		ErrorOtherCount)
}
