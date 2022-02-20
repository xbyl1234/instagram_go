package goinsta

import (
	"fmt"
	"makemoney/common"
	"strings"
	"time"
)

type Message struct {
	inst    *Instagram
	ChatMap map[int64]*RespCreateGroup
}

func newMessage(inst *Instagram) *Message {
	return &Message{
		inst:    inst,
		ChatMap: make(map[int64]*RespCreateGroup),
	}
}

type ChatHistory struct {
	Timestamp    string `json:"timestamp"`
	ItemId       string `json:"item_id"`
	ShhSeenState struct {
	} `json:"shh_seen_state"`
	CreatedAt string `json:"created_at"`
}

type InboxCursor struct {
	CursorTimestampSeconds int `json:"cursor_timestamp_seconds"`
	CursorRelevancyScore   int `json:"cursor_relevancy_score"`
	CursorThreadV2Id       int `json:"cursor_thread_v2_id"`
}

type ChatMediaBase struct {
	Id        string `json:"id"`
	MediaType int    `json:"media_type"`
}

type MessageBase struct {
	ItemType               string `json:"item_type"`
	ClientContext          string `json:"client_context"`
	IsSentByViewer         bool   `json:"is_sent_by_viewer"`
	IsShhMode              bool   `json:"is_shh_mode"`
	ItemId                 string `json:"item_id"`
	ShowForwardAttribution bool   `json:"show_forward_attribution"`
	Timestamp              int64  `json:"timestamp"`
	UserId                 int64  `json:"user_id"`
}

type TextMessage struct {
	MessageBase
	Text string `json:"text"`
}

type VoiceMessage struct {
	MessageBase
	VoiceMedia struct {
		Media              ChatMediaBase `json:"media"`
		SeenCount          int           `json:"seen_count"`
		IsShhMode          bool          `json:"is_shh_mode"`
		SeenUserIds        []interface{} `json:"seen_user_ids"`
		ReplayExpiringAtUs interface{}   `json:"replay_expiring_at_us"`
		ViewMode           string        `json:"view_mode"`
	} `json:"voice_media"`
}

type MediaMessage struct {
	MessageBase
	Media ChatMediaBase `json:"media"`
}

type VideoCallMessage struct {
	MessageBase
	VideoCallEvent struct {
		Action                 string        `json:"action"`
		VcId                   int64         `json:"vc_id"`
		EncodedServerDataInfo  string        `json:"encoded_server_data_info"`
		Description            string        `json:"description"`
		TextAttributes         []interface{} `json:"text_attributes"`
		DidJoin                bool          `json:"did_join"`
		ThreadHasAudioOnlyCall bool          `json:"thread_has_audio_only_call"`
		ThreadHasDropIn        bool          `json:"thread_has_drop_in"`
		FeatureSetStr          string        `json:"feature_set_str"`
		CallDuration           int           `json:"call_duration"`
		CallStartTime          int           `json:"call_start_time"`
		CallEndTime            int           `json:"call_end_time"`
	} `json:"video_call_event"`
}

type ChatItems struct {
	VideoCallMessage
	MediaMessage
	VoiceMessage
	TextMessage
}

