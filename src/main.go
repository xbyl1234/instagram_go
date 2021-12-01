package main

import "makemoney/common/log"

func main() {
	log.InitDefaultLog("test", true, true)
	log.Info("sadasdadadasd")
	log.Warn("sadasdadadasd")
	log.Error("sadasdadadasd")
}
