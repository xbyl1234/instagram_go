package main

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxys"
	"makemoney/common/verification"
	"makemoney/goinsta"
	"sync/atomic"
	"time"
)

func RegisterByEmail() {
	var mail verification.VerificationCodeProvider
	if config.ProviderName == "gmail" {
		mail = verification.GetGMails()
	} else {
		mail = Guerrilla
	}

	if mail == nil {
		log.Error("get mail error,so return!")
		return
	}

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

		account, err := mail.RequireAccount()
		//account := "admin1@followmebsix.com"
		if err != nil {
			log.Error("require account error: %v", err)
			break
		}

		username := common.Resource.ChoiceUsername()
		password := common.GenString(common.CharSet_abc, 4) +
			common.GenString(common.CharSet_123, 4)
		//password := "xbyl1234"

		inst := goinsta.New(username, password, _proxy)
		regisert := goinsta.Register{
			Inst:         inst,
			RegisterType: "email",
			Account:      account,
			Username:     username,
			Password:     password,
			Year:         fmt.Sprintf("%02d", common.GenNumber(1995, 2000)),
			Month:        fmt.Sprintf("%02d", common.GenNumber(1, 11)),
			Day:          fmt.Sprintf("%02d", common.GenNumber(1, 27)),
		}
		inst.AccountInfo.Register.RegisterIpCountry = _proxy.Country
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
		code, err := mail.RequireCode(account)
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
		_, err = regisert.GetSteps()

		if err == nil {
			log.Info("email: %s register success!   account: %s password: %s", account, inst.User, inst.Pass)
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
			log.Error("email %s create error: %v,   account info: %s", account, err, accInfo)
			statError(err)
			ErrorCreateCount++
		}
	}
	WaitAll.Done()
}
