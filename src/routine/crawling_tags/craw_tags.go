package main

import (
	"container/list"
	mapset "github.com/deckarep/golang-set"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"strings"
)

var TagList = list.New()
var TagIDSet = mapset.NewSet()
var KeyWord []string

var CrawTagsOperName = "craw_media"
var CrawTagsAccountTag = "craw_media"

func LoadKeyWord() error {
	data, err := os.ReadFile(config.KeyWordPath)
	if err != nil {
		log.Error("read key word file error: %v", err)
		return err
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
	return nil
}

func LoadTags() {
	retTagList, err := routine.LoadTags()
	if err != nil {
		log.Error("preCrawTags LoadTags error:%v", err)
		os.Exit(0)
	}

	for index := range retTagList {
		TagIDSet.Add(retTagList[index].Id)
		TagList.PushBack(&retTagList[index])
	}
}

func CheckTagsAndRunCrawTags() {
	search, err := routine.LoadSearch()
	if err != nil {
		log.Error("LoadTags error:%v", err)
		os.Exit(0)
	}
	err = LoadKeyWord()
	if err != nil {
		log.Error("LoadKeyWord error:%v", err)
		os.Exit(0)
	}

	for _, key := range KeyWord {
		find := false
		for _, searchItem := range search {
			if searchItem.Q == key {
				find = true
				break
			}
		}

		if !find {
			search = append(search, &goinsta.Search{
				Inst:    nil,
				Q:       key,
				HasMore: true,
			})
		}
	}

	for _, item := range search {
		if !item.HasMore {
			log.Info("pass search tag %s", item.Q)
			continue
		}
		CrawTags(item)
	}
}

func CrawTags(search *goinsta.Search) {
	var currAccount *goinsta.Instagram

	var RequireAccount = func(search *goinsta.Search) *goinsta.Search {
		inst := routine.ReqAccount(CrawTagsOperName, CrawTagsAccountTag)
		currAccount = inst
		if inst == nil {
			log.Error("CrawTags req account error")
			return nil
		}

		if search.Inst != nil {
			oldUser := search.Inst.User
			goinsta.AccountPool.ReleaseOne(search.Inst)
			search.SetAccount(inst)
			log.Warn("CrawTags replace account %s->%s", oldUser, inst.User)
		} else {
			search.SetAccount(inst)
			log.Info("CrawTags set account %s", inst.User)
		}
		return search
	}

	search = RequireAccount(search)
	defer func() {
		if currAccount != nil {
			goinsta.AccountPool.ReleaseOne(currAccount)
		}
	}()

	for true {
		searchResult, err := search.NextTags()
		if err != nil {
			if common.IsNoMoreError(err) {
				log.Info("tags %s has craw finish!", search.Q)
				break
			} else if common.IsError(err, common.ChallengeRequiredError) || common.IsError(err, common.LoginRequiredError) {
				search = RequireAccount(search)
				continue
			} else if common.IsError(err, common.RequestError) {
				log.Warn("CrawMedias retrying...user: %s, err: %v", search.Inst.User, err)
				continue
			} else {
				log.Error("search next unknow error: %v", err)
				continue
			}
		}

		tags := searchResult.GetTags()
		for index := range tags {
			if TagIDSet.Contains(tags[index].Id) {
				log.Info("")
				continue
			}
			TagIDSet.Add(tags[index].Id)
			TagList.PushBack(&tags[index])
			err = routine.SaveTags(&tags[index])
			if err != nil {
				log.Error("SaveTags error:%v", err)
			}
		}
		_ = routine.SaveSearch(search)
	}
}
