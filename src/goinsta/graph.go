package goinsta

import "time"

type Graph struct {
	inst *Instagram
	seq  int
}

func (this *Graph) SendBeforeSendSMS() {
	request1 := this.makeRequest(beforeSendSMS1, nil)
	request2 := this.makeRequest(beforeSendSMS2, nil)
}

func (this *Graph) SendRequest(action []string, params []map[string]interface{}) error {
	request := this.makeRequest(action, params)
	resp := &GraphResp{}
	this.inst.HttpRequestJson(&reqOptions{
		ApiPath:       "",
		IsPost:        true,
		IsApiGraph:    true,
		Signed:        false,
		Query:         nil,
		Body:          nil,
		Header:        nil,
		DisAutoHeader: false,
	}, resp)
}

func (this *Graph) makeRequest(action []string, params []map[string]interface{}) *GraphRequest {
	var request = &GraphRequest{}
	var graphData = make([]*GraphData, len(action))
	for index := range action {
		var param map[string]interface{} = nil
		if params != nil {
			param = params[index]
		}
		graphData[index].Extra = extraMap[action[index]](this.inst, param)
		graphData[index].Tags = 1
		graphData[index].SamplingRate = 1
		graphData[index].Time = float32(time.Now().Unix())
	}
	request.Data = graphData
	request.AppId = AppID
	request.Channel = "regular"
	request.Time = float32(time.Now().Unix())
	request.AppVer = goInstaVersion
	request.DeviceId = this.inst.uuid
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

type GraphRequest struct {
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
