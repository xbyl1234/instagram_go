package goinsta

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"strings"
	"time"
)

type UserOperate struct {
	inst             *Instagram
	cameraEntryPoint int
}

func newUserOperate(inst *Instagram) *UserOperate {
	return &UserOperate{inst: inst}
}

type RespLikeUser struct {
	BaseApiResp
	PreviousFollowing bool       `json:"previous_following,omitempty"`
	FriendshipStatus  Friendship `json:"friendship_status"`
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

type UploadVideoInfo struct {
	UploadMediaInfo
	// "camera"
	// "library"
	from       string
	durationMs float64
}

type Location2 struct {
	Name             string  `json:"name"`
	ExternalId       int64   `json:"external_id"`
	ExternalIdSource string  `json:"external_id_source"`
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	Address          string  `json:"address"`
	MinimumAge       int     `json:"minimum_age"`
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
		"camera_entry_point":              this.cameraEntryPoint,
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

var PostImgNavChain = []string{
	"IGCameraNavigationController:camera_nav:24,IGMediaCaptureViewController:media_capture:25,IGEditorViewController:photo_edit:26,IGShareViewController:media_broadcast_share:27,IGBroadcastShareManager:media_broadcast_share:28",
}

func (this *UserOperate) ConfigurePost(caption string, mediaInfo *UploadMediaInfo, locationReqID string, location Location2) error {
	iso := 30
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
		"device_id":                     this.inst.Device.DeviceID,
		"additional_exif_data":          camera,
		"lens_make":                     "Apple",
		"nav_chain":                     PostImgNavChain[common.GenNumber(0, len(PostImgNavChain))],
		"_uuid":                         this.inst.Device.DeviceID,
		"like_and_view_counts_disabled": false,
		"geotag_enabled":                true,
		"location": map[string]interface{}{
			"facebook_places_id": location.ExternalId,
			"external_source":    location.ExternalIdSource,
			"address":            location.Address,
			"lat":                location.Lat,
			"lng":                location.Lng,
			"follow_status":      0,
			"name":               location.Name,
		},
		"client_timestamp":           time.Now().Unix(),
		"edits":                      map[string]interface{}{},
		"scene_type":                 1,
		"lens_model":                 this.inst.Device.LensModel,
		"iso":                        iso,
		"disable_comments":           false,
		"upload_id":                  mediaInfo.uploadID,
		"caption_list":               []string{caption},
		"source_type":                "library",
		"caption":                    caption,
		"camera_entry_point":         this.cameraEntryPoint,
		"timezone_offset":            this.inst.Device.TimezoneOffset,
		"date_time_original":         timeStr,
		"foursquare_request_id":      locationReqID,
		"waterfall_id":               mediaInfo.waterfall,
		"date_time_digitized":        timeStr,
		"creation_logger_session_id": common.GenString(common.CharSet_16_Num, 32),
		"camera_position":            "back",
		"media_altitude":             common.GenNumber(0, 50),
		"_uid":                       fmt.Sprintf("%d", this.inst.ID),
		"media_longitude":            location.Lng,
		"media_latitude":             location.Lat,
		"allow_multi_configures":     false,
		"container_module":           "photo_edit",
		"video_subtitles_enabled":    true,
		"scene_capture_type":         "standard",
		"software":                   this.inst.Device.SystemVersion,
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

type RespLocation struct {
	BaseApiResp
	Venues    []Location2 `json:"venues"`
	RequestId string      `json:"request_id"`
	RankToken string      `json:"rank_token"`
}

func (this *UserOperate) LocationSearch(latitude float32, longitude float32) (*RespLocation, error) {
	params := map[string]interface{}{
		"latitude":  latitude,
		"longitude": longitude,
		"timestamp": float64(time.Now().Unix()),
		"rankToken": strings.ToUpper(common.GenUUID()),
	}

	resp := &RespLocation{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlLocationSearch,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *UserOperate) SetReelSettings() error {
	params := map[string]interface{}{
		"_uuid":             this.inst.Device.DeviceID,
		"_uid":              fmt.Sprintf("%d", this.inst.ID),
		"reel_auto_archive": 0,
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSetReelSettings,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *UserOperate) CreateReel(title string, mediaID string) error {
	sources := []string{"story_viewer_profile", "self_profile"}
	params := map[string]interface{}{
		"_uuid":       this.inst.Device.DeviceID,
		"_uid":        fmt.Sprintf("%d", this.inst.ID),
		"source":      sources[common.GenNumber(0, len(sources))],
		"creation_id": fmt.Sprintf("%d", time.Now().Unix()),
		"title":       title,
		"media_ids":   mediaID,
		"cover":       "{\"crop_rect\":\"[0,0.21889055472263869,1,0.78110944527736126]\",\"media_id\":\"" + mediaID + "\"}",
	}

	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCreateReel,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *UserOperate) ConfigureToClips(caption string, audioTitle string, mediaInfo *UploadVideoInfo) error {
	var cameraPosition int
	var sourceType string
	if mediaInfo.from == "camera" {
		cameraPosition = 2
		sourceType = "1"
	} else {
		cameraPosition = 3
		sourceType = "0"
	}

	params := map[string]interface{}{
		"remixed_original_sound_params": map[string]interface{}{
			"original_media_id": "",
		},
		"_uuid":             this.inst.Device.DeviceID,
		"_uid":              fmt.Sprintf("%d", this.inst.ID),
		"internal_features": "clips_format,clips_launch",
		"source_type":       sourceType,
		"is_clips_edited":   "0",
		"camera_session_id": common.GenString(common.CharSet_16_Num, 32),
		"additional_audio_info": map[string]interface{}{
			"has_voiceover_attribution": "0",
		},
		"device_id":                   this.inst.Device.DeviceID,
		"client_timestamp":            fmt.Sprintf("%d", time.Now().Unix()),
		"effect_ids":                  []int{},
		"waterfall_id":                mediaInfo.waterfall,
		"caption":                     caption,
		"camera_entry_point":          this.cameraEntryPoint,
		"text_overlay":                []int{},
		"clips_share_preview_to_feed": "1",
		"upload_id":                   mediaInfo.uploadID,
		"sticker_ids":                 []int{},
		"timezone_offset":             this.inst.Device.TimezoneOffset,
		"capture_type":                "clips_v2",
		"clips_audio_metadata": map[string]interface{}{
			"original": map[string]interface{}{
				"volume_level": 1,
			},
			"original_audio_title": audioTitle,
		},
		"clips_segments_metadata": map[string]interface{}{
			"num_segments": 1,
			"clips_segments": []map[string]interface{}{{
				"speed":               100,
				"index":               0,
				"from_draft":          "0",
				"media_type":          "video",
				"original_media_type": "video",
				"source_type":         sourceType,
				"media_folder":        "",
				"audio_type":          "original",
				"face_effect_id":      "",
				"source":              mediaInfo.from,
				"camera_position":     cameraPosition,
				"duration_ms":         mediaInfo.durationMs,
			}},
		},
		"overlay_data": []int{},
	}

	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlConfigureToClips,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *UserOperate) ClipsInfoForCreation() error {

	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlClipsInfoForCreation,
	}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *UserOperate) ClipsAssets(latitude float32, longitude float32) error {
	params := map[string]interface{}{
		"verticalAccuracy":   "10.000000",
		"speed":              "-1.000000",
		"_uuid":              this.inst.Device.DeviceID,
		"timezone_offset":    this.inst.Device.TimezoneOffset,
		"horizontalAccuracy": "65.000000",
		"alt":                "36.008301",
		"_uid":               fmt.Sprintf("%d", this.inst.ID),
		"type":               "static_stickers",
		"lng":                longitude,
		"lat":                latitude,
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlClipsAssets,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

func (this *UserOperate) VerifyOriginalAudioTitle(originalAudioName string) error {
	params := map[string]interface{}{
		"_uuid":               this.inst.Device.DeviceID,
		"original_audio_name": originalAudioName,
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlVerifyOriginalAudioTitle,
		IsPost:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}

type QualityInfo struct {
	OriginalVideoCodec       string  `json:"original_video_codec"`
	EncodedVideoCodec        string  `json:"encoded_video_codec"`
	OriginalColorPrimaries   string  `json:"original_color_primaries"`
	OriginalWidth            int     `json:"original_width"`
	OriginalFrameRate        float64 `json:"original_frame_rate"`
	OriginalTransferFunction string  `json:"original_transfer_function"`
	EncodedHeight            int     `json:"encoded_height"`
	OriginalBitRate          int     `json:"original_bit_rate"`
	EncodedColorPrimaries    string  `json:"encoded_color_primaries"`
	OriginalHeight           int     `json:"original_height"`
	EncodedBitRate           float64 `json:"encoded_bit_rate"`
	EncodedFrameRate         float64 `json:"encoded_frame_rate"`
	EncodedYcbcrMatrix       string  `json:"encoded_ycbcr_matrix"`
	OriginalYcbcrMatrix      string  `json:"original_ycbcr_matrix"`
	EncodedWidth             int     `json:"encoded_width"`
	MeasuredFrames           []struct {
		Ssim      float64 `json:"ssim"`
		Timestamp float64 `json:"timestamp"`
	} `json:"measured_frames"`
	EncodedTransferFunction string `json:"encoded_transfer_function"`
}

func (this *UserOperate) UpdateVideoWithQualityInfo(uploadID string, qualityInfo QualityInfo) error {
	qualityInfoStr, _ := json.Marshal(qualityInfo)
	params := map[string]interface{}{
		"_uuid":        this.inst.Device.DeviceID,
		"_uid":         fmt.Sprintf("%d", this.inst.ID),
		"quality_info": common.B2s(qualityInfoStr),
		"uploadID":     uploadID,
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUpdateVideoWithQualityInfo,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
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
