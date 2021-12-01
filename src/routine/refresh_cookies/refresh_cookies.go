package main

import (
	"container/list"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/config"
	"makemoney/goinsta"
	"os"
)

type tmpAccount struct {
	username string
	passwd   string
}

func main() {
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

	accounts := list.New()
	accounts.PushBack(&tmpAccount{"badrgirl67", "XBYLxbyl1234"})
	accounts.PushBack(&tmpAccount{"badrgirl21", "XBYLxbyl1234"})
	accounts.PushBack(&tmpAccount{"badrgirl21", "XBYLxbyl1234"})
	//accounts.PushBack(&tmpAccount{"lovergirl5289", "XBYLxbyl1234"})
	//accounts.PushBack(&tmpAccount{"badrgirl5", "XBYLxbyl1234"})
	//accounts.PushBack(&tmpAccount{"badrgirl6", "XBYLxbyl1234"})
	//accounts.PushBack(&tmpAccount{"badrgirl67", "XBYLxbyl1234"})
	config.UseCharles = false
	main.InitTest()
	for item := accounts.Front(); item != nil; item = item.Next() {
		acc := item.Value.(*tmpAccount)
		inst := goinsta.New(acc.username, acc.passwd, common.ProxyPool.GetOne())
		inst.PrepareNewClient()
		err := inst.Login()
		if err != nil {
			log.Warn("username: %s, login error: %v", acc.username, err.Error())
		} else {
			log.Info("username: %s, login success", acc.username)
		}
		_ = goinsta.SaveInstToDB(inst)
	}
}
