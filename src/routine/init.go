package routine

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"makemoney/goinsta"
	math_rand "math/rand"
	"time"
)

type DataBaseConfig struct {
	MogoUri string `json:"mogo_uri"`
}

var dbConfig DataBaseConfig

func InitRoutine(proxyPath string) {
	math_rand.Seed(time.Now().UnixNano())
	err := common.LoadJsonFile("./config/dbconfig.json", &dbConfig)
	if err != nil {
		log.Error("load db config error:%v", err)
		panic(err)
	}
	goinsta.InitMogoDB(dbConfig.MogoUri)
	err = proxy.InitProxyPool(proxyPath)
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
	err = goinsta.InitInstagramConst()
	if err != nil {
		log.Error("load const json error:%v", err)
		panic(err)
	}
	goinsta.ProxyCallBack = ProxyCallBack
	goinsta.InitGraph()
	//common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
}

func ReqAccount(block bool) (*goinsta.Instagram, error) {
	var reqAccount = func() (*goinsta.Instagram, error) {
		inst := goinsta.AccountPool.GetOne(block)
		if inst == nil {
			return nil, &common.MakeMoneyError{ErrType: common.NoMoreError, ErrStr: "no more account"}
		}
		if !SetProxy(inst) {
			return nil, &common.MakeMoneyError{ErrType: common.NoMoreError, ErrStr: "no more proxy"}
		}

		return inst, nil
	}

	for true {
		inst, err := reqAccount()
		if err != nil {
			log.Warn("try require account error: %v", err)
			if block {
				time.Sleep(time.Second * 5)
			} else {
				return nil, err
			}
		} else {
			log.Info("require account: %v", inst.User)
			return inst, nil
		}
	}
	return nil, nil
}

func SetProxy(inst *goinsta.Instagram) bool {
	var _proxy *proxy.Proxy
	if inst.Proxy != nil {
		if inst.Proxy.ID != "" {
			_proxy = proxy.ProxyPool.Get(inst.RegisterIpCountry, inst.Proxy.ID)
			if _proxy == nil {
				log.Warn("find insta proxy %s error!", inst.Proxy.ID)
			}
		}
	}

	if _proxy == nil {
		_proxy = proxy.ProxyPool.GetNoRisk(inst.RegisterIpCountry, false, false)
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

func ProxyCallBack(country string, id string) (*proxy.Proxy, error) {
	var _proxy *proxy.Proxy
	if country != "" {
		_proxy = proxy.ProxyPool.Get(country, id)
		if _proxy == nil {
			log.Warn("find insta proxy %s error!", country)
		}
	}

	if _proxy == nil {
		_proxy = proxy.ProxyPool.GetNoRisk(country, false, false)
	}

	if _proxy == nil {
		log.Error("get insta proxy error!")
		return nil, &common.MakeMoneyError{
			ErrStr:  "no more proxy",
			ErrType: common.PorxyError,
		}
	}

	return _proxy, nil
}
