package goinsta

import "makemoney/common"

type Message struct {
	inst *Instagram
}

type RespSendMsg struct {
	BaseApiResp
	Action  string `json:"action"`
	Payload struct {
		ClientContext string `json:"client_context"`
		ItemId        string `json:"item_id"`
		ThreadId      string `json:"thread_id"`
		Timestamp     string `json:"timestamp"`
	} `json:"payload"`
	StatusCode string `json:"status_code"`
}

func (this *Message) SendTextMessage(id string, msg string) error {
	msgID := common.GenString(common.CharSet_123, 19)
	params := map[string]interface{}{
		"recipient_users":  "[[" + id + "]]",
		"action":           "send_item",
		"is_shh_mode":      0,
		"send_attribution": "inbox",
		"client_context":   msgID,
		"text":             msg,
		//"device_id":            this.inst.androidID,
		"mutation_token":       msgID,
		"_uuid":                this.inst.deviceID,
		"nav_chain":            "8Of:self_profile:42,82Y:account_switch_fragment:43,8Of:self_profile:44,TRUNCATEDx10,6Hh:direct_sticker_tab_tray_fragment:71,4tf:direct_thread:72,4tf:direct_thread:73,6Hh:direct_sticker_tab_tray_fragment:74,4tf:direct_thread:75,4tf:direct_thread:76,4tf:direct_thread:77",
		"offline_threading_id": msgID,
	}
	resp := &RespSendMsg{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendText,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Message) SendImgMessage(id string, imageID string) error {
	msgID := common.GenString(common.CharSet_123, 19)
	params := map[string]interface{}{
		"recipient_users":  "[[" + id + "]]",
		"action":           "send_item",
		"is_shh_mode":      0,
		"send_attribution": "inbox",
		"client_context":   msgID,
		//"device_id":               this.inst.androidID,
		"mutation_token":          msgID,
		"_uuid":                   this.inst.deviceID,
		"allow_full_aspect_ratio": true,
		"nav_chain":               "8Of:self_profile:42,82Y:account_switch_fragment:43,8Of:self_profile:44,TRUNCATEDx10,6Hh:direct_sticker_tab_tray_fragment:71,4tf:direct_thread:72,4tf:direct_thread:73,6Hh:direct_sticker_tab_tray_fragment:74,4tf:direct_thread:75,4tf:direct_thread:76,4tf:direct_thread:77",
		"upload_id":               imageID,
		"offline_threading_id":    msgID,
	}
	resp := &RespSendMsg{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendImage,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Message) SendVoiceMessage(id string, videoID string) error {
	msgID := common.GenString(common.CharSet_123, 19)
	params := map[string]interface{}{
		"recipient_users":  "[[" + id + "]]",
		"action":           "send_item",
		"is_shh_mode":      0,
		"send_attribution": "direct_thread",
		"client_context":   msgID,
		"video_result":     "",
		//"device_id":               this.inst.androidID,
		"mutation_token":          msgID,
		"_uuid":                   this.inst.deviceID,
		"allow_full_aspect_ratio": true,
		"nav_chain":               "8Of:self_profile:42,82Y:account_switch_fragment:43,8Of:self_profile:44,TRUNCATEDx10,6Hh:direct_sticker_tab_tray_fragment:71,4tf:direct_thread:72,4tf:direct_thread:73,6Hh:direct_sticker_tab_tray_fragment:74,4tf:direct_thread:75,4tf:direct_thread:76,4tf:direct_thread:77",
		"upload_id":               videoID,
		"offline_threading_id":    msgID,
	}
	resp := &RespSendMsg{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendImage,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)
	return err
}
