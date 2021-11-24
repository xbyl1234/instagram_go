package tests

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/goinsta/dbhelper"
	"os"
	"testing"
)

func InitTestSearch() {
	log.InitLogger()
	dbhelper.InitMogoDB()
	err := common.InitProxyPool("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\zone2_ips_us.txt")
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
	common.InitResource("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture", "C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\user_nameraw.txt")
}

func TestSearch(t *testing.T) {
	InitTestSearch()
	inst := goinsta.AccountPool.GetOne()
	_proxy := common.ProxyPool.Get(inst.Proxy.ID)
	if _proxy == nil {
		log.Error("find insta proxy error!")
		os.Exit(0)
	}
	inst.SetProxy(_proxy)
	inst.Login()
	//inst.Account.Sync()
	//inst.Account.ChangeProfilePicture(common.Resource.ChoiceIco())

	search := inst.GetSearch("china")
	for true {
		tag, err := search.NextTags()
		if err != nil {
			log.Error("search next error: %v", err)
			break
		}
		log.Info("%v", tag)
	}
}
