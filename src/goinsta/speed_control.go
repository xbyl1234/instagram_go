package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
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

var SpeedControlConfigMap map[string]*SpeedControlConfig

func InitSpeedControl(path string) error {
	var config SpeedControlJson
	err := common.LoadJsonFile(path, &config)
	if err != nil || len(config.SpeedControl) == 0 {
		return err
	}
	SpeedControlConfigMap = make(map[string]*SpeedControlConfig)
	for _, item := range config.SpeedControl {
		SpeedControlConfigMap[item.OperName] = &item
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

func (this *LimitRate) Increase() int {
	this.Count++
	return this.Count
}

func (this *LimitRate) IsLimit(block bool) bool {
	result := false
	this.Lock.Lock()
	if this.Count >= this.Rate {
		if time.Now().Sub(this.Begin) >= this.Interval {
			this.Begin = time.Now()
			this.Count = 0
		} else {
			result = true
		}
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

func (this *SpeedControl) Increase() (int, int, int, int) {
	return this.EachSecond.Increase(),
		this.EachMinute.Increase(),
		this.EachHour.Increase(),
		this.EachDay.Increase()

}

func (this *SpeedControl) GetSpeed() (int, int, int, int) {
	return this.EachSecond.Count,
		this.EachMinute.Count,
		this.EachHour.Count,
		this.EachDay.Count
}

func (this *SpeedControl) IsSpeedLimit() bool {
	var ret = false
	if this.EachSecond.Rate != 0 {
		if this.EachSecond.IsLimit(true) {
			ret = true
		}
	}
	if this.EachMinute.Rate != 0 {
		if this.EachMinute.IsLimit(false) {
			ret = true
		}
	}
	if this.EachHour.Rate != 0 {
		if this.EachHour.IsLimit(false) {
			ret = true
		}
	}
	if this.EachDay.Rate != 0 {
		if this.EachDay.IsLimit(false) {
			ret = true
		}
	}
	return ret
}

func GetSpeedControl(OperName string) *SpeedControl {
	config := SpeedControlConfigMap[OperName]
	if config == nil {
		log.Warn("not find %s speed control,set no limit", OperName)
		config = &SpeedControlConfig{
			OperName:   OperName,
			EachSecond: 0,
			EachMinute: 0,
			EachHour:   0,
			EachDay:    0,
		}
	}

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
	if config == nil {
		log.Warn("not find %s speed control,set no limit", OperName)
		config = &SpeedControlConfig{
			OperName:   OperName,
			EachSecond: 0,
			EachMinute: 0,
			EachHour:   0,
			EachDay:    0,
		}
	}
	sp.EachSecond.Rate = config.EachSecond
	sp.EachMinute.Rate = config.EachMinute
	sp.EachHour.Rate = config.EachHour
	sp.EachDay.Rate = config.EachDay

	sp.EachSecond.IsLimit(false)
	sp.EachMinute.IsLimit(false)
	sp.EachHour.IsLimit(false)
	sp.EachDay.IsLimit(false)
}