type RespInbox struct {
	Inbox struct {
		Threads []struct {
			HasOlder                         bool                   `json:"has_older"`
			HasNewer                         bool                   `json:"has_newer"`
			Pending                          bool                   `json:"pending"`
			Items                            []ChatItems            `json:"items"`
			Canonical                        bool                   `json:"canonical"`
			ThreadId                         string                 `json:"thread_id"`
			ThreadV2Id                       string                 `json:"thread_v2_id"`
			Users                            UserSimp               `json:"users"`
			ViewerId                         int64                  `json:"viewer_id"`
			LastActivityAt                   int64                  `json:"last_activity_at"`
			Muted                            bool                   `json:"muted"`
			VcMuted                          bool                   `json:"vc_muted"`
			EncodedServerDataInfo            string                 `json:"encoded_server_data_info"`
			AdminUserIds                     []interface{}          `json:"admin_user_ids"`
			ApprovalRequiredForNewMembers    bool                   `json:"approval_required_for_new_members"`
			Archived                         bool                   `json:"archived"`
			ThreadHasAudioOnlyCall           bool                   `json:"thread_has_audio_only_call"`
			PendingUserIds                   []interface{}          `json:"pending_user_ids"`
			LastSeenAt                       map[string]ChatHistory `json:"last_seen_at"`
			RelevancyScore                   int                    `json:"relevancy_score"`
			RelevancyScoreExpr               int                    `json:"relevancy_score_expr"`
			OldestCursor                     string                 `json:"oldest_cursor"`
			NewestCursor                     string                 `json:"newest_cursor"`
			Named                            bool                   `json:"named"`
			NextCursor                       string                 `json:"next_cursor"`
			PrevCursor                       string                 `json:"prev_cursor"`
			ThreadTitle                      string                 `json:"thread_title"`
			LeftUsers                        []interface{}          `json:"left_users"`
			Spam                             bool                   `json:"spam"`
			BcPartnership                    bool                   `json:"bc_partnership"`
			MentionsMuted                    bool                   `json:"mentions_muted"`
			ThreadType                       string                 `json:"thread_type"`
			ThreadHasDropIn                  bool                   `json:"thread_has_drop_in"`
			VideoCallId                      interface{}            `json:"video_call_id"`
			ShhModeEnabled                   bool                   `json:"shh_mode_enabled"`
			ShhTogglerUserid                 interface{}            `json:"shh_toggler_userid"`
			ShhReplayEnabled                 bool                   `json:"shh_replay_enabled"`
			IsGroup                          bool                   `json:"is_group"`
			InputMode                        int                    `json:"input_mode"`
			ReadState                        int                    `json:"read_state"`
			AssignedAdminId                  int                    `json:"assigned_admin_id"`
			Folder                           int                    `json:"folder"`
			LastNonSenderItemAt              int                    `json:"last_non_sender_item_at"`
			BusinessThreadFolder             int                    `json:"business_thread_folder"`
			ThreadLabel                      int                    `json:"thread_label"`
			MarkedAsUnread                   bool                   `json:"marked_as_unread"`
			IsCloseFriendThread              bool                   `json:"is_close_friend_thread"`
			HasGroupsXacIneligibleUser       bool                   `json:"has_groups_xac_ineligible_user"`
			ThreadImage                      interface{}            `json:"thread_image"`
			IsXacThread                      bool                   `json:"is_xac_thread"`
			IsTranslationEnabled             bool                   `json:"is_translation_enabled"`
			TranslationBannerImpressionCount int                    `json:"translation_banner_impression_count"`
			SystemFolder                     int                    `json:"system_folder"`
			IsFanclubSubscriberThread        bool                   `json:"is_fanclub_subscriber_thread"`
			JoinableGroupLink                string                 `json:"joinable_group_link"`
			GroupLinkJoinableMode            int                    `json:"group_link_joinable_mode"`
			RtcFeatureSetStr                 string                 `json:"rtc_feature_set_str"`
		} `json:"threads"`
		HasOlder            bool        `json:"has_older"`
		UnseenCount         int         `json:"unseen_count"`
		UnseenCountTs       int64       `json:"unseen_count_ts"`
		PrevCursor          InboxCursor `json:"prev_cursor"`
		NextCursor          InboxCursor `json:"next_cursor"`
		BlendedInboxEnabled bool        `json:"blended_inbox_enabled"`
	} `json:"inbox"`
	SeqId                 int    `json:"seq_id"`
	SnapshotAtMs          int64  `json:"snapshot_at_ms"`
	PendingRequestsTotal  int    `json:"pending_requests_total"`
	HasPendingTopRequests bool   `json:"has_pending_top_requests"`
	Status                string `json:"status"`
}

