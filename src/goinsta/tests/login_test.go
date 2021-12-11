package tests

import (
	"container/list"
	"makemoney/common/log"
	"makemoney/common/proxy"
	"makemoney/config"
	"makemoney/goinsta"
	"testing"
)

type tmpAccount struct {
	username string
	passwd   string
}

func TestLogin(t *testing.T) {
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
		inst := goinsta.New(acc.username, acc.passwd, proxy.ProxyPool.GetNoRisk())
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
