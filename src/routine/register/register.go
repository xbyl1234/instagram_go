package main

import (
	"flag"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"makemoney/common/verification"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	ProxyPath       string                   `json:"proxy_path"`
	ResIcoPath      string                   `json:"res_ico_path"`
	ResUsernamePath string                   `json:"res_username_path"`
	Coro            int                      `json:"coro"`
	Country         string                   `json:"country"`
	ProviderName    string                   `json:"provider_name"`
	Provider        []*verification.Provider `json:"provider"`
}

var ConfigPath = flag.String("config", "./config/register.json", "")
var RegisterCount = flag.Int("count", 0, "")

var config Config

var Count int32 = 0
var SuccessCount = 0
var ErrorCreateCount = 0
var ErrorSendCodeCount = 0
var ErrorRecvCodeCount = 0
var ErrorCodeCount = 0
var ErrorCheckAccountCount = 0

var ErrorChallengeRequired = 0
var ErrorFeedback = 0
var ErrorOther = 0

var WaitAll sync.WaitGroup

var logTicker *time.Ticker

func LogStatus() {
	for range logTicker.C {
		log.Info("success: %d, create err: %d, send err: %d, recv err: %d, challenge: %d, feedback: %d, check err: %d",
			SuccessCount,
			ErrorCreateCount,
			ErrorSendCodeCount,
			ErrorRecvCodeCount,
			ErrorChallengeRequired,
			ErrorFeedback,
			ErrorCheckAccountCount,
		)
	}
}

func statError(err error) {
	if common.IsError(err, common.ChallengeRequiredError) {
		ErrorChallengeRequired++
	} else if common.IsError(err, common.FeedbackError) {
		ErrorFeedback++
	} else {
		ErrorOther++
	}
}

func RegisterByPhone() {

}

func RegisterByEmail() {
	provider := verification.VerificationProvider[config.ProviderName]

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

		account, err := provider.RequireAccount()
		if err != nil {
			log.Error("require account error: %v", err)
			break
		}

		username := common.Resource.ChoiceUsername()
		password := common.GenString(common.CharSet_ABC, 4) +
			common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)
		inst := goinsta.New("", "", _proxy)
		regisert := goinsta.Register{
			Inst:     inst,
			Account:  account,
			Username: username,
			Password: password,
			Year:     fmt.Sprintf("%d", common.GenNumber(1995, 2000)),
			Month:    fmt.Sprintf("%d", common.GenNumber(1, 11)),
			Day:      fmt.Sprintf("%d", common.GenNumber(1, 27)),
		}

		inst.PrepareNewClient()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))
		err = regisert.GetSignupConfig()

		err = regisert.GetCommonEmailDomains()
		err = regisert.PrecheckCloudId()
		err = regisert.IgUser()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))
		_, err = regisert.CheckEmail()
		if err != nil {
			ErrorCheckAccountCount++
			statError(err)
			log.Error("email %s check error: %v", account, err)
			continue
		}

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(1000, 2000)))
		_, err = regisert.SendVerifyEmail()
		if err != nil {
			ErrorSendCodeCount++
			statError(err)
			log.Error("email %s send error: %v", account, err)
			continue
		}
		code, err := provider.RequireCode(account)
		if err != nil {
			ErrorRecvCodeCount++
			statError(err)
			log.Error("email %s require code error: %v", account, err)
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckConfirmationCode(code)
		if err != nil {
			ErrorCodeCount++
			statError(err)
			log.Error("email %s check code error: %v", account, err)
			continue
		}

		regisert.GenUsername()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckAgeEligibility()
		_, err = regisert.NewUserFlowBegins()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CreateEmail()
		if err != nil {
			ErrorCreateCount++
			statError(err)
			log.Error("email %s create error: %v", account, err)
			continue
		}

		_, err = regisert.NewAccountNuxSeen()
		_, err = regisert.GetSteps()

		err = goinsta.SaveInstToDB(inst)

		var uploadID string
		uploadID, err = inst.GetUpload().RuploadPhotoFromPath(common.Resource.ChoiceIco())
		err = inst.GetAccount().ChangeProfilePicture(uploadID)

		if err != nil {
			statError(err)
			if common.IsError(err, common.ChallengeRequiredError) {
				log.Error("email: %s had been challenge_required", account)
				ErrorCreateCount++
				continue
			} else if common.IsError(err, common.FeedbackError) {
				ErrorCreateCount++
				log.Error("email: %s had been feedback_required", account)
				continue
			}

			log.Warn("email: %s change ico error: %v", account, err)
		}

		SuccessCount++
		log.Info("email: %s register success!", account)
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
	config2.UseCharles = false
	config2.UseTruncation = true

	initParams()
	routine.InitRoutine(config.ProxyPath)

	err := verification.InitVerificationProviderByJson(config.Provider)
	if err != nil {
		log.Error("create phone provider error!%v", err)
		os.Exit(0)
	}
	if verification.VerificationProvider[config.ProviderName] == nil {
		log.Error("no find phone provider %s", config.ProviderName)
		os.Exit(0)
	}

	err = common.InitResource(config.ResIcoPath, config.ResUsernamePath)
	if err != nil {
		log.Error("InitResource error!%v", err)
		os.Exit(0)
	}

	WaitAll.Add(config.Coro)

	if config.ProviderName == "guerrilla" {
		for i := 0; i < config.Coro; i++ {
			go RegisterByEmail()
		}
	} else if config.ProviderName == "taxin" {
		for i := 0; i < config.Coro; i++ {
			go RegisterByPhone()
		}
	}

	logTicker = time.NewTicker(time.Second * 10)
	go LogStatus()
	WaitAll.Wait()
	logTicker.Stop()
	log.Info("success: %d, create err: %d, send err: %d, recv err: %d, challenge: %d, feedback: %d, check err: %d",
		SuccessCount,
		ErrorCreateCount,
		ErrorSendCodeCount,
		ErrorRecvCodeCount,
		ErrorChallengeRequired,
		ErrorFeedback,
		ErrorCheckAccountCount)
	log.Info("task finish!")
}
