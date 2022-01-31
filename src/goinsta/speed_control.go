package goinsta

import (
	"makemoney/common"
	"sync"
	"time"
)

type SpeedControlConfig struct {
	OperName   string `json:"oper_name"`
	EachSecond int    `json:"each_second"`
	EachMinute int    `json:"each_minute"`
	EachHour   int    `json:"each_hour"`
	EachDay    int    `json:"each_day"`
}

type SpeedControlJson struct {
	SpeedControl []SpeedControlConfig `json:"speed_control"`
}

var SpeedControlConfigMap map[string]SpeedControlConfig

func InitSpeedControl(path string) error {
	var config SpeedControlJson
	err := common.LoadJsonFile(path, &config)
	if err != nil || len(config.SpeedControl) == 0 {
		return err
	}
	SpeedControlConfigMap = make(map[string]SpeedControlConfig)
	for _, item := range config.SpeedControl {
		SpeedControlConfigMap[item.OperName] = item
	}
	return nil
}

type LimitRate struct {
	Rate     int           `db:"-"`
	Begin    time.Time     `db:"begin"`
	Count    int           `db:"count"`
	Interval time.Duration `db:"interval"`
	Lock     sync.Mutex    `db:"-"`
}

func (this *LimitRate) Limit(block bool) bool {
	result := true
	this.Lock.Lock()
	if this.Count == this.Rate {
		if time.Now().Sub(this.Begin) >= this.Interval {
			this.Begin = time.Now()
			this.Count = 0
		} else {
			result = false
		}
	} else {
		this.Count++
	}
	this.Lock.Unlock()
	return result
}

type SpeedControl struct {
	EachSecond LimitRate `db:"each_second"`
	EachMinute LimitRate `db:"each_minute"`
	EachHour   LimitRate `db:"each_hour"`
	EachDay    LimitRate `db:"each_day"`
}

func (this *SpeedControl) TestSpeed() bool {
	if this.EachSecond.Rate != 0 {
		if this.EachSecond.Limit(true) {
			return false
		}
	}
	if this.EachMinute.Rate != 0 {
		if this.EachMinute.Limit(false) {
			return false
		}
	}
	if this.EachHour.Rate != 0 {
		if this.EachHour.Limit(false) {
			return false
		}
	}
	if this.EachDay.Rate != 0 {
		if this.EachDay.Limit(false) {
			return false
		}
	}
	return true
}

func GetSpeedControl(OperName string) *SpeedControl {
	config := SpeedControlConfigMap[OperName]
	ret := &SpeedControl{
		EachSecond: LimitRate{
			Rate:     config.EachSecond,
			Begin:    time.Now(),
			Interval: time.Second,
		},
		EachMinute: LimitRate{
			Rate:     config.EachMinute,
			Begin:    time.Now(),
			Interval: time.Minute,
		},
		EachHour: LimitRate{
			Rate:     config.EachHour,
			Begin:    time.Now(),
			Interval: time.Hour,
		},
		EachDay: LimitRate{
			Rate:     config.EachDay,
			Begin:    time.Now(),
			Interval: time.Hour * 24,
		},
	}
	return ret
}

func ReSetRate(sp *SpeedControl, OperName string) {
	config := SpeedControlConfigMap[OperName]
	sp.EachSecond.Rate = config.EachSecond
	sp.EachMinute.Rate = config.EachMinute
	sp.EachHour.Rate = config.EachHour
	sp.EachDay.Rate = config.EachDay
}
