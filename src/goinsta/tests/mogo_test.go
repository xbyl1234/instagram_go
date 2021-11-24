package tests

import (
	"makemoney/common"
	"makemoney/goinsta"
	"testing"
)
import "makemoney/goinsta/dbhelper"

func TestPhoneDb(t *testing.T) {
	dbhelper.InitMogoDB()
	ins := goinsta.New("as", "as", &common.Proxy{})
	goinsta.SaveInstToDB(ins)
}
