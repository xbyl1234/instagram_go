package goinsta

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"strings"
	"time"
)

type UserOperate struct {
	inst *Instagram
	//for post
	cameraEntryPoint int
	//for comment
	commentSessionId string
}

func newUserOperate(inst *Instagram) *UserOperate {
	return &UserOperate{
		inst:             inst,
		cameraEntryPoint: 0,
		commentSessionId: common.GenString(common.CharSet_16_Num, 32),
	}
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
	UploadID  string
	Waterfall string
	High      int
	Width     int
}

type UploadVideoInfo struct {
	UploadMediaInfo

	from       string
	durationMs float64
}

type LocationSearch struct {
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
			FocalLength:  this.inst.AccountInfo.Device.FocalLength,
			Aperture:     this.inst.AccountInfo.Device.Aperture,
			Iso:          []int{iso},
			ShutterSpeed: float64(common.GenNumber(1, 10)),
			MeteringMode: 3,
			ExposureTime: float64(common.GenNumber(1, 10)) / 100.0,
			Software:     this.inst.AccountInfo.Device.SystemVersion,
			LensModel:    this.inst.AccountInfo.Device.LensModel,
			FlashStatus:  0,
		},
	}

	timeStr := common.GetNewYorkTimeString()
	params := map[string]interface{}{
		"device_id":                       this.inst.AccountInfo.Device.DeviceID,
		"private_mention_sharing_enabled": false,
		"additional_exif_data":            camera,
		"lens_make":                       "Apple",
		"_uuid":                           this.inst.AccountInfo.Device.DeviceID,
		"like_and_view_counts_disabled":   false,
		"capture_type":                    "normal",
		"geotag_enabled":                  false,
		"archived_media_id":               "",
		"client_timestamp":                fmt.Sprintf("%d", time.Now().Unix()),
		"edits":                           map[string]interface{}{},
		"original_media_size":             fmt.Sprintf("{%d, %d}", mediaInfo.Width, mediaInfo.High),
		"scene_type":                      1,
		"lens_model":                      this.inst.AccountInfo.Device.LensModel,
		"camera_session_id":               common.GenString(common.CharSet_16_Num, 32),
		"iso":                             iso,
		"has_animated_sticker":            false,
		"upload_id":                       mediaInfo.UploadID,
		"camera_entry_point":              this.cameraEntryPoint,
		"source_type":                     sourceType,
		"configure_mode":                  1,
		"disable_comments":                false,
		"timezone_offset":                 this.inst.AccountInfo.Location.Timezone,
		"date_time_original":              timeStr,
		"waterfall_id":                    mediaInfo.Waterfall,
		"composition_id":                  strings.ToUpper(common.GenUUID()),
		"date_time_digitized":             timeStr,
		"camera_position":                 "back",
		"_uid":                            fmt.Sprintf("%d", this.inst.ID),
		"client_context":                  mediaInfo.UploadID,
		"original_media_type":             mediaType,
		//"client_shared_at":                fmt.Sprintf("%d", time.Now().Unix()),
		"client_shared_at":        mediaInfo.UploadID[:10],
		"allow_multi_configures":  true,
		"container_module":        "direct_story_audience_picker",
		"creation_surface":        "camera",
		"video_subtitles_enabled": true,
		"from_drafts":             false,
		"software":                this.inst.AccountInfo.Device.SystemVersion,
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

type RespConfigure struct {
	BaseApiResp
	Media struct {
		TakenAt         int    `json:"taken_at"`
		Pk              int64  `json:"pk"`
		Id              string `json:"id"`
		DeviceTimestamp int64  `json:"device_timestamp"`
		MediaType       int    `json:"media_type"`
		Code            string `json:"code"`
		//ClientCacheKey  string `json:"client_cache_key"`
		//FilterType      int    `json:"filter_type"`
		//IsUnifiedVideo  bool   `json:"is_unified_video"`
		//Location        struct {
		//	Pk                  int64   `json:"pk"`
		//	ShortName           string  `json:"short_name"`
		//	FacebookPlacesId    int64   `json:"facebook_places_id"`
		//	ExternalSource      string  `json:"external_source"`
		//	Name                string  `json:"name"`
		//	Address             string  `json:"address"`
		//	City                string  `json:"city"`
		//	HasViewerSaved      bool    `json:"has_viewer_saved"`
		//	Lng                 float64 `json:"lng"`
		//	Lat                 float64 `json:"lat"`
		//	IsEligibleForGuides bool    `json:"is_eligible_for_guides"`
		//} `json:"location"`
		//Lat float64 `json:"lat"`
		//Lng float64 `json:"lng"`
		//User struct {
		//	Pk                         int64  `json:"pk"`
		//	Username                   string `json:"username"`
		//	FullName                   string `json:"full_name"`
		//	IsPrivate                  bool   `json:"is_private"`
		//	ProfilePicUrl              string `json:"profile_pic_url"`
		//	ProfilePicId               string `json:"profile_pic_id"`
		//	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		//	CanBoostPost               bool   `json:"can_boost_post"`
		//	CanSeeOrganicInsights      bool   `json:"can_see_organic_insights"`
		//	ShowInsightsTerms          bool   `json:"show_insights_terms"`
		//	ReelAutoArchive            string `json:"reel_auto_archive"`
		//	IsUnpublished              bool   `json:"is_unpublished"`
		//	AllowedCommenterType       string `json:"allowed_commenter_type"`
		//	HasHighlightReels          bool   `json:"has_highlight_reels"`
		//	InteropMessagingUserFbid   int64  `json:"interop_messaging_user_fbid"`
		//	FbidV2                     int64  `json:"fbid_v2"`
		//} `json:"user"`
		//CanViewerReshare                    bool          `json:"can_viewer_reshare"`
		//CaptionIsEdited                     bool          `json:"caption_is_edited"`
		//LikeAndViewCountsDisabled           bool          `json:"like_and_view_counts_disabled"`
		//CommercialityStatus                 string        `json:"commerciality_status"`
		//IsPaidPartnership                   bool          `json:"is_paid_partnership"`
		//IsVisualReplyCommenterNoticeEnabled bool          `json:"is_visual_reply_commenter_notice_enabled"`
		//CommentLikesEnabled                 bool          `json:"comment_likes_enabled"`
		//CommentThreadingEnabled             bool          `json:"comment_threading_enabled"`
		//HasMoreComments                     bool          `json:"has_more_comments"`
		//MaxNumVisiblePreviewComments        int           `json:"max_num_visible_preview_comments"`
		//PreviewComments                     []interface{} `json:"preview_comments"`
		//CanViewMorePreviewComments          bool          `json:"can_view_more_preview_comments"`
		//CommentCount                        int           `json:"comment_count"`
		//HideViewAllCommentEntrypoint        bool          `json:"hide_view_all_comment_entrypoint"`
		//ImageVersions2                      struct {
		//	Candidates []struct {
		//		Width        int    `json:"width"`
		//		Height       int    `json:"height"`
		//		Url          string `json:"url"`
		//		ScansProfile string `json:"scans_profile"`
		//	} `json:"candidates"`
		//} `json:"image_versions2"`
		//OriginalWidth         int  `json:"original_width"`
		//OriginalHeight        int  `json:"original_height"`
		//PhotoOfYou            bool `json:"photo_of_you"`
		//CanSeeInsightsAsBrand bool `json:"can_see_insights_as_brand"`
		//Caption               struct {
		//	Pk              int64  `json:"pk"`
		//	UserId          int64  `json:"user_id"`
		//	Text            string `json:"text"`
		//	Type            int    `json:"type"`
		//	CreatedAt       int    `json:"created_at"`
		//	CreatedAtUtc    int    `json:"created_at_utc"`
		//	ContentType     string `json:"content_type"`
		//	Status          string `json:"status"`
		//	BitFlags        int    `json:"bit_flags"`
		//	DidReportAsSpam bool   `json:"did_report_as_spam"`
		//	ShareEnabled    bool   `json:"share_enabled"`
		//	User            struct {
		//		Pk                         int64  `json:"pk"`
		//		Username                   string `json:"username"`
		//		FullName                   string `json:"full_name"`
		//		IsPrivate                  bool   `json:"is_private"`
		//		ProfilePicUrl              string `json:"profile_pic_url"`
		//		ProfilePicId               string `json:"profile_pic_id"`
		//		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		//		CanBoostPost               bool   `json:"can_boost_post"`
		//		CanSeeOrganicInsights      bool   `json:"can_see_organic_insights"`
		//		ShowInsightsTerms          bool   `json:"show_insights_terms"`
		//		ReelAutoArchive            string `json:"reel_auto_archive"`
		//		IsUnpublished              bool   `json:"is_unpublished"`
		//		AllowedCommenterType       string `json:"allowed_commenter_type"`
		//		HasHighlightReels          bool   `json:"has_highlight_reels"`
		//		InteropMessagingUserFbid   int64  `json:"interop_messaging_user_fbid"`
		//		FbidV2                     int64  `json:"fbid_v2"`
		//	} `json:"user"`
		//	IsCovered          bool  `json:"is_covered"`
		//	MediaId            int64 `json:"media_id"`
		//	PrivateReplyStatus int   `json:"private_reply_status"`
		//} `json:"caption"`
		//FbUserTags struct {
		//	In []interface{} `json:"in"`
		//} `json:"fb_user_tags"`
		//CanViewerSave        bool   `json:"can_viewer_save"`
		//OrganicTrackingToken string `json:"organic_tracking_token"`
		//SharingFrictionInfo  struct {
		//	ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
		//	BloksAppUrl               interface{} `json:"bloks_app_url"`
		//} `json:"sharing_friction_info"`
		//CommentInformTreatment struct {
		//	ShouldHaveInformTreatment bool   `json:"should_have_inform_treatment"`
		//	Text                      string `json:"text"`
		//} `json:"comment_inform_treatment"`
		//ProductType               string `json:"product_type"`
		//IsInProfileGrid           bool   `json:"is_in_profile_grid"`
		//ProfileGridControlEnabled bool   `json:"profile_grid_control_enabled"`
		//DeletedReason             int    `json:"deleted_reason"`
		//IntegrityReviewDecision   string `json:"integrity_review_decision"`
		//MusicMetadata             struct {
		//	MusicCanonicalId  string      `json:"music_canonical_id"`
		//	AudioType         interface{} `json:"audio_type"`
		//	MusicInfo         interface{} `json:"music_info"`
		//	OriginalSoundInfo interface{} `json:"original_sound_info"`
		//} `json:"music_metadata"`
		//PostFrictionInfo struct {
		//	ShouldShowPostFriction bool   `json:"should_show_post_friction"`
		//	PostFrictionReason     string `json:"post_friction_reason"`
		//} `json:"post_friction_info"`
	} `json:"media"`
	UploadId string `json:"upload_id"`
}

func (this *UserOperate) ConfigurePost(caption string, mediaInfo *UploadMediaInfo, locationReqID string, location *LocationSearch) (*RespConfigure, error) {
	iso := 30
	camera := map[string]interface{}{
		"camera_settings": &CameraSettings{
			FocalLength:  this.inst.AccountInfo.Device.FocalLength,
			Aperture:     this.inst.AccountInfo.Device.Aperture,
			Iso:          []int{iso},
			ShutterSpeed: float64(common.GenNumber(1, 10)),
			MeteringMode: 3,
			ExposureTime: float64(common.GenNumber(1, 10)) / 100.0,
			Software:     this.inst.AccountInfo.Device.SystemVersion,
			LensModel:    this.inst.AccountInfo.Device.LensModel,
			FlashStatus:  0,
		},
	}

	timeStr := common.GetNewYorkTimeString()
	params := map[string]interface{}{
		"device_id":                     this.inst.AccountInfo.Device.DeviceID,
		"additional_exif_data":          camera,
		"lens_make":                     "Apple",
		"nav_chain":                     PostImgNavChain[common.GenNumber(0, len(PostImgNavChain))],
		"_uuid":                         this.inst.AccountInfo.Device.DeviceID,
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
		"lens_model":                 this.inst.AccountInfo.Device.LensModel,
		"iso":                        iso,
		"disable_comments":           false,
		"upload_id":                  mediaInfo.UploadID,
		"caption_list":               []string{caption},
		"source_type":                "library",
		"caption":                    caption,
		"camera_entry_point":         this.cameraEntryPoint,
		"timezone_offset":            this.inst.AccountInfo.Location.Timezone,
		"date_time_original":         timeStr,
		"foursquare_request_id":      locationReqID,
		"waterfall_id":               mediaInfo.Waterfall,
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
		"software":                   this.inst.AccountInfo.Device.SystemVersion,
	}
	resp := &RespConfigure{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlConfigure,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespLocation struct {
	BaseApiResp
	Venues    []LocationSearch `json:"venues"`
	RequestId string           `json:"request_id"`
	RankToken string           `json:"rank_token"`
}

func (this *UserOperate) LocationSearch(longitude float32, latitude float32) (*RespLocation, error) {
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

type RespSetReel struct {
	BaseApiResp
	ReelAutoArchive string      `json:"reel_auto_archive"`
	MessagePrefs    interface{} `json:"message_prefs"`
}

func (this *UserOperate) SetReelSettings() (*RespSetReel, error) {
	params := map[string]interface{}{
		"_uuid":             this.inst.AccountInfo.Device.DeviceID,
		"_uid":              fmt.Sprintf("%d", this.inst.ID),
		"reel_auto_archive": 0,
	}
	resp := &RespSetReel{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSetReelSettings,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespCreateReel struct {
	BaseApiResp
	Reel struct {
		Id              string      `json:"id"`
		LatestReelMedia int         `json:"latest_reel_media"`
		Seen            interface{} `json:"seen"`
		//CanReply                  bool        `json:"can_reply"`
		//CanGifQuickReply          bool        `json:"can_gif_quick_reply"`
		//CanReshare                bool        `json:"can_reshare"`
		ReelType string `json:"reel_type"`
		//AdExpiryTimestampInMillis interface{} `json:"ad_expiry_timestamp_in_millis"`
		//IsCtaStickerAvailable     interface{} `json:"is_cta_sticker_available"`
		//CoverMedia                struct {
		//	CroppedImageVersion struct {
		//		Width  int    `json:"width"`
		//		Height int    `json:"height"`
		//		Url    string `json:"url"`
		//	} `json:"cropped_image_version"`
		//	CropRect         interface{} `json:"crop_rect"`
		//	FullImageVersion struct {
		//		Width  int    `json:"width"`
		//		Height int    `json:"height"`
		//		Url    string `json:"url"`
		//	} `json:"full_image_version"`
		//} `json:"cover_media"`
		//User struct {
		//	Pk                 int64  `json:"pk"`
		//	Username           string `json:"username"`
		//	FullName           string `json:"full_name"`
		//	IsPrivate          bool   `json:"is_private"`
		//	ProfilePicUrl      string `json:"profile_pic_url"`
		//	ProfilePicId       string `json:"profile_pic_id"`
		//	IsVerified         bool   `json:"is_verified"`
		//	FollowFrictionType int    `json:"follow_friction_type"`
		//} `json:"user"`
		//Items []struct {
		//TakenAt         int    `json:"taken_at"`
		//Pk              int64  `json:"pk"`
		//Id              string `json:"id"`
		//DeviceTimestamp int64  `json:"device_timestamp"`
		//MediaType       int    `json:"media_type"`
		//Code            string `json:"code"`
		//ClientCacheKey  string `json:"client_cache_key"`
		//FilterType      int    `json:"filter_type"`
		//IsUnifiedVideo  bool   `json:"is_unified_video"`
		//User            struct {
		//	Pk                         int64  `json:"pk"`
		//	Username                   string `json:"username"`
		//	FullName                   string `json:"full_name"`
		//	IsPrivate                  bool   `json:"is_private"`
		//	ProfilePicUrl              string `json:"profile_pic_url"`
		//	ProfilePicId               string `json:"profile_pic_id"`
		//	IsVerified                 bool   `json:"is_verified"`
		//	FollowFrictionType         int    `json:"follow_friction_type"`
		//	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		//	CanBoostPost               bool   `json:"can_boost_post"`
		//	CanSeeOrganicInsights      bool   `json:"can_see_organic_insights"`
		//	ShowInsightsTerms          bool   `json:"show_insights_terms"`
		//	ReelAutoArchive            string `json:"reel_auto_archive"`
		//	IsUnpublished              bool   `json:"is_unpublished"`
		//	AllowedCommenterType       string `json:"allowed_commenter_type"`
		//	HasHighlightReels          bool   `json:"has_highlight_reels"`
		//	InteropMessagingUserFbid   int64  `json:"interop_messaging_user_fbid"`
		//	FbidV2                     int64  `json:"fbid_v2"`
		//} `json:"user"`
		//CaptionIsEdited                     bool   `json:"caption_is_edited"`
		//LikeAndViewCountsDisabled           bool   `json:"like_and_view_counts_disabled"`
		//CommercialityStatus                 string `json:"commerciality_status"`
		//IsPaidPartnership                   bool   `json:"is_paid_partnership"`
		//IsVisualReplyCommenterNoticeEnabled bool   `json:"is_visual_reply_commenter_notice_enabled"`
		//ImageVersions2                      struct {
		//	Candidates []struct {
		//		Width        int    `json:"width"`
		//		Height       int    `json:"height"`
		//		Url          string `json:"url"`
		//		ScansProfile string `json:"scans_profile"`
		//	} `json:"candidates"`
		//} `json:"image_versions2"`
		//OriginalWidth         int         `json:"original_width"`
		//OriginalHeight        int         `json:"original_height"`
		//CaptionPosition       float64     `json:"caption_position"`
		//IsReelMedia           bool        `json:"is_reel_media"`
		//TimezoneOffset        int         `json:"timezone_offset"`
		//PhotoOfYou            bool        `json:"photo_of_you"`
		//CanSeeInsightsAsBrand bool        `json:"can_see_insights_as_brand"`
		//Caption               interface{} `json:"caption"`
		//FbUserTags            struct {
		//	In []interface{} `json:"in"`
		//} `json:"fb_user_tags"`
		//CanViewerSave        bool   `json:"can_viewer_save"`
		//OrganicTrackingToken string `json:"organic_tracking_token"`
		//SharingFrictionInfo  struct {
		//	ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
		//	BloksAppUrl               interface{} `json:"bloks_app_url"`
		//} `json:"sharing_friction_info"`
		//CommentInformTreatment struct {
		//	ShouldHaveInformTreatment bool   `json:"should_have_inform_treatment"`
		//	Text                      string `json:"text"`
		//} `json:"comment_inform_treatment"`
		//ProductType               string        `json:"product_type"`
		//IsInProfileGrid           bool          `json:"is_in_profile_grid"`
		//ProfileGridControlEnabled bool          `json:"profile_grid_control_enabled"`
		//DeletedReason             int           `json:"deleted_reason"`
		//IntegrityReviewDecision   string        `json:"integrity_review_decision"`
		//MusicMetadata             interface{}   `json:"music_metadata"`
		//CanReshare                bool          `json:"can_reshare"`
		//CanReply                  bool          `json:"can_reply"`
		//StoryIsSavedToArchive     bool          `json:"story_is_saved_to_archive"`
		//StoryStaticModels         []interface{} `json:"story_static_models"`
		//HighlightReelIds          []string      `json:"highlight_reel_ids"`
		//Viewers                   []struct {
		//	Pk                 int64  `json:"pk"`
		//	Username           string `json:"username"`
		//	FullName           string `json:"full_name"`
		//	IsPrivate          bool   `json:"is_private"`
		//	ProfilePicUrl      string `json:"profile_pic_url"`
		//	ProfilePicId       string `json:"profile_pic_id"`
		//	IsVerified         bool   `json:"is_verified"`
		//	FollowFrictionType int    `json:"follow_friction_type"`
		//} `json:"viewers"`
		//ViewerCount              int           `json:"viewer_count"`
		//FbViewerCount            interface{}   `json:"fb_viewer_count"`
		//ViewerCursor             interface{}   `json:"viewer_cursor"`
		//TotalViewerCount         int           `json:"total_viewer_count"`
		//MultiAuthorReelNames     []interface{} `json:"multi_author_reel_names"`
		//SupportsReelReactions    bool          `json:"supports_reel_reactions"`
		//CanSendCustomEmojis      bool          `json:"can_send_custom_emojis"`
		//ShowOneTapFbShareTooltip bool          `json:"show_one_tap_fb_share_tooltip"`
		//HasSharedToFb            int           `json:"has_shared_to_fb"`
		//HasSharedToFbDating      int           `json:"has_shared_to_fb_dating"`
		//SourceType               int           `json:"source_type"`
		//} `json:"items"`
		//RankedPosition                   int     `json:"ranked_position"`
		//Title                            string  `json:"title"`
		CreatedAt         int  `json:"created_at"`
		IsPinnedHighlight bool `json:"is_pinned_highlight"`
		//SeenRankedPosition               int     `json:"seen_ranked_position"`
		//PrefetchCount                    int     `json:"prefetch_count"`
		//MediaCount                       int     `json:"media_count"`
		MediaIds []int64 `json:"media_ids"`
		//ContainsStitchedMediaBlockedByRm bool    `json:"contains_stitched_media_blocked_by_rm"`
		//IsConvertedToClips               bool    `json:"is_converted_to_clips"`
	} `json:"reel"`
}

func (this *UserOperate) CreateReel(title string, mediaID string) (*RespCreateReel, error) {
	sources := []string{"story_viewer_profile", "self_profile"}
	params := map[string]interface{}{
		"_uuid":       this.inst.AccountInfo.Device.DeviceID,
		"_uid":        fmt.Sprintf("%d", this.inst.ID),
		"source":      sources[common.GenNumber(0, len(sources))],
		"creation_id": fmt.Sprintf("%d", time.Now().Unix()),
		"title":       title,
		"media_ids":   mediaID,
		"cover":       "{\"crop_rect\":\"[0,0.21889055472263869,1,0.78110944527736126]\",\"media_id\":\"" + mediaID + "\"}",
	}

	resp := &RespCreateReel{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCreateReel,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespConfigureClips struct {
	BaseApiResp
	Media struct {
		TakenAt         int    `json:"taken_at"`
		Pk              int64  `json:"pk"`
		Id              string `json:"id"`
		DeviceTimestamp int64  `json:"device_timestamp"`
		MediaType       int    `json:"media_type"`
		Code            string `json:"code"`
		ClientCacheKey  string `json:"client_cache_key"`
		FilterType      int    `json:"filter_type"`
		//IsUnifiedVideo  bool   `json:"is_unified_video"`
		//User            struct {
		//	Pk                         int64  `json:"pk"`
		//	Username                   string `json:"username"`
		//	FullName                   string `json:"full_name"`
		//	IsPrivate                  bool   `json:"is_private"`
		//	ProfilePicUrl              string `json:"profile_pic_url"`
		//	ProfilePicId               string `json:"profile_pic_id"`
		//	IsVerified                 bool   `json:"is_verified"`
		//	FollowFrictionType         int    `json:"follow_friction_type"`
		//	HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		//	CanBoostPost               bool   `json:"can_boost_post"`
		//	CanSeeOrganicInsights      bool   `json:"can_see_organic_insights"`
		//	ShowInsightsTerms          bool   `json:"show_insights_terms"`
		//	ReelAutoArchive            string `json:"reel_auto_archive"`
		//	IsUnpublished              bool   `json:"is_unpublished"`
		//	AllowedCommenterType       string `json:"allowed_commenter_type"`
		//	HasHighlightReels          bool   `json:"has_highlight_reels"`
		//	InteropMessagingUserFbid   int64  `json:"interop_messaging_user_fbid"`
		//	FbidV2                     int64  `json:"fbid_v2"`
		//} `json:"user"`
		//CanViewerReshare                    bool          `json:"can_viewer_reshare"`
		//CaptionIsEdited                     bool          `json:"caption_is_edited"`
		//LikeAndViewCountsDisabled           bool          `json:"like_and_view_counts_disabled"`
		//CommercialityStatus                 string        `json:"commerciality_status"`
		//IsPaidPartnership                   bool          `json:"is_paid_partnership"`
		//IsVisualReplyCommenterNoticeEnabled bool          `json:"is_visual_reply_commenter_notice_enabled"`
		//CommentLikesEnabled                 bool          `json:"comment_likes_enabled"`
		//CommentThreadingEnabled             bool          `json:"comment_threading_enabled"`
		//HasMoreComments                     bool          `json:"has_more_comments"`
		//MaxNumVisiblePreviewComments        int           `json:"max_num_visible_preview_comments"`
		//PreviewComments                     []interface{} `json:"preview_comments"`
		//CanViewMorePreviewComments          bool          `json:"can_view_more_preview_comments"`
		//CommentCount                        int           `json:"comment_count"`
		//HideViewAllCommentEntrypoint        bool          `json:"hide_view_all_comment_entrypoint"`
		//ImageVersions2                      struct {
		//	Candidates []struct {
		//		Width        int    `json:"width"`
		//		Height       int    `json:"height"`
		//		Url          string `json:"url"`
		//		ScansProfile string `json:"scans_profile"`
		//	} `json:"candidates"`
		//	AdditionalCandidates struct {
		//		IgtvFirstFrame struct {
		//			Width        int    `json:"width"`
		//			Height       int    `json:"height"`
		//			Url          string `json:"url"`
		//			ScansProfile string `json:"scans_profile"`
		//		} `json:"igtv_first_frame"`
		//		FirstFrame struct {
		//			Width        int    `json:"width"`
		//			Height       int    `json:"height"`
		//			Url          string `json:"url"`
		//			ScansProfile string `json:"scans_profile"`
		//		} `json:"first_frame"`
		//	} `json:"additional_candidates"`
		//	AnimatedThumbnailSpritesheetInfoCandidates struct {
		//		Default struct {
		//			VideoLength                float64  `json:"video_length"`
		//			ThumbnailWidth             int      `json:"thumbnail_width"`
		//			ThumbnailHeight            int      `json:"thumbnail_height"`
		//			ThumbnailDuration          float64  `json:"thumbnail_duration"`
		//			SpriteUrls                 []string `json:"sprite_urls"`
		//			ThumbnailsPerRow           int      `json:"thumbnails_per_row"`
		//			TotalThumbnailNumPerSprite int      `json:"total_thumbnail_num_per_sprite"`
		//			MaxThumbnailsPerSprite     int      `json:"max_thumbnails_per_sprite"`
		//			SpriteWidth                int      `json:"sprite_width"`
		//			SpriteHeight               int      `json:"sprite_height"`
		//			RenderedWidth              int      `json:"rendered_width"`
		//			FileSizeKb                 int      `json:"file_size_kb"`
		//		} `json:"default"`
		//	} `json:"animated_thumbnail_spritesheet_info_candidates"`
		//} `json:"image_versions2"`
		//OriginalWidth         int           `json:"original_width"`
		//OriginalHeight        int           `json:"original_height"`
		//LikeCount             int           `json:"like_count"`
		//HasLiked              bool          `json:"has_liked"`
		//Likers                []interface{} `json:"likers"`
		//PhotoOfYou            bool          `json:"photo_of_you"`
		//CanSeeInsightsAsBrand bool          `json:"can_see_insights_as_brand"`
		//IsDashEligible        int           `json:"is_dash_eligible"`
		//VideoDashManifest     string        `json:"video_dash_manifest"`
		//VideoCodec            string        `json:"video_codec"`
		//NumberOfQualities     int           `json:"number_of_qualities"`
		//VideoVersions         []struct {
		//	Type   int    `json:"type"`
		//	Width  int    `json:"width"`
		//	Height int    `json:"height"`
		//	Url    string `json:"url"`
		//	Id     string `json:"id"`
		//} `json:"video_versions"`
		//HasAudio      bool    `json:"has_audio"`
		//VideoDuration float64 `json:"video_duration"`
		//PlayCount     int     `json:"play_count"`
		//Caption       struct {
		//	Pk              int64  `json:"pk"`
		//	UserId          int64  `json:"user_id"`
		//	Text            string `json:"text"`
		//	Type            int    `json:"type"`
		//	CreatedAt       int    `json:"created_at"`
		//	CreatedAtUtc    int    `json:"created_at_utc"`
		//	ContentType     string `json:"content_type"`
		//	Status          string `json:"status"`
		//	BitFlags        int    `json:"bit_flags"`
		//	DidReportAsSpam bool   `json:"did_report_as_spam"`
		//	ShareEnabled    bool   `json:"share_enabled"`
		//	User            struct {
		//		Pk                         int64  `json:"pk"`
		//		Username                   string `json:"username"`
		//		FullName                   string `json:"full_name"`
		//		IsPrivate                  bool   `json:"is_private"`
		//		ProfilePicUrl              string `json:"profile_pic_url"`
		//		ProfilePicId               string `json:"profile_pic_id"`
		//		IsVerified                 bool   `json:"is_verified"`
		//		FollowFrictionType         int    `json:"follow_friction_type"`
		//		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		//		CanBoostPost               bool   `json:"can_boost_post"`
		//		CanSeeOrganicInsights      bool   `json:"can_see_organic_insights"`
		//		ShowInsightsTerms          bool   `json:"show_insights_terms"`
		//		ReelAutoArchive            string `json:"reel_auto_archive"`
		//		IsUnpublished              bool   `json:"is_unpublished"`
		//		AllowedCommenterType       string `json:"allowed_commenter_type"`
		//		HasHighlightReels          bool   `json:"has_highlight_reels"`
		//		InteropMessagingUserFbid   int64  `json:"interop_messaging_user_fbid"`
		//		FbidV2                     int64  `json:"fbid_v2"`
		//	} `json:"user"`
		//	IsCovered          bool  `json:"is_covered"`
		//	MediaId            int64 `json:"media_id"`
		//	PrivateReplyStatus int   `json:"private_reply_status"`
		//} `json:"caption"`
		//CanViewerSave        bool   `json:"can_viewer_save"`
		//OrganicTrackingToken string `json:"organic_tracking_token"`
		//SharingFrictionInfo  struct {
		//	ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
		//	BloksAppUrl               interface{} `json:"bloks_app_url"`
		//} `json:"sharing_friction_info"`
		//CommentInformTreatment struct {
		//	ShouldHaveInformTreatment bool   `json:"should_have_inform_treatment"`
		//	Text                      string `json:"text"`
		//} `json:"comment_inform_treatment"`
		//ProductType               string      `json:"product_type"`
		//IsInProfileGrid           bool        `json:"is_in_profile_grid"`
		//ProfileGridControlEnabled bool        `json:"profile_grid_control_enabled"`
		//DeletedReason             int         `json:"deleted_reason"`
		//IntegrityReviewDecision   string      `json:"integrity_review_decision"`
		//MusicMetadata             interface{} `json:"music_metadata"`
		//ClipsMetadata             struct {
		//	MusicInfo         interface{} `json:"music_info"`
		//	OriginalSoundInfo struct {
		//		AudioAssetId           int64  `json:"audio_asset_id"`
		//		ProgressiveDownloadUrl string `json:"progressive_download_url"`
		//		DashManifest           string `json:"dash_manifest"`
		//		IgArtist               struct {
		//			Pk                 int64  `json:"pk"`
		//			Username           string `json:"username"`
		//			FullName           string `json:"full_name"`
		//			IsPrivate          bool   `json:"is_private"`
		//			ProfilePicUrl      string `json:"profile_pic_url"`
		//			ProfilePicId       string `json:"profile_pic_id"`
		//			IsVerified         bool   `json:"is_verified"`
		//			FollowFrictionType int    `json:"follow_friction_type"`
		//		} `json:"ig_artist"`
		//		ShouldMuteAudio    bool   `json:"should_mute_audio"`
		//		OriginalMediaId    int64  `json:"original_media_id"`
		//		HideRemixing       bool   `json:"hide_remixing"`
		//		DurationInMs       int    `json:"duration_in_ms"`
		//		TimeCreated        int    `json:"time_created"`
		//		OriginalAudioTitle string `json:"original_audio_title"`
		//		ConsumptionInfo    struct {
		//			IsBookmarked              bool        `json:"is_bookmarked"`
		//			ShouldMuteAudioReason     string      `json:"should_mute_audio_reason"`
		//			IsTrendingInClips         bool        `json:"is_trending_in_clips"`
		//			ShouldMuteAudioReasonType interface{} `json:"should_mute_audio_reason_type"`
		//		} `json:"consumption_info"`
		//		AllowCreatorToRename           bool          `json:"allow_creator_to_rename"`
		//		CanRemixBeSharedToFb           bool          `json:"can_remix_be_shared_to_fb"`
		//		FormattedClipsMediaCount       interface{}   `json:"formatted_clips_media_count"`
		//		AudioParts                     []interface{} `json:"audio_parts"`
		//		IsExplicit                     bool          `json:"is_explicit"`
		//		OriginalAudioSubtype           string        `json:"original_audio_subtype"`
		//		IsAudioAutomaticallyAttributed bool          `json:"is_audio_automatically_attributed"`
		//	} `json:"original_sound_info"`
		//	AudioType        string      `json:"audio_type"`
		//	MusicCanonicalId string      `json:"music_canonical_id"`
		//	FeaturedLabel    interface{} `json:"featured_label"`
		//	MashupInfo       struct {
		//		MashupsAllowed                      bool        `json:"mashups_allowed"`
		//		CanToggleMashupsAllowed             bool        `json:"can_toggle_mashups_allowed"`
		//		HasBeenMashedUp                     bool        `json:"has_been_mashed_up"`
		//		FormattedMashupsCount               interface{} `json:"formatted_mashups_count"`
		//		OriginalMedia                       interface{} `json:"original_media"`
		//		NonPrivacyFilteredMashupsMediaCount interface{} `json:"non_privacy_filtered_mashups_media_count"`
		//	} `json:"mashup_info"`
		//	NuxInfo                   interface{} `json:"nux_info"`
		//	ViewerInteractionSettings interface{} `json:"viewer_interaction_settings"`
		//	BrandedContentTagInfo     struct {
		//		CanAddTag bool `json:"can_add_tag"`
		//	} `json:"branded_content_tag_info"`
		//	ShoppingInfo        interface{} `json:"shopping_info"`
		//	AdditionalAudioInfo struct {
		//		AdditionalAudioUsername interface{} `json:"additional_audio_username"`
		//		AudioReattributionInfo  struct {
		//			ShouldAllowRestore bool `json:"should_allow_restore"`
		//		} `json:"audio_reattribution_info"`
		//	} `json:"additional_audio_info"`
		//	IsSharedToFb            bool        `json:"is_shared_to_fb"`
		//	BreakingContentInfo     interface{} `json:"breaking_content_info"`
		//	ChallengeInfo           interface{} `json:"challenge_info"`
		//	ReelsOnTheRiseInfo      interface{} `json:"reels_on_the_rise_info"`
		//	BreakingCreatorInfo     interface{} `json:"breaking_creator_info"`
		//	AssetRecommendationInfo interface{} `json:"asset_recommendation_info"`
		//	ContextualHighlightInfo interface{} `json:"contextual_highlight_info"`
		//	ClipsCreationEntryPoint string      `json:"clips_creation_entry_point"`
		//	AudioRankingInfo        struct {
		//		BestAudioClusterId string `json:"best_audio_cluster_id"`
		//	} `json:"audio_ranking_info"`
		//	TemplateInfo interface{} `json:"template_info"`
		//} `json:"clips_metadata"`
		//MediaCroppingInfo struct {
		//	SquareCrop struct {
		//		CropLeft   float64 `json:"crop_left"`
		//		CropRight  float64 `json:"crop_right"`
		//		CropTop    float64 `json:"crop_top"`
		//		CropBottom float64 `json:"crop_bottom"`
		//	} `json:"square_crop"`
		//} `json:"media_cropping_info"`
		//LoggingInfoToken interface{} `json:"logging_info_token"`
	} `json:"media"`
	UploadId string `json:"upload_id"`
}

func (this *UserOperate) ConfigureToClips(mediaInfo *RawVideoMedia) (*RespConfigureClips, error) {
	var cameraPosition int
	var sourceType string
	if mediaInfo.From == FromCamera {
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
		"_uuid":             this.inst.AccountInfo.Device.DeviceID,
		"_uid":              fmt.Sprintf("%d", this.inst.ID),
		"internal_features": "clips_format,clips_launch",
		"source_type":       sourceType,
		"is_clips_edited":   "0",
		"camera_session_id": common.GenString(common.CharSet_16_Num, 32),
		"additional_audio_info": map[string]interface{}{
			"has_voiceover_attribution": "0",
		},
		"device_id":                   this.inst.AccountInfo.Device.DeviceID,
		"client_timestamp":            fmt.Sprintf("%d", time.Now().Unix()),
		"effect_ids":                  []int{},
		"waterfall_id":                mediaInfo.Waterfall,
		"caption":                     mediaInfo.Caption,
		"camera_entry_point":          this.cameraEntryPoint,
		"text_overlay":                []int{},
		"clips_share_preview_to_feed": "1",
		"upload_id":                   mediaInfo.UploadId,
		"sticker_ids":                 []int{},
		"timezone_offset":             this.inst.AccountInfo.Location.Timezone,
		"capture_type":                "clips_v2",
		"clips_audio_metadata": map[string]interface{}{
			"original": map[string]interface{}{
				"volume_level": 1,
			},
			"original_audio_title": mediaInfo.AudioTitle,
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
				"source":              mediaInfo.From,
				"camera_position":     cameraPosition,
				"duration_ms":         mediaInfo.Duration,
			}},
		},
		"overlay_data": []int{},
	}

	resp := &RespConfigureClips{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlConfigureToClips,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

func (this *UserOperate) ClipsInfoForCreation() error {
	resp := &BaseApiResp{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlClipsInfoForCreation,
	}, resp)

	err = resp.CheckError(err)
	return err
}

type RespClipsAssets struct {
	BaseApiResp
	StaticStickers []struct {
		Id       interface{} `json:"id"`
		Stickers []struct {
			Id                  interface{} `json:"id"`
			Type                string      `json:"type,omitempty"`
			Name                string      `json:"name,omitempty"`
			Text                string      `json:"text,omitempty"`
			FontSize            int         `json:"font_size,omitempty"`
			TextColor           string      `json:"text_color,omitempty"`
			TextBackgroundColor string      `json:"text_background_color,omitempty"`
			TextBackgroundAlpha int         `json:"text_background_alpha,omitempty"`
			TapStateStrId       string      `json:"tap_state_str_id,omitempty"`
			ImageUrl            string      `json:"image_url,omitempty"`
			ImageWidthRatio     float64     `json:"image_width_ratio,omitempty"`
			TrayImageWidthRatio float64     `json:"tray_image_width_ratio,omitempty"`
			ImageWidth          int         `json:"image_width,omitempty"`
			ImageHeight         int         `json:"image_height,omitempty"`
		} `json:"stickers"`
		Keywords          []string    `json:"keywords"`
		IncludeInRecent   bool        `json:"include_in_recent,omitempty"`
		HasAttribution    interface{} `json:"has_attribution"`
		AvailableInDirect bool        `json:"available_in_direct,omitempty"`
	} `json:"static_stickers"`
	Version        int `json:"version"`
	ComposerConfig struct {
		SwipeUpUrls                     bool `json:"swipe_up_urls"`
		FelixLinks                      bool `json:"felix_links"`
		TotalArEffects                  int  `json:"total_ar_effects"`
		ProfileShopLinks                bool `json:"profile_shop_links"`
		ShoppingLinkMoreOptions         bool `json:"shopping_link_more_options"`
		ShoppingCollectionLinks         bool `json:"shopping_collection_links"`
		ShoppingProductCollectionLinks  bool `json:"shopping_product_collection_links"`
		ShoppingProductLinks            bool `json:"shopping_product_links"`
		ShoppingMultiProductLinks       bool `json:"shopping_multi_product_links"`
		ShoppingMultiProductMaxProducts int  `json:"shopping_multi_product_max_products"`
	} `json:"composer_config"`
}

func (this *UserOperate) ClipsAssets(latitude float32, longitude float32) (*RespClipsAssets, error) {
	params := map[string]interface{}{
		"verticalAccuracy":   "10.000000",
		"speed":              "-1.000000",
		"_uuid":              this.inst.AccountInfo.Device.DeviceID,
		"timezone_offset":    this.inst.AccountInfo.Location.Timezone,
		"horizontalAccuracy": "65.000000",
		"alt":                "36.008301",
		"_uid":               fmt.Sprintf("%d", this.inst.ID),
		"type":               "static_stickers",
		"lng":                longitude,
		"lat":                latitude,
	}
	resp := &RespClipsAssets{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlClipsAssets,
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type RespVerifyAudioTitle struct {
	BaseApiResp
	IsValid bool `json:"is_valid"`
}

func (this *UserOperate) VerifyOriginalAudioTitle(originalAudioName string) (*RespVerifyAudioTitle, error) {
	params := map[string]interface{}{
		"_uuid":               this.inst.AccountInfo.Device.DeviceID,
		"original_audio_name": originalAudioName,
	}
	resp := &RespVerifyAudioTitle{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlVerifyOriginalAudioTitle,
		IsPost:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return resp, err
}

type MeasuredFrames struct {
	Ssim      float64 `json:"ssim"`
	Timestamp float64 `json:"timestamp"`
}
type QualityInfo struct {
	OriginalVideoCodec       string           `json:"original_video_codec,omitempty"`
	EncodedVideoCodec        string           `json:"encoded_video_codec,omitempty"`
	OriginalColorPrimaries   string           `json:"original_color_primaries,omitempty"`
	OriginalWidth            int              `json:"original_width,omitempty"`
	OriginalFrameRate        float64          `json:"original_frame_rate,omitempty"`
	OriginalTransferFunction string           `json:"original_transfer_function,omitempty"`
	EncodedHeight            int              `json:"encoded_height,omitempty"`
	OriginalBitRate          float64          `json:"original_bit_rate,omitempty"`
	EncodedColorPrimaries    string           `json:"encoded_color_primaries,omitempty"`
	OriginalHeight           int              `json:"original_height,omitempty"`
	EncodedBitRate           float64          `json:"encoded_bit_rate,omitempty"`
	EncodedFrameRate         float64          `json:"encoded_frame_rate,omitempty"`
	EncodedYcbcrMatrix       string           `json:"encoded_ycbcr_matrix,omitempty"`
	OriginalYcbcrMatrix      string           `json:"original_ycbcr_matrix,omitempty"`
	EncodedWidth             int              `json:"encoded_width,omitempty"`
	MeasuredFrames           []MeasuredFrames `json:"measured_frames,omitempty"`
	EncodedTransferFunction  string           `json:"encoded_transfer_function,omitempty"`
}

func (this *UserOperate) UpdateVideoWithQualityInfo(uploadID string, qualityInfo *QualityInfo) error {
	qualityInfoStr, _ := json.Marshal(qualityInfo)
	params := map[string]interface{}{
		"_uuid":        this.inst.AccountInfo.Device.DeviceID,
		"_uid":         fmt.Sprintf("%d", this.inst.ID),
		"quality_info": common.B2s(qualityInfoStr),
		"upload_id":    uploadID,
	}
	resp := &BaseApiResp{}
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
		"_uuid":            this.inst.AccountInfo.Device.DeviceID,
		"_uid":             this.inst.ID,
		"user_id":          userID,
		"device_id":        this.inst.AccountInfo.Device.DeviceID,
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

type AddCommentParams struct {
	ParentCommentId string
	UserName        string

	LoggingInfoToken string
	MediaId          string
	CommentText      string
}

var AddCommentNavChain = []string{"IGSundialHomeViewController:clips_viewer_clips_tab:2,IGCommentThreadV2ViewController:comments_v2_clips_viewer_clips_tab:3"}
var AddSubCommentNavChain = []string{"IGSundialHomeViewController:clips_viewer_clips_tab:2"}

type RespCheckOffensiveComment struct {
	BaseApiResp
	IsOffensive     bool    `json:"is_offensive"`
	TextLanguage    string  `json:"text_language"`
	BullyClassifier float64 `json:"bully_classifier"`
}

type RespAddComment struct {
	BaseApiResp
	Comment struct {
		ContentType  string             `json:"content_type"`
		User         AddCommentRespUser `json:"user"`
		Pk           int64              `json:"pk"`
		Text         string             `json:"text"`
		Type         int                `json:"type"`
		CreatedAt    int                `json:"created_at"`
		CreatedAtUtc int                `json:"created_at_utc"`
		MediaId      int64              `json:"media_id"`
		Status       string             `json:"status"`
		ShareEnabled bool               `json:"share_enabled"`
	} `json:"comment"`
}

func (this *UserOperate) AddComment(params *AddCommentParams) error {
	checkParams := map[string]interface{}{
		"media_id":           params.MediaId,
		"_uuid":              this.inst.AccountInfo.Device,
		"_uid":               fmt.Sprintf("%d", this.inst.ID),
		"comment_session_id": this.commentSessionId,
	}

	commentParam := map[string]interface{}{
		"delivery_class":          "organic",
		"logging_info_token":      params.LoggingInfoToken,
		"_uuid":                   this.inst.AccountInfo.Device,
		"_uid":                    fmt.Sprintf("%d", this.inst.ID),
		"idempotence_token":       common.GenString(common.CharSet_16_Num, 32),
		"is_carousel_bumped_post": "false",
		"carousel_index":          "0",
		"container_module":        "comments_v2_clips_viewer_clips_tab",
	}
	if params.ParentCommentId != "" {
		text := "@" + params.UserName + " " + params.CommentText
		checkParams["comment_text"] = text
		commentParam["comment_text"] = text
		commentParam["parent_comment_id"] = params.ParentCommentId
		commentParam["replied_to_comment_id"] = params.ParentCommentId
		commentParam["nav_chain"] = AddSubCommentNavChain[common.GenNumber(0, len(AddSubCommentNavChain))]
	} else {
		checkParams["comment_text"] = params.CommentText
		commentParam["comment_text"] = params.CommentText
		commentParam["nav_chain"] = AddCommentNavChain[common.GenNumber(0, len(AddCommentNavChain))]
	}
	checkResp := &RespCheckOffensiveComment{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:         true,
		Signed:         true,
		ApiPath:        checkOffensiveComment,
		HeaderSequence: LoginHeaderMap[checkOffensiveComment],
		Query:          checkParams,
	}, checkResp)

	err = checkResp.CheckError(err)
	if err != nil {
		return err
	}

	addResp := &RespAddComment{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:         true,
		Signed:         true,
		ApiPath:        fmt.Sprintf(urlCommentAdd, params.MediaId),
		HeaderSequence: LoginHeaderMap[urlCommentAdd],
		Query:          commentParam,
	}, addResp)

	err = addResp.CheckError(err)
	return err
}
