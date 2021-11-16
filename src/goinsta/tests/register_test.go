package tests

import (
	"makemoney/goinsta"
	"makemoney/goinsta/dbhelper"
	"makemoney/log"
	"makemoney/phone"
	"os"
	"path/filepath"
	"testing"
)

func SetCurrPath() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(dir)
}
func InitAll() {
	log.InitLogger()
	dbhelper.InitMogoDB()
}

func TestRegister(t *testing.T) {
	SetCurrPath()
	InitAll()
	curPath, _ := os.Getwd()
	log.Info("cur path %s", curPath)

	provider, err := phone.NewPhoneVerificationCode("do889")
	if err != nil {
		log.Error("create phone provider error!")
		os.Exit(0)
	}
	//err = provider.Login()
	//if err != nil {
	//	log.Error("provider login error!")
	//	os.Exit(0)
	//}

	regisert := goinsta.NewRegister("+86", provider)
	inst, err := regisert.Do("lovergirl", "lovergirl", "XBYLxbyl1234")
	if err != nil {
		log.Warn("register error, %v", err)
	} else {
		log.Info("register success, username %s, passwd %s", inst.User, inst.Pass)
	}
}
