package goinsta

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"strings"
	"time"
)

const clipsTab = "clips_tab"
const containerModule = "clips_viewer_clips_tab"
const pctReels = "0"

type Feed struct {
	Inst        *Instagram
	LastMedias  *VideosFeedResp
	lastReqTime time.Time

	TabType         string
	ContainerModule string
	PctReels        string
	SessionInfo     string
	SessionId       string
	SeenReels       string
	MaxId           string
	MoreAvailable   bool
}

type NumLoop struct {
	Value         int `json:"value"`
	LastLoopEndTs int `json:"last_loop_end_ts"`
}
type TotalWatchTime struct {
	Value           int     `json:"value"`
	LatestPlayEndTs float64 `json:"latest_play_end_ts"`
}

type SeenInfoItem struct {
	NumLoops         NumLoop        `json:"num_loops"`
	TotalWatchTimeMs TotalWatchTime `json:"total_watch_time_ms"`
}

type SeenInfo struct {
	Items map[string]*SeenInfoItem
}

func newFeed(inst *Instagram) *Feed {
	return &Feed{
		Inst:            inst,
		MoreAvailable:   true,
		TabType:         clipsTab,
		ContainerModule: containerModule,
		PctReels:        pctReels,
		SessionId:       fmt.Sprintf("%d_%s", inst.ID, strings.ToUpper(common.GenUUID())),
	}
}

func (this *Feed) Next() (*VideosFeedResp, error) {
	if this.MoreAvailable {
		return nil, &common.MakeMoneyError{ErrStr: "no more", ErrType: common.NoMoreError}
	}

	params := map[string]interface{}{
		"tab_type":         this.TabType,
		"session_id":       this.SessionId,
		"_uuid":            this.Inst.AccountInfo.Device.DeviceID,
		"container_module": this.ContainerModule,
		"pct_reels":        this.PctReels,
	}

	if this.MaxId != "" {
		params["max_id"] = this.MaxId
		type SeenReel struct {
			Id string `json:"id"`
		}
		var seenReels []SeenReel
		var seenInfo SeenInfo
		used := 0
		for _, item := range this.LastMedias.Items {
			//item.Media.Pk
			id := fmt.Sprintf("%d", item.Media.Pk)
			seenReels = append(seenReels, SeenReel{Id: id})
			used += common.GenNumber(100, 2000)
			seenInfo.Items[id] = &SeenInfoItem{
				NumLoops: NumLoop{
					Value:         0,
					LastLoopEndTs: 0,
				},
				TotalWatchTimeMs: TotalWatchTime{
					Value:           used,
					LatestPlayEndTs: float64(this.lastReqTime.Add(time.Duration(used) * time.Millisecond).Unix()),
				},
			}
		}
		marshal, err := json.Marshal(seenReels)
		if err == nil {
			params["seen_reels"] = common.B2s(marshal)
		}
		marshal, err = json.Marshal(seenInfo.Items)
		if err == nil {
			params["session_info"] = common.B2s(marshal)
		}
	}

	ret := &VideosFeedResp{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		IsPost:         false,
		ApiPath:        urlDiscoverVideosFeed,
		HeaderSequence: LoginHeaderMap[urlDiscoverVideosFeed],
		Query:          params,
	}, ret)

	this.lastReqTime = time.Now()
	err = ret.CheckError(err)
	if err == nil {
		this.MaxId = ret.PagingInfo.MaxId
		this.MoreAvailable = ret.PagingInfo.MoreAvailable
	}
	return ret, err
}