func (this *Message) FetchInbox() {

}

type RespCreateGroup struct {
	BaseApiResp
	ThreadId       string `json:"thread_id"`
	ThreadV2Id     string `json:"thread_v2_id"`
	LastActivityAt int    `json:"last_activity_at"`
	BcPartnership  bool   `json:"bc_partnership"`
	ThreadType     string `json:"thread_type"`
	ViewerId       int64  `json:"viewer_id"`
	IsGroup        bool   `json:"is_group"`
}

var navChain = []string{
	"IGProfileViewController:self_profile:2,IGFollowListTabPageViewController:self_unified_follow_lists:12,IGProfileViewController:profile:13",
	"IGMainFeedViewController:feed_timeline:1,IGDirectInboxNavigationController:direct_inbox:3,IGDirectInboxViewController:direct_inbox:4",
	"IGExploreViewController:explore_popular:19,IGProfileViewController:profile:21",
	"IGMainFeedViewController:feed_timeline:1,IGDirectInboxNavigationController:direct_inbox:3,IGDirectInboxViewController:direct_inbox:4",
	"IGProfileViewController:self_profile:2,IGFollowListTabPageViewController:self_unified_follow_lists:12,IGProfileViewController:profile:13",
}

var sendAttribution = []string{
	"message_button",
	"inbox",
	"direct_inbox",
	"thread_view",
}

var waveform = [][]float32{
	{0.2969, 0.3027, 0.3111, 0.4661, 0.4467, 0.4105, 0.3816, 0.4424, 0.4384, 0.3946, 0.3677, 0.4171, 0.4220, 0.3782, 0.3150, 0.2716, 0.4209, 0.3486, 0.3082, 0.3758, 0.3453, 0.3208, 0.3052, 0.2639, 0.2182, 0.3558, 0.2750, 0.3337, 0.3347, 0.2555, 0.2774, 0.3700, 0.3604, 0.3395, 0.3033, 0.3551, 0.2919, 0.2630, 0.2123},
}

func (this *Message) CreateGroup(id int64) (*RespCreateGroup, error) {
	params := map[string]interface{}{
		"client_context":  common.GenUUID(),
		"_uuid":           this.inst.Device.DeviceID,
		"recipient_users": "[[" + fmt.Sprintf("%d", id) + "]]",
		"_uid":            this.inst.ID,
	}
	resp := &RespCreateGroup{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		Signed:  true,
		ApiPath: urlCreateGroupThread,
		Header: map[string]string{
			"Source": "thread_presenter_message_button",
		},
		Query: params,
	}, resp)
	err = resp.CheckError(err)
	return resp, err
}

func GenChatID() string {
	return fmt.Sprintf("%d", (time.Now().UnixMilli()<<22)|(int64(common.GenNumber(0, 9999999))&4194303)&0x7fffffffffffffff)
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
}

//58
func (this *Message) GetThreadId(id int64) (string, error) {
	var err error
	var chatInfo = this.ChatMap[id]
	if chatInfo == nil {
		chatInfo, err = this.CreateGroup(id)
		if err != nil {
			return "", err
		}
		this.ChatMap[id] = chatInfo
	}
	return chatInfo.ThreadId, err
}

