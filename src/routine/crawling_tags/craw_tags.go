package main

import (
	mapset "github.com/deckarep/golang-set"
	"makemoney/common/log"
	"makemoney/goinsta"
)

func CrawTags(inst *goinsta.Instagram, keys []string) []*goinsta.Tags {
	defer func() {
		if err := recover(); err != nil {
			log.Error("account: %s CrawTags panic error: %v", inst.User, err)
		}
	}()
	var TagIDSet = mapset.NewSet()
	var ret = make([]*goinsta.Tags, 0)
	for _, key := range keys {
		search := inst.NewSearch(key)
		searchResult, err := search.NextTags()
		if err != nil {
			log.Error("search tag error: %v", err)
		}

		for _, item := range searchResult.Tags {
			if TagIDSet.Contains(item.Id) {
				continue
			}
			TagIDSet.Add(item.Id)
			ret = append(ret, item)
		}
	}
	return ret
}
