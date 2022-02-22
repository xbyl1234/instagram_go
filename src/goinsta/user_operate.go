package goinsta

import (
	"fmt"
	"makemoney/common"
	"strings"
	"time"
)

type UserOperate struct {
	inst *Instagram
}

func newUserOperate(inst *Instagram) *UserOperate {
	return &UserOperate{inst: inst}
}

type RespLikeUser struct {
	BaseApiResp
	PreviousFollowing bool       `json:"previous_following,omitempty"`
	FriendshipStatus  Friendship `json:"friendship_status"`
}

func (this *UserOperate) LikeUser(userID int64) error {
	this.inst.Increase(OperNameLikeUser)
	params := map[string]interface{}{
		"_uuid":            this.inst.Device.DeviceID,
		"_uid":             this.inst.ID,
		"user_id":          userID,
		"device_id":        this.inst.Device.DeviceID,
		"container_module": "profile",
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        fmt.Sprintf(urlUserFollow, userID),
		HeaderSequence: LoginHeaderMap[urlUserFollow],
		IsPost:         true,
		Signed:         true,
		Query:          params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

type CameraSettings struct {
	FocalLength  float64 `bson:"focal_length"`
	Aperture     float64 `bson:"aperture"`
	Iso          []int   `bson:"iso"`
	ShutterSpeed float64 `bson:"shutter_speed"`
	MeteringMode int     `bson:"metering_mode"`
	ExposureTime float64 `bson:"exposure_time"`
	Software     string  `bson:"software"`
	LensModel    string  `bson:"lens_model"`
	FlashStatus  int     `bson:"flash_status"`
}

type UploadMediaInfo struct {
	uploadID  string
	waterfall string
	high      int
	width     int
}

func (this *UserOperate) ConfigureToStory(mediaInfo *UploadMediaInfo) error {
	iso := 30
	sourceType := "camera"
	//sourceType:="library"
	mediaType := "photo"
	camera := map[string]interface{}{
		"camera_settings": &CameraSettings{
			FocalLength:  this.inst.Device.FocalLength,
			Aperture:     this.inst.Device.Aperture,
			Iso:          []int{iso},
			ShutterSpeed: float64(common.GenNumber(1, 10)),
			MeteringMode: 3,
			ExposureTime: float64(common.GenNumber(1, 10)) / 100.0,
			Software:     this.inst.Device.SystemVersion,
			LensModel:    this.inst.Device.LensModel,
			FlashStatus:  0,
		},
	}

	timeStr := common.GetNewYorkTimeString()
	params := map[string]interface{}{
		"device_id":                       this.inst.Device.DeviceID,
		"private_mention_sharing_enabled": false,
		"additional_exif_data":            camera,
		"lens_make":                       "Apple",
		"_uuid":                           this.inst.Device.DeviceID,
		"like_and_view_counts_disabled":   false,
		"capture_type":                    "normal",
		"geotag_enabled":                  false,
		"archived_media_id":               "",
		"client_timestamp":                fmt.Sprintf("%d", time.Now().Unix()),
		"edits":                           map[string]interface{}{},
		"original_media_size":             fmt.Sprintf("{%d, %d}", mediaInfo.width, mediaInfo.high),
		"scene_type":                      1,
		"lens_model":                      this.inst.Device.LensModel,
		"camera_session_id":               common.GenString(common.CharSet_16_Num, 32),
		"iso":                             iso,
		"has_animated_sticker":            false,
		"upload_id":                       mediaInfo.uploadID,
		"camera_entry_point":              13,
		"source_type":                     sourceType,
		"configure_mode":                  1,
		"disable_comments":                false,
		"timezone_offset":                 this.inst.Device.TimezoneOffset,
		"date_time_original":              timeStr,
		"waterfall_id":                    mediaInfo.waterfall,
		"composition_id":                  strings.ToUpper(common.GenUUID()),
		"date_time_digitized":             timeStr,
		"camera_position":                 "back",
		"_uid":                            fmt.Sprintf("%d", this.inst.ID),
		"client_context":                  mediaInfo.uploadID,
		"original_media_type":             mediaType,
		//"client_shared_at":                fmt.Sprintf("%d", time.Now().Unix()),
		"client_shared_at":        mediaInfo.uploadID[:10],
		"allow_multi_configures":  true,
		"container_module":        "direct_story_audience_picker",
		"creation_surface":        "camera",
		"video_subtitles_enabled": true,
		"from_drafts":             false,
		"software":                this.inst.Device.SystemVersion,
		"media_gesture":           0,
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlConfigureToStory,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

//
//func (this *UserOperate) CreateReel() {
//	params := map[string]interface{}{
//		"_uuid":            this.inst.Device.DeviceID,
//		"_uid":             this.inst.ID,
//		"user_id":          userID,
//		"device_id":        this.inst.Device.DeviceID,
//		"container_module": "profile",
//	}
//	resp := &RespLikeUser{}
//	err := this.inst.HttpRequestJson(&reqOptions{
//		ApiPath:        urlCreateReel
//		HeaderSequence: LoginHeaderMap[urlUserFollow],
//		IsPost:         true,
//		Signed:         true,
//		Query:          params,
//	}, resp)
//
//	err = resp.CheckError(err)
//	return err
//}