func (this *Message) SendTextMessage(id int64, msg string) error {
	msgID := GenChatID()
	threadID, err := this.GetThreadId(id)
	if err != nil {
		return err
	}

	this.inst.Increase(OperNameSendMsg)
	params := map[string]interface{}{
		"thread_id":            threadID,
		"action":               "send_item",
		"is_shh_mode":          0,
		"send_attribution":     sendAttribution[common.GenNumber(0, len(sendAttribution))],
		"client_context":       msgID,
		"text":                 msg,
		"device_id":            this.inst.Device.DeviceID,
		"mutation_token":       msgID,
		"_uuid":                this.inst.Device.DeviceID,
		"nav_chain":            navChain[common.GenNumber(0, len(navChain))],
		"offline_threading_id": msgID,
	}
	resp := &RespSendMsg{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendText,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Message) SendImgMessage(id int64, imageID string) error {
	msgID := GenChatID()
	threadID, err := this.GetThreadId(id)
	if err != nil {
		return err
	}
	this.inst.Increase(OperNameSendMsg)
	params := map[string]interface{}{
		"thread_id":        threadID,
		"client_timestamp": time.Now().Unix(),
		"timezone_offset":  this.inst.Device.TimezoneOffset,
		"mutation_token":   msgID,
		//"nav_chain":               navChain[common.GenNumber(0, len(navChain))],
		"content_type":            "photo",
		"_uuid":                   this.inst.Device.DeviceID,
		"action":                  "send_item",
		"allow_full_aspect_ratio": 1,
		"waterfall_id":            "",
		"offline_threading_id":    msgID,
		"upload_id":               imageID,
		"device_id":               this.inst.Device.DeviceID,
		"send_attribution":        sendAttribution[common.GenNumber(0, len(sendAttribution))],
		"client_context":          msgID,
		"is_shh_mode":             0,
	}
	resp := &RespSendMsg{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendImage,
		Header: map[string]string{
			"X-Ig-Connection-Speed":      "429kbps",
			"X-Ig-Eu-Configure-Disabled": "true",
		},
		Query: params,
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Message) SendVoiceMessage(id int64, voiceID string) error {
	msgID := GenChatID()
	threadID, err := this.GetThreadId(id)
	if err != nil {
		return err
	}
	this.inst.Increase(OperNameSendMsg)
	params := map[string]interface{}{
		"is_shh_mode":                    0,
		"thread_id":                      threadID,
		"client_timestamp":               time.Now().Unix(),
		"timezone_offset":                this.inst.Device.TimezoneOffset,
		"mutation_token":                 msgID,
		"nav_chain":                      navChain[common.GenNumber(0, len(navChain))],
		"content_type":                   "audio",
		"_uuid":                          this.inst.Device.DeviceID,
		"action":                         "send_item",
		"allow_full_aspect_ratio":        1,
		"waterfall_id":                   "",
		"offline_threading_id":           msgID,
		"upload_id":                      voiceID,
		"device_id":                      this.inst.Device.DeviceID,
		"waveform_sampling_frequency_hz": 10,
		"client_context":                 msgID,
		"waveform":                       waveform[common.GenNumber(0, len(waveform))],
		"send_attribution":               sendAttribution[common.GenNumber(0, len(sendAttribution))],
	}
	resp := &RespSendMsg{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlShareVoice,
		Header: map[string]string{
			"Priority": "u=2, i",
			//"X-Ig-Salt-Ids: 42139649",
		},
		Query: params,
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Message) SendLinkMessage(id int64, link string) error {
	msgID := GenChatID()
	threadID, err := this.GetThreadId(id)
	if err != nil {
		return err
	}

	this.inst.Increase(OperNameSendMsg)
	params := map[string]interface{}{
		"thread_id":            threadID,
		"mutation_token":       msgID,
		"nav_chain":            navChain[common.GenNumber(0, len(navChain))],
		"link_urls":            "[\"" + strings.ReplaceAll(link, "/", "\\/") + "\"]",
		"_uuid":                this.inst.Device.DeviceID,
		"link_text":            link,
		"action":               "send_item",
		"offline_threading_id": msgID,
		"text":                 link,
		"is_shh_mode":          0,
		"client_context":       msgID,
		"device_id":            this.inst.Device.DeviceID,
		"send_attribution":     sendAttribution[common.GenNumber(0, len(sendAttribution))],
	}
	resp := &RespSendMsg{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:  true,
		ApiPath: urlSendLink,
		Query:   params,
	}, resp)
	err = resp.CheckError(err)
	return err
}
