package common

import (
	"math/rand"
)

var AntiFlag float64 = 0.45841

func RandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func AntiJudge() bool {
	if AntiFlag > 0.5 {
		return false
	}
	return true
}

func InitAnti(wait chan int) {
	tarDevice := GetDeviceID()
	if DeviceID-tarDevice < 10 && DeviceID-tarDevice >= 0 {
		AntiFlag = RandFloats(0, 0.5)
	} else {
		AntiFlag = RandFloats(0.5, 1)
	}
	close(wait)
}
