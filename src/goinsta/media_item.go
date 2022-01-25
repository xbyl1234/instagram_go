package goinsta

import (
	"makemoney/common/log"
)

type MediaType int

var (
	MediaType_Photo    MediaType = 1
	MediaType_Video    MediaType = 2
	MediaType_Carousel MediaType = 8
)

type Item struct {
	Inst *Instagram `bson:"-"`
	//Comments *Comments  `json:"-" bson:"-"`

	CanSeeInsightsAsBrand      bool   `json:"can_see_insights_as_brand" bson:"can_see_insights_as_brand"`
	CanViewMorePreviewComments bool   `json:"can_view_more_preview_comments" bson:"can_view_more_preview_comments"`
	CommercialityStatus        string `json:"commerciality_status" bson:"commerciality_status"`
	DeletedReason              int    `json:"deleted_reason" bson:"deleted_reason"`
	FundraiserTag              struct {
		HasStandaloneFundraiser bool `json:"has_standalone_fundraiser" bson:"has_standalone_fundraiser"`
	} `json:"fundraiser_tag" bson:"fundraiser_tag"`
	HideViewAllCommentEntrypoint bool   `json:"hide_view_all_comment_entrypoint" bson:"hide_view_all_comment_entrypoint"`
	IntegrityReviewDecision      string `json:"integrity_review_decision" bson:"integrity_review_decision"`
	IsCommercial                 bool   `json:"is_commercial" bson:"is_commercial"`
	IsInProfileGrid              bool   `json:"is_in_profile_grid" bson:"is_in_profile_grid"`
	IsPaidPartnership            bool   `json:"is_paid_partnership" bson:"is_paid_partnership"`
	IsUnifiedVideo               bool   `json:"is_unified_video" bson:"is_unified_video"`
	LikeAndViewCountsDisabled    bool   `json:"like_and_view_counts_disabled" bson:"like_and_view_counts_disabled"`
	NextMaxId                    int64  `json:"next_max_id" bson:"next_max_id"`
	ProductType                  string `json:"product_type" bson:"product_type"`
	ProfileGridControlEnabled    bool   `json:"profile_grid_control_enabled" bson:"profile_grid_control_enabled"`

	TakenAt          int64   `json:"taken_at" bson:"taken_at"`
	Pk               int64   `json:"pk" bson:"pk"`
	ID               string  `json:"id" bson:"id"`
	CommentsDisabled bool    `json:"comments_disabled" bson:"comments_disabled"`
	DeviceTimestamp  int64   `json:"device_timestamp" bson:"device_timestamp"`
	MediaType        int     `json:"media_type" bson:"media_type"`
	Code             string  `json:"code" bson:"code"`
	ClientCacheKey   string  `json:"client_cache_key" bson:"client_cache_key"`
	FilterType       int     `json:"filter_type" bson:"filter_type"`
	CarouselParentID string  `json:"carousel_parent_id" bson:"carousel_parent_id"`
	CarouselMedia    []Item  `json:"carousel_media,omitempty" bson:"carousel_media"`
	User             User    `json:"user" bson:"user"`
	CanViewerReshare bool    `json:"can_viewer_reshare" bson:"can_viewer_reshare"`
	Caption          Caption `json:"caption" bson:"caption"`
	CaptionIsEdited  bool    `json:"caption_is_edited" bson:"caption_is_edited"`
	Likes            int     `json:"like_count" bson:"likes"`
	HasLiked         bool    `json:"has_liked" bson:"has_liked"`
	// Toplikers can be `string` or `[]string`.
	// Use TopLikers function instead of getting it directly.
	Toplikers                    interface{} `json:"top_likers" bson:"toplikers"`
	Likers                       []User      `json:"likers" bson:"likers"`
	CommentLikesEnabled          bool        `json:"comment_likes_enabled" bson:"comment_likes_enabled"`
	CommentThreadingEnabled      bool        `json:"comment_threading_enabled" bson:"comment_threading_enabled"`
	HasMoreComments              bool        `json:"has_more_comments" bson:"has_more_comments"`
	MaxNumVisiblePreviewComments int         `json:"max_num_visible_preview_comments" bson:"max_num_visible_preview_comments"`
	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	// Use PreviewComments function instead of getting it directly.
	Previewcomments interface{} `json:"preview_comments,omitempty" bson:"previewcomments"`
	CommentCount    int         `json:"comment_count" bson:"comment_count"`
	PhotoOfYou      bool        `json:"photo_of_you" bson:"photo_of_you"`
	// Tags are tagged people in photo
	Tags struct {
		In []UserTag `json:"in" bson:"in"`
	} `json:"usertags,omitempty" bson:"tags"`
	FbUserTags           UserTag `json:"fb_user_tags" bson:"fb_user_tags"`
	CanViewerSave        bool    `json:"can_viewer_save" bson:"can_viewer_save"`
	OrganicTrackingToken string  `json:"organic_tracking_token" bson:"organic_tracking_token"`
	// Images contains URL images in different versions.
	// Version = quality.
	Images          Images   `json:"image_versions2,omitempty" bson:"images"`
	OriginalWidth   int      `json:"original_width,omitempty" bson:"original_width"`
	OriginalHeight  int      `json:"original_height,omitempty" bson:"original_height"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty" bson:"imported_taken_at"`
	Location        Location `json:"location,omitempty" bson:"location"`
	Lat             float64  `json:"lat,omitempty" bson:"lat"`
	Lng             float64  `json:"lng,omitempty" bson:"lng"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty" bson:"videos"`
	HasAudio          bool    `json:"has_audio,omitempty" bson:"has_audio"`
	VideoDuration     float64 `json:"video_duration,omitempty" bson:"video_duration"`
	ViewCount         float64 `json:"view_count,omitempty" bson:"view_count"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty" bson:"is_dash_eligible"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty" bson:"video_dash_manifest"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty" bson:"number_of_qualities"`

	// Only for stories
	//StoryEvents              []interface{}      `json:"story_events"`
	//StoryHashtags            []interface{}      `json:"story_hashtags"`
	//StoryPolls               []interface{}      `json:"story_polls"`
	//StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	//StorySoundOn             []interface{}      `json:"story_sound_on"`
	//CreativeConfig           interface{}        `json:"creative_config"`
	//StoryLocations           []interface{}      `json:"story_locations"`
	//StorySliders             []interface{}      `json:"story_sliders"`
	//StoryQuestions           []interface{}      `json:"story_questions"`
	//StoryProductItems        []interface{}      `json:"story_product_items"`
	//StoryCTA                 []StoryCTA         `json:"story_cta"`
	//ReelMentions             []StoryReelMention `json:"reel_mentions"`
	//SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	//ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	//HasSharedToFb            int64              `json:"has_shared_to_fb"`
	//Mentions                 []Mentions
	//Audience                 string `json:"audience,omitempty"`
	//StoryMusicStickers       []struct {
	//	X              float64 `json:"x"`
	//	Y              float64 `json:"y"`
	//	Z              int     `json:"z"`
	//	Width          float64 `json:"width"`
	//	Height         float64 `json:"height"`
	//	Rotation       float64 `json:"rotation"`
	//	IsPinned       int     `json:"is_pinned"`
	//	IsHidden       int     `json:"is_hidden"`
	//	IsSticker      int     `json:"is_sticker"`
	//	MusicAssetInfo struct {
	//		ID                       string `json:"id"`
	//		Title                    string `json:"title"`
	//		Subtitle                 string `json:"subtitle"`
	//		DisplayArtist            string `json:"display_artist"`
	//		CoverArtworkURI          string `json:"cover_artwork_uri"`
	//		CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
	//		ProgressiveDownloadURL   string `json:"progressive_download_url"`
	//		HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
	//		IsExplicit               bool   `json:"is_explicit"`
	//		DashManifest             string `json:"dash_manifest"`
	//		HasLyrics                bool   `json:"has_lyrics"`
	//		AudioAssetID             string `json:"audio_asset_id"`
	//		IgArtist                 struct {
	//			Pk            int    `json:"pk"`
	//			Username      string `json:"username"`
	//			FullName      string `json:"full_name"`
	//			IsPrivate     bool   `json:"is_private"`
	//			ProfilePicURL string `json:"profile_pic_url"`
	//			ProfilePicID  string `json:"profile_pic_id"`
	//			IsVerified    bool   `json:"is_verified"`
	//		} `json:"ig_artist"`
	//		PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
	//		ShouldMuteAudio          bool   `json:"should_mute_audio"`
	//		ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
	//		OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
	//		AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
	//	} `json:"music_asset_info"`
	//} `json:"story_music_stickers,omitempty"`
}

func (this *Item) GetMediaType() MediaType {
	switch this.MediaType {
	case 1:
		return MediaType_Photo
	case 2:
		return MediaType_Video
	case 8:
		return MediaType_Carousel
	}
	log.Error("GetMediaType error: %v", this.MediaType)
	return MediaType_Photo
}
func (this *Item) SetAccount(inst *Instagram) {
	this.Inst = inst
}

func (this *Item) GetComments() *Comments {
	return &Comments{media: this, Inst: this.Inst, MediaID: this.ID, HasMore: true}
}
