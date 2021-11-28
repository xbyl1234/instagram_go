package tests

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	main.InitTest()
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
		searchResult, err := search.NextTags()
		if err != nil {
			log.Error("search next error: %v", err)
			break
		}
		log.Info("%v", searchResult)

		tags := searchResult.GetTags()
		for index := range tags {
			tag := tags[index]
			tag.Sync(goinsta.TabRecent)
			tag.Stories()
			tagResult, err := tag.Next()
			if err != nil {
				log.Error("next media error: %v", err)
				break
			}
			medias := tagResult.GetAllMedias()
			for mindex := range medias {
				media := medias[mindex]
				comments := media.GetComments()
				commResp, err := comments.NextComments()
				if err != nil {
					log.Error("next comm error: %v", err)
					break
				}

				for cindex := range commResp.Comments {
					log.Info("comment user id: %v", commResp.Comments[cindex].User.ID)
				}
			}
		}
	}
}
