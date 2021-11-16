package tests

import "testing"
import "makemoney/goinsta/dbhelper"

func TestPhoneDb(t *testing.T) {
	dbhelper.InitMogoDB()

}
