package main

import (
	mapset "github.com/deckarep/golang-set"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"os"
	"strings"
)

func LoadKeyWord() ([]string, error) {
	var KeyWord []string
	data, err := os.ReadFile(config.KeyWordPath)
	if err != nil {
		log.Error("read key word file error: %v", err)
		return nil, err
	}
	sp := strings.Split(common.B2s(data), "\n")
	KeyWord = make([]string, len(sp))
	index := 0
	for _, item := range sp {
		item = strings.Trim(item, " ")
		item = strings.ReplaceAll(item, "\r", "")
		item = strings.ReplaceAll(item, "\n", "")
		if item != "" {
			KeyWord[index] = item
			index++
		}
	}
	KeyWord = KeyWord[:index]
	return KeyWord, err
}

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
