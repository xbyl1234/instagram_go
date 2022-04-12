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
	SpeedControl []*SpeedControlConfig `json:"speed_control"`
}

const (
	OperNameCrawMedia   = "craw_media"
	OperNameCrawComment = "craw_comment"
	OperNameSendMsg     = "send_msg"
	OperNameLikeUser    = "like_user"
	OperNameComment     = "comment"
	OperNamePostImg     = "post_img"
	OperNamePostVideo   = "post_video"
	OperNameFeedVideo   = "feed_video"
)

var SpeedControlConfigMap map[string]*SpeedControlConfig

func InitSpeedControl(path string) error {
	var config SpeedControlJson
	err := common.LoadJsonFile(path, &config)
	if err != nil || len(config.SpeedControl) == 0 {
		return err
	}
	SpeedControlConfigMap = make(map[string]*SpeedControlConfig)
	for _, item := range config.SpeedControl {
		SpeedControlConfigMap[item.OperName] = item
	}
	return nil
}

type LimitRate struct {
	Rate     int           `bson:"-"`
	Begin    time.Time     `bson:"begin"`
	Count    int           `bson:"count"`
	Interval time.Duration `bson:"interval"`
	Lock     sync.Mutex    `bson:"-"`
}

func (this *LimitRate) Increase() int {
	this.Count++
	return this.Count
}

func (this *LimitRate) GetCool() time.Duration {
	if this.Count < this.Rate {
		return -1
	}

	return this.Interval - time.Now().Sub(this.Begin)
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
	EachSecond     LimitRate `bson:"each_second"`
	EachMinute     LimitRate `bson:"each_minute"`
	EachHour       LimitRate `bson:"each_hour"`
	EachDay        LimitRate `bson:"each_day"`
	History        int       `bson:"history"`
	SuccessHistory int       `bson:"success_history"`
}

func (this *SpeedControl) IncreaseSuccess() {
	this.SuccessHistory++
}

func (this *SpeedControl) Increase() (int, int, int, int) {
	this.History++
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

func (this *SpeedControl) GetCoolTime() time.Duration {
	if this.EachDay.GetCool() > 0 {
		return this.EachDay.GetCool()
	}
	if this.EachHour.GetCool() > 0 {
		return this.EachHour.GetCool()
	}
	if this.EachMinute.GetCool() > 0 {
		return this.EachMinute.GetCool()
	}
	if this.EachSecond.GetCool() > 0 {
		return this.EachSecond.GetCool()
	}
	return -1
}

func (this *SpeedControl) IsSpeedLimitInDay() bool {
	if this.EachDay.Rate != 0 {
		if this.EachDay.IsLimit(false) {
			return true
		}
	}
	return false
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

func GetSpeedControl(OperName string) (*SpeedControl, bool) {
	var isCtrl = true
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
		isCtrl = false
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
	return ret, isCtrl
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
