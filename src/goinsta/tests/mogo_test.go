package tests

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	math_rand "math/rand"
	"testing"
	"time"
)

type DataBaseConfig struct {
	MogoUri string `json:"mogo_uri"`
}

var dbConfig DataBaseConfig

func TestAccount(t *testing.T) {
	math_rand.Seed(time.Now().UnixNano())
	log.InitDefaultLog("register", true, true)
	err := common.LoadJsonFile("./config/dbconfig.json", &dbConfig)
	if err != nil {
		log.Error("load db config error:%v", err)
		panic(err)
	}
	goinsta.InitMogoDB(dbConfig.MogoUri)

	goinsta.InitInstagramConst()
	instas := goinsta.LoadAllAccount()
	for _, item := range instas {
		goinsta.SaveInstToDB(item)
	}
}
