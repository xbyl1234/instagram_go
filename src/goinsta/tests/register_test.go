package tests

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/phone"
	"makemoney/goinsta"
	"makemoney/goinsta/dbhelper"
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
	err := common.InitProxyPool("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\zone2_ips_us.txt")
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
	common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
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
	_proxy := common.ProxyPool.GetOne()
	_proxy = common.ProxyPool.GetOne()
	_proxy = common.ProxyPool.GetOne()
	if _proxy == nil {
		log.Error("get proxy error: %v", _proxy)
	}
	regisert := goinsta.NewRegister(_proxy, provider)
	inst, err := regisert.Do("badrgirl", "badrgirl", "XBYLxbyl1234")
	if err != nil {
		log.Warn("register error, %v", err)
	} else {
		log.Info("register success, username %s, passwd %s", inst.User, inst.Pass)
		goinsta.SaveInstToDB(inst)
		err := inst.Account.ChangeProfilePicture(common.Resource.ChoiceIco())
		log.Info("ch ico %v", err)
	}
}
