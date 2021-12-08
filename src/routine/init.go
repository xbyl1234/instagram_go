package routine

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"os"
)

type DataBaseConfig struct {
	MogoUri string `json:"mogo_uri"`
}

var dbConfig DataBaseConfig

func InitRoutine(proxyPath string) {
	err := common.LoadJsonFile("./config/dbconfig.json", &dbConfig)
	if err != nil {
		log.Error("load db config error:%v", err)
		panic(err)
	}
	common.InitMogoDB(dbConfig.MogoUri)
	err = common.InitProxyPool(proxyPath)
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
	intas := goinsta.LoadAllAccount()
	if len(intas) == 0 {
		log.Error("there have no account!")
		os.Exit(0)
	}
	goinsta.InitAccountPool(intas)

	log.Info("load account count: %d", goinsta.AccountPool.Available.Len())
	//common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
}

func ReqAccount() *goinsta.Instagram {
	inst := goinsta.AccountPool.GetOne()
	if inst == nil {
		return nil
	}
	SetProxy(inst)
	return inst
}

func SetProxy(inst *goinsta.Instagram) bool {
	var _proxy *common.Proxy
	if inst.Proxy.ID != "" {
		_proxy = common.ProxyPool.Get(inst.Proxy.ID)
		if _proxy == nil {
			log.Warn("find insta proxy %s error!", inst.Proxy.ID)
		}
	}

	if _proxy == nil {
		_proxy = common.ProxyPool.GetNoRisk(false, false)
		if _proxy == nil {
			log.Error("get insta proxy error!")
		}
	}

	if _proxy != nil {
		inst.SetProxy(_proxy)
	} else {
		return false
	}
	return true
}
