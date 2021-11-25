package tests

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	InitTest()
	inst := goinsta.AccountPool.GetOne()
	_proxy := common.ProxyPool.Get(inst.Proxy.ID)
	if _proxy == nil {
		log.Error("find insta proxy error!")
		os.Exit(0)
	}
	inst.SetProxy(_proxy)

	//inst.Account.Sync()
	//err := inst.Account.ChangeProfilePicture(common.Resource.ChoiceIco())
	//if err != nil {
	//	log.Warn("%v", err)
	//}

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
