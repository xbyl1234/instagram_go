package goinsta

import (
	"fmt"
	"makemoney/common"
)

var (
	goInstaHost       = "https://i.instagram.com/"
	goInstaHost_B     = "https://b.i.instagram.com/"
	goInstaHost_Graph = "https://graph.instagram.com/"
	//goInstaUserAgent  = "Instagram 107.0.0.27.121 Android (24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom; en_US)"
	goInstaUserAgent = "Instagram 140.0.0.30.126 Android (27/8.1.0; 420dpi; 1080x1794; Google/google; Pixel; sailfish; sailfish; en_US; %s)"
	//goInstaUserAgent   = "Instagram 187.0.0.32.120 Android (27/8.1.0; 560dpi; 1440x2712; Google/google; Pixel 2 XL; taimen; taimen; en_US; %s)"
	goInstaExperiments = "ig_android_fci_onboarding_friend_search,ig_android_device_detection_info_upload,ig_android_account_linking_upsell_universe,ig_android_direct_main_tab_universe_v2,ig_android_sign_in_help_only_one_account_family_universe,ig_android_sms_retriever_backtest_universe,ig_android_direct_add_direct_to_android_native_photo_share_sheet,ig_growth_android_profile_pic_prefill_with_fb_pic_2,ig_account_identity_logged_out_signals_global_holdout_universe,ig_android_login_identifier_fuzzy_match,ig_android_mas_remove_close_friends_entrypoint,ig_android_video_render_codec_low_memory_gc,ig_android_email_fuzzy_matching_universe,ig_android_direct_send_like_from_notification,ig_android_suma_landing_page,ig_android_prefetch_debug_dialog,ig_android_smartlock_hints_universe,ig_activation_global_discretionary_sms_holdout,ig_android_video_ffmpegutil_pts_fix,ig_android_multi_tap_login_new,ig_android_caption_typeahead_fix_on_o_universe,ig_android_enable_keyboardlistener_redesign,ig_android_nux_add_email_device,ig_android_direct_remove_view_mode_stickiness_universe,ig_android_new_users_one_tap_holdout_universe,ig_android_mas_notification_badging_universe,ig_android_secondary_account_creation_universe,ig_android_account_recovery_auto_login,ig_android_sim_info_upload,ig_android_mobile_http_flow_device_universe,ig_android_gmail_oauth_in_reg,ig_android_vc_interop_use_test_igid_universe,ig_android_notification_unpack_universe,ig_android_quickcapture_keep_screen_on,ig_android_device_based_country_verification,ig_android_reg_modularization_universe,ig_android_device_verification_separate_endpoint,ig_android_one_login_toast_universe,ig_android_retry_create_account_universe,ig_android_family_apps_user_values_provider_universe,ig_android_reg_nux_headers_cleanup_universe,ig_android_get_cookie_with_concurrent_session_universe,ig_android_device_info_foreground_reporting,ig_android_shortcuts_2019,ig_android_device_verification_fb_signup,ig_android_passwordless_account_password_creation_universe,ig_android_black_out_toggle_universe,ig_video_debug_overlay,ig_android_ask_for_permissions_on_reg,ig_assisted_login_universe,ig_android_security_intent_switchoff,ig_android_passwordless_auth,ig_android_recovery_one_tap_holdout_universe,ig_android_modularized_dynamic_nux_universe,ig_android_fb_account_linking_sampling_freq_universe,ig_android_fix_sms_read_lollipop,ig_android_access_flow_prefill"
	goInstaVersion     = "140.0.0.30.126"
	goInstaBuildNum    = "289692181"
	goInstaLocation    = "en_US"
	DeviceList         = []string{"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G935F; hero2lte; samsungexynos8890",
		"24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom",
		"23/6.0.1; 640dpi; 1440x2392; LGE/lge; RS988; h1; h1",
		"24/7.0; 640dpi; 1440x2560; HUAWEI; LON-L29; HWLON; hi3660",
		"23/6.0.1; 640dpi; 1440x2560; ZTE; ZTE A2017U; ailsa_ii; qcom",
		"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G930F; herolte; samsungexynos8890"}
	BloksVersionID = "a28d5c7230ceed88159f332dce4ad89ff4ceb589502350df7965ce295cdce4bb"
	AppID          = "567067343352427"
)

