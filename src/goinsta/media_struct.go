package goinsta

type MediaUser struct {
	Pk               int64  `json:"pk"`
	Username         string `json:"username"`
	FullName         string `json:"full_name"`
	IsPrivate        bool   `json:"is_private"`
	ProfilePicUrl    string `json:"profile_pic_url"`
	ProfilePicId     string `json:"profile_pic_id"`
	FriendshipStatus struct {
		Following       bool `json:"following"`
		OutgoingRequest bool `json:"outgoing_request"`
		IsBestie        bool `json:"is_bestie"`
		IsRestricted    bool `json:"is_restricted"`
		IsFeedFavorite  bool `json:"is_feed_favorite"`
	} `json:"friendship_status"`
	IsVerified         bool `json:"is_verified"`
	FollowFrictionType int  `json:"follow_friction_type"`
	GrowthFrictionInfo struct {
		HasActiveInterventions bool `json:"has_active_interventions"`
		Interventions          struct {
		} `json:"interventions"`
	} `json:"growth_friction_info"`
	AccountBadges              []interface{} `json:"account_badges"`
	HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture"`
	IsUnpublished              bool          `json:"is_unpublished"`
	IsFavorite                 bool          `json:"is_favorite"`
	HasHighlightReels          bool          `json:"has_highlight_reels"`
	HasPrimaryCountryInFeed    bool          `json:"has_primary_country_in_feed"`
	HasPrimaryCountryInProfile bool          `json:"has_primary_country_in_profile"`
}
type CommentUser struct {
	Pk                 int64  `json:"pk"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	IsPrivate          bool   `json:"is_private"`
	ProfilePicUrl      string `json:"profile_pic_url"`
	ProfilePicId       string `json:"profile_pic_id,omitempty"`
	IsVerified         bool   `json:"is_verified"`
	FollowFrictionType int    `json:"follow_friction_type"`
	GrowthFrictionInfo struct {
		HasActiveInterventions bool `json:"has_active_interventions"`
		Interventions          struct {
		} `json:"interventions"`
	} `json:"growth_friction_info"`
	AccountBadges []interface{} `json:"account_badges"`
}

type AddCommentRespUser struct {
	Pk                 int64  `json:"pk"`
	Username           string `json:"username"`
	FullName           string `json:"full_name"`
	IsPrivate          bool   `json:"is_private"`
	ProfilePicUrl      string `json:"profile_pic_url"`
	IsVerified         bool   `json:"is_verified"`
	FollowFrictionType int    `json:"follow_friction_type"`
	GrowthFrictionInfo struct {
		HasActiveInterventions bool `json:"has_active_interventions"`
		Interventions          struct {
		} `json:"interventions"`
	} `json:"growth_friction_info"`
	AccountBadges              []interface{} `json:"account_badges"`
	HasAnonymousProfilePicture bool          `json:"has_anonymous_profile_picture"`
	ReelAutoArchive            string        `json:"reel_auto_archive"`
	AllowedCommenterType       string        `json:"allowed_commenter_type"`
	HasHighlightReels          bool          `json:"has_highlight_reels"`
	InteropMessagingUserFbid   int64         `json:"interop_messaging_user_fbid"`
	HasPrimaryCountryInFeed    bool          `json:"has_primary_country_in_feed"`
	HasPrimaryCountryInProfile bool          `json:"has_primary_country_in_profile"`
	FbidV2                     int64         `json:"fbid_v2"`
}

type PreviewComment struct {
	Pk                 int64       `json:"pk"`
	UserId             int64       `json:"user_id"`
	Text               string      `json:"text"`
	Type               int         `json:"type"`
	CreatedAt          int         `json:"created_at"`
	CreatedAtUtc       int         `json:"created_at_utc"`
	ContentType        string      `json:"content_type"`
	Status             string      `json:"status"`
	BitFlags           int         `json:"bit_flags"`
	DidReportAsSpam    bool        `json:"did_report_as_spam"`
	ShareEnabled       bool        `json:"share_enabled"`
	User               CommentUser `json:"user"`
	IsCovered          bool        `json:"is_covered"`
	MediaId            int64       `json:"media_id"`
	PrivateReplyStatus int         `json:"private_reply_status"`
	HasTranslation     bool        `json:"has_translation,omitempty"`
	IsPinned           bool        `json:"is_pinned,omitempty"`
}
type ImageVersion2 struct {
	Candidates []struct {
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		Url          string `json:"url"`
		ScansProfile string `json:"scans_profile"`
	} `json:"candidates"`
	AdditionalCandidates struct {
		IgtvFirstFrame struct {
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			Url          string `json:"url"`
			ScansProfile string `json:"scans_profile"`
		} `json:"igtv_first_frame"`
		FirstFrame struct {
			Width        int    `json:"width"`
			Height       int    `json:"height"`
			Url          string `json:"url"`
			ScansProfile string `json:"scans_profile"`
		} `json:"first_frame"`
	} `json:"additional_candidates"`
	AnimatedThumbnailSpritesheetInfoCandidates struct {
		Default struct {
			VideoLength                float64  `json:"video_length"`
			ThumbnailWidth             int      `json:"thumbnail_width"`
			ThumbnailHeight            int      `json:"thumbnail_height"`
			ThumbnailDuration          float64  `json:"thumbnail_duration"`
			SpriteUrls                 []string `json:"sprite_urls"`
			ThumbnailsPerRow           int      `json:"thumbnails_per_row"`
			TotalThumbnailNumPerSprite int      `json:"total_thumbnail_num_per_sprite"`
			MaxThumbnailsPerSprite     int      `json:"max_thumbnails_per_sprite"`
			SpriteWidth                int      `json:"sprite_width"`
			SpriteHeight               int      `json:"sprite_height"`
			RenderedWidth              int      `json:"rendered_width"`
			FileSizeKb                 int      `json:"file_size_kb"`
		} `json:"default"`
	} `json:"animated_thumbnail_spritesheet_info_candidates"`
}

type UserTag struct {
	In []struct {
		User                  MediaUser   `json:"user"`
		Position              []float64   `json:"position"`
		StartTimeInVideoInSec interface{} `json:"start_time_in_video_in_sec"`
		DurationInVideoInSec  interface{} `json:"duration_in_video_in_sec"`
	} `json:"in"`
}

type VideoVersion struct {
	Type   int    `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
	Id     string `json:"id"`
}

