package main

import (
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
)

func NewTags(tag string, count int) {
	insts := goinsta.LoadAllAccount()
	goinsta.InitAccountPool(insts)
	for i := 0; i < count; {
		inst := goinsta.AccountPool.GetOneNoWait("", "")
		if inst == nil {
			log.Error("req account error!")
			break
		}

		i++
		inst.Tags = tag
		goinsta.SaveInstToDB(inst)
		log.Info("set %s", inst.User)
	}

	log.Info("NewTags finish!")
}

func main() {
	log.InitDefaultLog("tools", true, false)
	routine.InitRoutine("./config/proxy_config_ios.json")
	NewTags("craw_media", 10)
	NewTags("craw_comment", 10)
}