func GenUserAgent() string {
	//return goInstaUserAgent
	//return fmt.Sprintf("Instagram %s Android (%s; %s)", goInstaVersion, DeviceList[rand.Intn(len(DeviceList))], goInstaLocation)
	return fmt.Sprintf(goInstaUserAgent, common.GenString(common.CharSet_123, 9))
}

type muteOption string

const (
	MuteAll   muteOption = "all"
	MuteStory muteOption = "story"
	MuteFeed  muteOption = "feed"
)

// Endpoints (with format vars)
//注册流程
//	i
///api/v1/users/check_username/
///api/v1/consent/new_user_flow_begins/
///api/v1/dynamic_onboarding/get_steps/
///api/v1/accounts/create_validated/
//	b
///api/v1/launcher/sync/
///api/v1/zr/token/result/
///api/v1/accounts/contact_point_prefill/
///api/v1/multiple_accounts/get_account_family/
///api/v1/dynamic_onboarding/get_steps/
///api/v1/nux/new_account_nux_seen/

const (
	// login
	urlMsisdnHeader      = "api/v1/accounts/read_msisdn_header/"
	urlContactPrefill    = "api/v1/accounts/contact_point_prefill/"
	urlZrToken           = "api/v1/zr/token/result/"
	urlLogin             = "api/v1/accounts/login/"
	urlLogout            = "api/v1/accounts/logout/"
	urlAutoComplete      = "api/v1/friendships/autocomplete_user_list/"
	urlQeSync            = "api/v1/qe/sync/"
	urlLogAttribution    = "api/v1/attribution/log_attribution/"
	urlMegaphoneLog      = "api/v1/megaphone/log/"
	urlExpose            = "api/v1/qe/expose/"
	urlPrefillCandidates = "api/v1/accounts/get_prefill_candidates/"
	//register  v1
	urlCheckPhoneNumber      = "api/v1/accounts/check_phone_number/"
	urlSendSignupSmsCode     = "api/v1/accounts/send_signup_sms_code/"
	urlValidateSignupSmsCode = "api/v1/accounts/validate_signup_sms_code/"
	urlUsernameSuggestions   = "api/v1/accounts/username_suggestions/"
	urlCreateValidated       = "api/v1/accounts/create_validated/"
	//urlCreateValidated       = "api/v1/accounts/create/"
	urlCheckUsername       = "api/v1/users/check_username/"
	urlLauncherSync        = "api/v1/launcher/sync/"
	urlCheckAgeEligibility = "api/v1/consent/check_age_eligibility/"
	urlNewUserFlowBegins   = "api/v1/consent/new_user_flow_begins/"
	urlGetSteps            = "api/v1/dynamic_onboarding/get_steps/"

	// account
	urlCurrentUser          = "api/v1/accounts/current_user/"
	urlChangePass           = "api/v1/accounts/change_password/"
	urlSetPrivate           = "api/v1/accounts/set_private/"
	urlSetPublic            = "api/v1/accounts/set_public/"
	urlRemoveProfPic        = "api/v1/accounts/remove_profile_picture/"
	urlFeedSaved            = "api/v1/feed/saved/"
	urlSetBiography         = "api/v1/accounts/set_biography/"
	urlEditProfile          = "api/v1/accounts/edit_profile"
	urlFeedLiked            = "api/v1/feed/liked/"
	urlChangeProfilePicture = "api/v1/accounts/change_profile_picture/"

	// account and profile
	urlFollowers = "api/v1/friendships/%d/followers/"
	urlFollowing = "api/v1/friendships/%d/following/"

	// users

	urlUserArchived      = "api/v1/feed/only_me_feed/"
	urlUserByName        = "api/v1/users/%s/usernameinfo/"
	urlUserByID          = "api/v1/users/%d/info/"
	urlUserBlock         = "api/v1/friendships/block/%d/"
	urlUserUnblock       = "api/v1/friendships/unblock/%d/"
	urlUserMute          = "api/v1/friendships/mute_posts_or_story_from_follow/"
	urlUserUnmute        = "api/v1/friendships/unmute_posts_or_story_from_follow/"
	urlUserFollow        = "api/v1/friendships/create/%d/"
	urlUserUnfollow      = "api/v1/friendships/destroy/%d/"
	urlUserFeed          = "api/v1/feed/user/%d/"
	urlFriendship        = "api/v1/friendships/show/%d/"
	urlFriendshipPending = "api/v1/friendships/pending/"
	urlUserStories       = "api/v1/feed/user/%d/reel_media/"
	urlUserTags          = "api/v1/usertags/%d/feed/"
	urlBlockedList       = "api/v1/users/blocked_list/"
	urlUserInfo          = "api/v1/users/%d/info/"
	urlUserHighlights    = "api/v1/highlights/%d/highlights_tray/"
	urlFriendFollowers   = "api/v1/friendships/%v/followers/"

	// timeline
	urlTimeline  = "api/v1/feed/timeline/"
	urlStories   = "api/v1/feed/reels_tray/"
	urlReelMedia = "api/v1/feed/reels_media/"

	// search
	urlSearchUser     = "api/v1/users/search/"
	urlSearchTag      = "api/v1/tags/search/"
	urlSearchLocation = "api/v1/fbsearch/places"
	urlSearchFacebook = "api/v1/fbsearch/topsearch/"

	// feeds
	urlFeedLocationID = "api/v1/feed/location/%d/"
	urlFeedLocations  = "api/v1/locations/%d/sections/"
	urlFeedTag        = "api/v1/feed/tag/%s/"

	// media
	urlMediaInfo   = "api/v1/media/%s/info/"
	urlMediaDelete = "api/v1/media/%s/delete/"
	urlMediaLike   = "api/v1/media/%s/like/"
	urlMediaUnlike = "api/v1/media/%s/unlike/"
	urlMediaSave   = "api/v1/media/%s/save/"
	urlMediaUnsave = "api/v1/media/%s/unsave/"
	urlMediaSeen   = "api/v1/media/seen/"
	urlMediaLikers = "api/v1/media/%s/likers/"

	// comments
	urlCommentAdd     = "api/v1/media/%d/comment/"
	urlCommentDelete  = "api/v1/media/%s/comment/%s/delete/"
	urlComment        = "api/v1/media/%s/comments/"
	urlCommentDisable = "api/v1/media/%s/disable_comments/"
	urlCommentEnable  = "api/v1/media/%s/enable_comments/"
	urlCommentLike    = "api/v1/media/%s/comment_like/"
	urlCommentUnlike  = "api/v1/media/%s/comment_unlike/"

	// activity
	urlActivityFollowing = "api/v1/news/"
	urlActivityRecent    = "api/v1/news/inbox/"

	// inbox
	urlInbox         = "api/v1/direct_v2/inbox/"
	urlInboxPending  = "api/v1/direct_v2/pending_inbox/"
	urlInboxSendLike = "api/v1/direct_v2/threads/broadcast/like/"
	urlReplyStory    = "api/v1/direct_v2/threads/broadcast/reel_share/"
	urlInboxThread   = "api/v1/direct_v2/threads/%s/"
	urlInboxMute     = "api/v1/direct_v2/threads/%s/mute/"
	urlInboxUnmute   = "api/v1/direct_v2/threads/%s/unmute/"

	// tags
	urlTagSync     = "api/v1/tags/%s/info/"
	urlTagStories  = "api/v1/tags/%s/story/"
	urlTagContent  = "api/v1/tags/%s/ranked_sections/"
	urlTagSections = "api/v1/tags/%s/sections/"

	// upload
	urlUploadPhone = "rupload_igphoto/"
	urlUploadVideo = "rupload_igvideo/"

	// msg
	urlSendText  = "api/v1/direct_v2/threads/broadcast/text/"
	urlSendImage = "api/v1/direct_v2/threads/broadcast/configure_photo/"

	//	graph
	urlLoggingClientEvents = "logging_client_events"
)