type MediaCaption struct {
	Pk                 int64       `json:"pk"`
	UserId             int64       `json:"user_id"`
	Text               string      `json:"text"`
	Type               int         `json:"type"`
	CreatedAt          int         `json:"created_at"`
	CreatedAtUtc       int         `json:"created_at_utc"`
	ContentType        string      `json:"content_type"`
	Status             string      `json:"status"`
	BitFlags           int         `json:"bit_flags"`
	DidReportAsSpam    bool        `json:"did_report_as_spam"`
	ShareEnabled       bool        `json:"share_enabled"`
	User               CommentUser `json:"user"`
	IsCovered          bool        `json:"is_covered"`
	MediaId            int64       `json:"media_id"`
	PrivateReplyStatus int         `json:"private_reply_status"`
	HasTranslation     bool        `json:"has_translation,omitempty"`
}

type ClipsMetadata struct {
	MusicInfo *struct {
		MusicAssetInfo struct {
			AudioClusterId                  string      `json:"audio_cluster_id"`
			Id                              string      `json:"id"`
			Title                           string      `json:"title"`
			Subtitle                        string      `json:"subtitle"`
			DisplayArtist                   string      `json:"display_artist"`
			ArtistId                        string      `json:"artist_id"`
			CoverArtworkUri                 string      `json:"cover_artwork_uri"`
			CoverArtworkThumbnailUri        string      `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadUrl          string      `json:"progressive_download_url"`
			ReactiveAudioDownloadUrl        interface{} `json:"reactive_audio_download_url"`
			FastStartProgressiveDownloadUrl string      `json:"fast_start_progressive_download_url"`
			Web30SPreviewDownloadUrl        interface{} `json:"web_30s_preview_download_url"`
			HighlightStartTimesInMs         []int       `json:"highlight_start_times_in_ms"`
			IsExplicit                      bool        `json:"is_explicit"`
			DashManifest                    interface{} `json:"dash_manifest"`
			HasLyrics                       bool        `json:"has_lyrics"`
			AudioAssetId                    string      `json:"audio_asset_id"`
			DurationInMs                    int         `json:"duration_in_ms"`
			DarkMessage                     interface{} `json:"dark_message"`
			AllowsSaving                    bool        `json:"allows_saving"`
			TerritoryValidityPeriods        struct {
			} `json:"territory_validity_periods"`
		} `json:"music_asset_info"`
		MusicConsumptionInfo struct {
			IgArtist struct {
				Pk                 int64  `json:"pk"`
				Username           string `json:"username"`
				FullName           string `json:"full_name"`
				IsPrivate          bool   `json:"is_private"`
				ProfilePicUrl      string `json:"profile_pic_url"`
				ProfilePicId       string `json:"profile_pic_id"`
				IsVerified         bool   `json:"is_verified"`
				FollowFrictionType int    `json:"follow_friction_type"`
				GrowthFrictionInfo struct {
					HasActiveInterventions bool `json:"has_active_interventions"`
					Interventions          struct {
					} `json:"interventions"`
				} `json:"growth_friction_info"`
				AccountBadges []interface{} `json:"account_badges"`
			} `json:"ig_artist"`
			PlaceholderProfilePicUrl    string      `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio             bool        `json:"should_mute_audio"`
			ShouldMuteAudioReason       string      `json:"should_mute_audio_reason"`
			ShouldMuteAudioReasonType   interface{} `json:"should_mute_audio_reason_type"`
			IsBookmarked                bool        `json:"is_bookmarked"`
			OverlapDurationInMs         int         `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs     int         `json:"audio_asset_start_time_in_ms"`
			AllowMediaCreationWithMusic bool        `json:"allow_media_creation_with_music"`
			IsTrendingInClips           bool        `json:"is_trending_in_clips"`
			FormattedClipsMediaCount    interface{} `json:"formatted_clips_media_count"`
			StreamingServices           interface{} `json:"streaming_services"`
			DisplayLabels               interface{} `json:"display_labels"`
		} `json:"music_consumption_info"`
		PushBlockingTest interface{} `json:"push_blocking_test"`
	} `json:"music_info"`
	OriginalSoundInfo *struct {
		AudioAssetId           int64  `json:"audio_asset_id"`
		ProgressiveDownloadUrl string `json:"progressive_download_url"`
		DashManifest           string `json:"dash_manifest"`
		IgArtist               struct {
			Pk                 int64  `json:"pk"`
			Username           string `json:"username"`
			FullName           string `json:"full_name"`
			IsPrivate          bool   `json:"is_private"`
			ProfilePicUrl      string `json:"profile_pic_url"`
			ProfilePicId       string `json:"profile_pic_id"`
			IsVerified         bool   `json:"is_verified"`
			FollowFrictionType int    `json:"follow_friction_type"`
			GrowthFrictionInfo struct {
				HasActiveInterventions bool `json:"has_active_interventions"`
				Interventions          struct {
				} `json:"interventions"`
			} `json:"growth_friction_info"`
			AccountBadges []interface{} `json:"account_badges"`
		} `json:"ig_artist"`
		ShouldMuteAudio    bool   `json:"should_mute_audio"`
		OriginalMediaId    int64  `json:"original_media_id"`
		HideRemixing       bool   `json:"hide_remixing"`
		DurationInMs       int    `json:"duration_in_ms"`
		TimeCreated        int    `json:"time_created"`
		OriginalAudioTitle string `json:"original_audio_title"`
		ConsumptionInfo    struct {
			IsBookmarked              bool        `json:"is_bookmarked"`
			ShouldMuteAudioReason     string      `json:"should_mute_audio_reason"`
			IsTrendingInClips         bool        `json:"is_trending_in_clips"`
			ShouldMuteAudioReasonType interface{} `json:"should_mute_audio_reason_type"`
		} `json:"consumption_info"`
		AllowCreatorToRename           bool          `json:"allow_creator_to_rename"`
		CanRemixBeSharedToFb           bool          `json:"can_remix_be_shared_to_fb"`
		FormattedClipsMediaCount       interface{}   `json:"formatted_clips_media_count"`
		AudioParts                     []interface{} `json:"audio_parts"`
		IsExplicit                     bool          `json:"is_explicit"`
		OriginalAudioSubtype           string        `json:"original_audio_subtype"`
		IsAudioAutomaticallyAttributed bool          `json:"is_audio_automatically_attributed"`
	} `json:"original_sound_info"`
	AudioType        string      `json:"audio_type"`
	MusicCanonicalId string      `json:"music_canonical_id"`
	FeaturedLabel    interface{} `json:"featured_label"`
	MashupInfo       struct {
		MashupsAllowed                      bool        `json:"mashups_allowed"`
		CanToggleMashupsAllowed             bool        `json:"can_toggle_mashups_allowed"`
		HasBeenMashedUp                     bool        `json:"has_been_mashed_up"`
		FormattedMashupsCount               interface{} `json:"formatted_mashups_count"`
		OriginalMedia                       interface{} `json:"original_media"`
		NonPrivacyFilteredMashupsMediaCount *int        `json:"non_privacy_filtered_mashups_media_count"`
	} `json:"mashup_info"`
	NuxInfo                   interface{} `json:"nux_info"`
	ViewerInteractionSettings interface{} `json:"viewer_interaction_settings"`
	BrandedContentTagInfo     struct {
		CanAddTag bool `json:"can_add_tag"`
	} `json:"branded_content_tag_info"`
	ShoppingInfo        interface{} `json:"shopping_info"`
	AdditionalAudioInfo struct {
		AdditionalAudioUsername interface{} `json:"additional_audio_username"`
		AudioReattributionInfo  struct {
			ShouldAllowRestore bool `json:"should_allow_restore"`
		} `json:"audio_reattribution_info"`
	} `json:"additional_audio_info"`
	IsSharedToFb            bool        `json:"is_shared_to_fb"`
	BreakingContentInfo     interface{} `json:"breaking_content_info"`
	ChallengeInfo           interface{} `json:"challenge_info"`
	ReelsOnTheRiseInfo      interface{} `json:"reels_on_the_rise_info"`
	BreakingCreatorInfo     interface{} `json:"breaking_creator_info"`
	AssetRecommendationInfo interface{} `json:"asset_recommendation_info"`
	ContextualHighlightInfo interface{} `json:"contextual_highlight_info"`
	ClipsCreationEntryPoint string      `json:"clips_creation_entry_point"`
	AudioRankingInfo        struct {
		BestAudioClusterId string `json:"best_audio_cluster_id"`
	} `json:"audio_ranking_info"`
	TemplateInfo        interface{} `json:"template_info"`
	IsFanClubPromoVideo interface{} `json:"is_fan_club_promo_video"`
}

type MediaLocation struct {
	Pk                  int64   `json:"pk"`
	ShortName           string  `json:"short_name"`
	FacebookPlacesId    int64   `json:"facebook_places_id"`
	ExternalSource      string  `json:"external_source"`
	Name                string  `json:"name"`
	Address             string  `json:"address"`
	City                string  `json:"city"`
	HasViewerSaved      bool    `json:"has_viewer_saved"`
	Lng                 float64 `json:"lng"`
	Lat                 float64 `json:"lat"`
	IsEligibleForGuides bool    `json:"is_eligible_for_guides"`
}

type CreativeConfig struct {
	EffectIds        []int64       `json:"effect_ids"`
	CaptureType      string        `json:"capture_type"`
	CreationToolInfo []interface{} `json:"creation_tool_info"`
	EffectConfigs    []struct {
		Name              string      `json:"name"`
		Id                string      `json:"id"`
		FailureReason     *string     `json:"failure_reason"`
		FailureCode       *string     `json:"failure_code"`
		Gatekeeper        interface{} `json:"gatekeeper"`
		AttributionUserId string      `json:"attribution_user_id"`
		AttributionUser   struct {
			InstagramUserId string `json:"instagram_user_id"`
			Username        string `json:"username"`
			ProfilePicture  struct {
				Uri string `json:"uri"`
			} `json:"profile_picture"`
		} `json:"attribution_user"`
		Gatelogic               interface{} `json:"gatelogic"`
		SaveStatus              string      `json:"save_status"`
		EffectActions           []string    `json:"effect_actions"`
		IsSpotRecognitionEffect bool        `json:"is_spot_recognition_effect"`
		IsSpotEffect            bool        `json:"is_spot_effect"`
		ThumbnailImage          struct {
			Uri string `json:"uri"`
		} `json:"thumbnail_image"`
		EffectActionSheet struct {
			PrimaryActions   []string `json:"primary_actions"`
			SecondaryActions []string `json:"secondary_actions"`
		} `json:"effect_action_sheet"`
		DevicePosition           interface{} `json:"device_position"`
		FanClub                  interface{} `json:"fan_club"`
		FormattedClipsMediaCount interface{} `json:"formatted_clips_media_count"`
	} `json:"effect_configs"`
}

type Media struct {
	TakenAt         int    `json:"taken_at"`
	Pk              int64  `json:"pk"`
	Id              string `json:"id"`
	DeviceTimestamp int64  `json:"device_timestamp"`
	MediaType       int    `json:"media_type"`
	Code            string `json:"code"`
	ClientCacheKey  string `json:"client_cache_key"`
	FilterType      int    `json:"filter_type"`
	IsUnifiedVideo  bool   `json:"is_unified_video"`

	User MediaUser `json:"user"`

	CanViewerReshare                    bool             `json:"can_viewer_reshare"`
	CaptionIsEdited                     bool             `json:"caption_is_edited"`
	LikeAndViewCountsDisabled           bool             `json:"like_and_view_counts_disabled"`
	CommercialityStatus                 string           `json:"commerciality_status"`
	IsPaidPartnership                   bool             `json:"is_paid_partnership"`
	IsVisualReplyCommenterNoticeEnabled bool             `json:"is_visual_reply_commenter_notice_enabled"`
	CommentLikesEnabled                 bool             `json:"comment_likes_enabled"`
	CommentThreadingEnabled             bool             `json:"comment_threading_enabled"`
	HasMoreComments                     bool             `json:"has_more_comments"`
	NextMaxId                           int64            `json:"next_max_id,omitempty"`
	MaxNumVisiblePreviewComments        int              `json:"max_num_visible_preview_comments"`
	PreviewComments                     []PreviewComment `json:"preview_comments"`
	CanViewMorePreviewComments          bool             `json:"can_view_more_preview_comments"`
	CommentCount                        int              `json:"comment_count"`
	HideViewAllCommentEntrypoint        bool             `json:"hide_view_all_comment_entrypoint"`
	ImageVersions2                      ImageVersion2    `json:"image_versions2"`
	OriginalWidth                       int              `json:"original_width"`
	OriginalHeight                      int              `json:"original_height"`
	LikeCount                           int              `json:"like_count"`
	HasLiked                            bool             `json:"has_liked"`
	PhotoOfYou                          bool             `json:"photo_of_you"`
	Usertags                            UserTag          `json:"usertags,omitempty"`
	IsOrganicProductTaggingEligible     bool             `json:"is_organic_product_tagging_eligible"`
	CanSeeInsightsAsBrand               bool             `json:"can_see_insights_as_brand"`
	IsDashEligible                      int              `json:"is_dash_eligible"`
	VideoDashManifest                   string           `json:"video_dash_manifest"`
	VideoCodec                          string           `json:"video_codec"`
	NumberOfQualities                   int              `json:"number_of_qualities"`
	VideoVersions                       []VideoVersion   `json:"video_versions"`
	HasAudio                            bool             `json:"has_audio"`
	VideoDuration                       float64          `json:"video_duration"`
	ViewCount                           int              `json:"view_count"`
	PlayCount                           int              `json:"play_count"`
	Caption                             *MediaCaption    `json:"caption"`
	CanViewerSave                       bool             `json:"can_viewer_save"`
	OrganicTrackingToken                string           `json:"organic_tracking_token"`
	HasSharedToFb                       int              `json:"has_shared_to_fb"`
	SharingFrictionInfo                 struct {
		ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
		BloksAppUrl               interface{} `json:"bloks_app_url"`
	} `json:"sharing_friction_info"`
	CommentInformTreatment struct {
		ShouldHaveInformTreatment bool   `json:"should_have_inform_treatment"`
		Text                      string `json:"text"`
	} `json:"comment_inform_treatment"`
	ProductType               string        `json:"product_type"`
	IsInProfileGrid           bool          `json:"is_in_profile_grid"`
	ProfileGridControlEnabled bool          `json:"profile_grid_control_enabled"`
	DeletedReason             int           `json:"deleted_reason"`
	IntegrityReviewDecision   string        `json:"integrity_review_decision"`
	MusicMetadata             interface{}   `json:"music_metadata"`
	ClipsMetadata             ClipsMetadata `json:"clips_metadata"`
	MediaCroppingInfo         struct {
		SquareCrop struct {
			CropLeft   float64 `json:"crop_left"`
			CropRight  float64 `json:"crop_right"`
			CropTop    float64 `json:"crop_top"`
			CropBottom float64 `json:"crop_bottom"`
		} `json:"square_crop"`
	} `json:"media_cropping_info"`
	MezqlToken                  string         `json:"mezql_token"`
	LoggingInfoToken            string         `json:"logging_info_token"`
	CommentingDisabledForViewer bool           `json:"commenting_disabled_for_viewer,omitempty"`
	CreativeConfig              CreativeConfig `json:"creative_config,omitempty"`
	Location                    MediaLocation  `json:"location,omitempty"`
	Lat                         float64        `json:"lat,omitempty"`
	Lng                         float64        `json:"lng,omitempty"`
}

type VideosFeedResp struct {
	BaseApiResp
	Items []struct {
		Media Media `json:"media"`
	} `json:"items"`
	PagingInfo struct {
		MaxId         string `json:"max_id"`
		MoreAvailable bool   `json:"more_available"`
	} `json:"paging_info"`
}
