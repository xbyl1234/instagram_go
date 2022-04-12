package main

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxys"
	"makemoney/goinsta"
	"sync/atomic"
	"time"
)

func RegisterByPhone() {
	provider := PhoneProvider
	for true {
		curCount := atomic.AddInt32(&Count, 1)
		if curCount > int32(*RegisterCount) {
			break
		}
		_proxy := proxys.ProxyPool.Get(config.Country, "")
		if _proxy == nil {
			log.Error("get proxy error: %v", _proxy)
			break
		}
		var err error

		username := common.Resource.ChoiceUsername()
		password := common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)

		inst := goinsta.New("", "", _proxy)
		regisert := goinsta.Register{
			Inst:         inst,
			RegisterType: "phone",
			//Account:      account,
			Username: username,
			Password: password,
			AreaCode: provider.GetArea(),
			Year:     fmt.Sprintf("%d", common.GenNumber(1995, 2000)),
			Month:    fmt.Sprintf("%02d", common.GenNumber(1, 11)),
			Day:      fmt.Sprintf("%02d", common.GenNumber(1, 27)),
		}
		for true {
			regisert.Account = "18501750803"
			_, err = regisert.SendSignupSmsCode()
			if err != nil {
				log.Error("phone %s send error: %v", "", err)
			}
			time.Sleep(10 * time.Second)
			regisert.Inst.ResetProxy()
		}

		inst.AccountInfo.Register.RegisterIpCountry = _proxy.Country
		prepare := inst.PrepareNewClient()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))

		err = regisert.GetSignupConfig()
		err = regisert.GetCommonEmailDomains()
		err = regisert.PrecheckCloudId()
		err = regisert.IgUser()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(2000, 3000)))

		var account string
		account, err = provider.RequireAccount()
		if err != nil {
			log.Error("require account error: %v", err)
			break
		}
		regisert.Account = account
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(1000, 2000)))

		_, err = regisert.SendSignupSmsCode()
		if err != nil {
			ErrorSendCodeCount++
			statError(err)
			provider.ReleaseAccount(regisert.Account)
			log.Error("phone %s send error: %v", account, err)
			continue
		}
		code, err := provider.RequireCode(account)
		if err != nil {
			ErrorRecvCodeCount++
			statError(err)
			log.Error("phone %s require code error: %v", account, err)
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.ValidateSignupSmsCode(code)
		if err != nil {
			ErrorCodeCount++
			statError(err)
			log.Error("phone %s check code error: %v", account, err)
			continue
		}

		regisert.GenUsername()
		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CheckAgeEligibility()
		_, err = regisert.NewUserFlowBegins()

		time.Sleep(time.Millisecond * time.Duration(common.GenNumber(0, 1000)))
		_, err = regisert.CreatePhone()
		if err != nil {
			ErrorCreateCount++
			statError(err)
			log.Error("phone %s create error: %v", account, err)
			continue
		}

		prepareStr, _ := json.Marshal(prepare)
		_, err = regisert.GetSteps()

		if err == nil {
			log.Info("phone: %s register success! prepare: %s, account: %s password: %s", account, prepareStr, inst.User, inst.Pass)
			_ = goinsta.SaveInstToDB(inst)
			_, err = regisert.NewAccountNuxSeen()
			_, err = inst.AddressBookLink(GenAddressBook())
			var uploadID string
			uploadID, _, err = inst.GetUpload().UploadPhotoFromPath(common.Resource.ChoiceIco(), nil)
			err = inst.GetAccount().ChangeProfilePicture(uploadID)
			SuccessCount++
		} else {
			_ = goinsta.SaveInstToDB(inst)
			accInfo, _ := json.Marshal(inst.AccountInfo)
			log.Error("phone %s create error: %v, prepare: %s, account info: %s", account, err, prepareStr, accInfo)
			statError(err)
			ErrorCreateCount++

			//if common.IsError(err, common.ChallengeRequiredError) {
			//	log.Error("phone: %s had been challenge_required", account)
			//	continue
			//} else if common.IsError(err, common.FeedbackError) {
			//	ErrorCreateCount++
			//	log.Error("phone: %s had been feedback_required", account)
			//	continue
			//}
		}

	}
	WaitAll.Done()
}
