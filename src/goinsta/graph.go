package goinsta

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"time"
)

type Graph struct {
	inst *Instagram
	seq  int
}

func (this *Graph) SendBeforeSendSMS() {
	err := this.SendRequest(beforeSendSMS1, nil)
	if err != nil {
		log.Warn("graph send beforeSendSMS1 err: %v", err)
	}
	err = this.SendRequest(beforeSendSMS2, nil)
	if err != nil {
		log.Warn("graph send beforeSendSMS2 err: %v", err)
	}
	err = this.SendRequest(beforeSendSMS3, nil)
	if err != nil {
		log.Warn("graph send beforeSendSMS3 err: %v", err)
	}
}
func (this *Graph) SendAfterSendSMS() {
	err := this.SendRequest(afterSendSMS, nil)
	if err != nil {
		log.Warn("graph send afterSendSMS err: %v", err)
	}
}

func (this *Graph) SendRequest(action []string, params []map[string]interface{}) error {
	bmsg, err := json.Marshal(this.makeRequest(action, params))
	if err != nil {
		log.Error("Marshal GraphMsg error: %v", err)
		return err
	}

	resp := &GraphResp{}
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath:    urlLoggingClientEvents,
		IsPost:     true,
		IsApiGraph: true,
		Signed:     false,
		Query: map[string]interface{}{
			"format":        "json",
			"message":       common.Base64Encode(common.GZipCompress(bmsg)),
			"system_uptime": "",
			"compressed":    "1",
			"access_token":  InstagramAccessToken,
			"sent_time":     fmt.Sprintf("%d", time.Now().Unix()) + ".00000",
		},
		Header: map[string]string{
			"User-Agent":                this.inst.UserAgent,
			"X-Ig-Bandwidth-Speed-Kbps": "0.000",
			"Accept-Language":           "en-US;q=1.0",
			"Content-Type":              "application/x-www-form-urlencoded; charset=UTF-8",
			"X-Tigon-Is-Retry":          "False",
			"Accept-Encoding":           "gzip, deflate",
			"X-Fb-Http-Engine":          "Liger",
			"X-Fb-Client-Ip":            "True",
			"X-Fb-Server-Cluster":       "True",
		},
		DisAutoHeader: true,
	}, resp)

	return err
}

func (this *Graph) makeRequest(action []string, params []map[string]interface{}) *GraphMsg {
	var request = &GraphMsg{}
	var graphData = make([]*GraphData, len(action))
	for index := range action {
		var param map[string]interface{} = nil
		if params != nil {
			param = params[index]
		}
		graphData[index] = extraMap[action[index]](this.inst, param)
		graphData[index].Tags = 1
		graphData[index].Time = float32(time.Now().Unix())
		graphData[index].SamplingRate = 1
	}
	request.Data = graphData
	request.AppId = InstagramAppID
	request.Channel = "regular"
	request.Time = float32(time.Now().Unix())
	request.AppVer = InstagramVersion
	request.DeviceId = this.inst.deviceID
	request.FamilyDeviceId = this.inst.familyID
	request.SessionId = this.inst.sessionID
	request.LogType = "client_event"
	request.AppUid = ""
	request.Seq = this.seq
	this.seq++
	return request
}

type GraphResp struct {
	Checksum      string `json:"checksum"`
	Config        string `json:"config"`
	ConfigOwnerId string `json:"config_owner_id"`
	AppData       string `json:"app_data"`
	QplVersion    string `json:"qpl_version"`
}

type GraphMsg struct {
	AppId          string       `json:"app_id"`
	Channel        string       `json:"channel"`
	Time           float32      `json:"time"`
	AppVer         string       `json:"app_ver"`
	DeviceId       string       `json:"device_id"`
	FamilyDeviceId string       `json:"family_device_id"`
	SessionId      string       `json:"session_id"`
	LogType        string       `json:"log_type"`
	AppUid         string       `json:"app_uid"`
	Seq            int          `json:"seq"`
	Data           []*GraphData `json:"data"`
}

type GraphData struct {
	Tags         int         `json:"tags"`
	Module       string      `json:"module"`
	Name         string      `json:"name"`
	Time         float32     `json:"time"`
	SamplingRate int         `json:"sampling_rate"`
	Extra        interface{} `json:"extra"`
}

var extraMap = make(map[string]func(*Instagram, map[string]interface{}) *GraphData)

var beforeSendSMS1 = []string{"next_button_tapped", "reg_field_interacted", "fx_sso_library", "fx_sso_library", "analytics_file_deleted"}
var beforeSendSMS2 = []string{"reg_field_interacted", "next_blocked", "analytics_file_deleted"}
var beforeSendSMS3 = []string{"time_spent_bit_array", "analytics_file_deleted"}
var afterSendSMS = []string{"next_button_tapped", "reg_field_interacted", "fx_sso_library", "fx_sso_library", "ig_emergency_push_did_set_initial_version", "reg_field_interacted", "proceed_with_phone_number", "step_view_loaded"}

func InitGraph() {
	extraMap["step_view_loaded"] = MakeStepViewLoaded
	extraMap["proceed_with_phone_number"] = MakeProceedWithPhoneNumber
	extraMap["ig_emergency_push_did_set_initial_version"] = MakeIgEmergencyPushDidSetInitialVersion
	extraMap["time_spent_bit_array"] = MakeTimeSpentBitArray
	extraMap["next_blocked"] = MakeNextBlocked
	extraMap["analytics_file_deleted"] = MakeAnalyticsFileDeleted
	extraMap["fx_sso_library"] = MakeFxSsoLibrary
	extraMap["reg_field_interacted"] = MakeRegFieldInteracted
	extraMap["next_button_tapped"] = MakeNextButtonTapped
}
