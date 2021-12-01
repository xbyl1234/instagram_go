package routine

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"os"
)

func InitRoutine(proxyPath string) {
	common.InitMogoDB()
	err := common.InitProxyPool(proxyPath)
	if err != nil {
		log.Error("init ProxyPool error:%v", err)
		panic(err)
	}
	intas := goinsta.LoadAllAccount()
	if len(intas) == 0 {
		log.Error("there have no account!")
		os.Exit(0)
	}
	log.Info("load account count: %d", len(intas))
	goinsta.InitAccountPool(intas)
	//common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
}
