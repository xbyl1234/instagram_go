package tests

import (
	"makemoney/goinsta"
	"makemoney/goinsta/dbhelper"
	"makemoney/log"
	"makemoney/phone"
	"makemoney/proxy"
	"os"
	"path/filepath"
	"testing"
)

func SetCurrPath() {
	dir, _ := filepath.Abs("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project")
	os.Chdir(dir)
}
func InitAll() {
	log.InitLogger()
	dbhelper.InitMogoDB()
	err := proxy.InitProxyPool("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\zone2_ips_us.txt")
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
}

func TestRegister(t *testing.T) {
	SetCurrPath()
	InitAll()
	curPath, _ := os.Getwd()
	log.Info("cur path %s", curPath)

	provider, err := phone.NewPhoneVerificationCode("do889")
	if err != nil {
		log.Error("create phone provider error!%v", err)
		os.Exit(0)
	}
	//err = provider.Login()
	//if err != nil {
	//	log.Error("provider login error!")
	//	os.Exit(0)
	//}
	_proxy, err := proxy.ProxyPool.GetOne()
	_proxy, err = proxy.ProxyPool.GetOne()
	_proxy, err = proxy.ProxyPool.GetOne()
	if err != nil {
		log.Error("get proxy error: %v", _proxy)
	}
	regisert := goinsta.NewRegister(_proxy, provider)
	inst, err := regisert.Do("badrgirl", "badrgirl", "XBYLxbyl1234")
	if err != nil {
		log.Warn("register error, %v", err)
	} else {
		log.Info("register success, username %s, passwd %s", inst.User, inst.Pass)
		goinsta.SaveInstToDB(inst)
	}
}
