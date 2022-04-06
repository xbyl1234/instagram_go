package main

import (
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/verification"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"time"
)

type Config struct {
	ProxyPath       string                        `json:"proxy_path"`
	ResIcoPath      string                        `json:"res_ico_path"`
	ResUsernamePath string                        `json:"res_username_path"`
	Coro            int                           `json:"coro"`
	Country         string                        `json:"country"`
	ProviderName    string                        `json:"provider_name"`
	Gmail           *verification.GmailConfig     `json:"gmail"`
	Guerrilla       *verification.GuerrillaConfig `json:"guerrilla"`
	Taxin           *verification.PhoneInfo       `json:"taxin"`
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

var PhoneProvider verification.VerificationCodeProvider
var Guerrilla verification.VerificationCodeProvider

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

func GenAddressBook() []goinsta.AddressBook {
	addr := make([]goinsta.AddressBook, common.GenNumber(20, 30))
	for index := range addr {
		addr[index].EmailAddresses = []string{common.GenString(common.CharSet_All, common.GenNumber(0, 10)) + "@gmail.com"}
		addr[index].PhoneNumbers = []string{"+1 " + "410 " + "895 " + common.GenString(common.CharSet_123, 4)}
		addr[index].LastName = common.GenString(common.CharSet_All, common.GenNumber(0, 10))
		addr[index].FirstName = common.GenString(common.CharSet_All, common.GenNumber(0, 10))
	}
	return addr[:]
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
	goinsta.UsePanic = false
	common.UseCharles = false
	common.UseTruncation = true
	initParams()
	routine.InitRoutine(config.ProxyPath)
	var err error
	switch config.ProviderName {
	case "taxin":
		PhoneProvider, err = verification.InitTaxin(config.Taxin)
		break
	case "gmail":
		err = verification.InitDefaultGMail(config.Gmail)
		break
	case "guerrilla":
		Guerrilla, err = verification.InitGuerrilla(config.Guerrilla)
		break
	}

	if err != nil {
		log.Error("create provider error! %v", err)
		os.Exit(0)
	}

	err = common.InitResource(config.ResIcoPath, config.ResUsernamePath)
	if err != nil {
		log.Error("InitResource error!%v", err)
		os.Exit(0)
	}

	WaitAll.Add(config.Coro)

	if config.ProviderName == "gmail" || config.ProviderName == "guerrilla" {
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
