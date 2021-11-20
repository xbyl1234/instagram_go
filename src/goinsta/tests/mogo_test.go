package tests

import (
	"makemoney/goinsta"
	"makemoney/proxy"
	"testing"
)
import "makemoney/goinsta/dbhelper"

func TestPhoneDb(t *testing.T) {
	dbhelper.InitMogoDB()
	ins := goinsta.New("as", "as", &proxy.Proxy{})
	goinsta.SaveInstToDB(ins)
}
