package goinsta

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
	urlMsisdnHeader          = "api/v1/accounts/read_msisdn_header/"
	urlContactPrefill        = "api/v1/accounts/contact_point_prefill/"
	urlZrToken               = "api/v1/zr/token/result/"
	urlLogin                 = "api/v1/accounts/login/"
	urlLogout                = "api/v1/accounts/logout/"
	urlAutoComplete          = "api/v1/friendships/autocomplete_user_list/"
	urlQeSync                = "api/v1/qe/sync/"
	urlLogAttribution        = "api/v1/attribution/log_attribution/"
	urlDeviceRegister        = "api/v1/push/register/"
	urlMegaphoneLog          = "api/v1/megaphone/log/"
	urlExpose                = "api/v1/qe/expose/"
	urlPrefillCandidates     = "api/v1/accounts/get_prefill_candidates/"
	urlBadge                 = "api/v1/notifications/badge/"
	urlMobileConfig          = "api/v1/launcher/mobileconfig/"
	urlWwwgraphql            = "api/v1/wwwgraphql/ig/query/"
	urlGetAccountFamily      = "api/v1/multiple_accounts/get_account_family/"
	urlGetCoolDowns          = "api/v1/qp/get_cooldowns/"
	urlCheckPhoneNumber      = "api/v1/accounts/check_phone_number/"
	urlSendSignupSmsCode     = "api/v1/accounts/send_signup_sms_code/"
	urlValidateSignupSmsCode = "api/v1/accounts/validate_signup_sms_code/"
	urlUsernameSuggestions   = "api/v1/accounts/username_suggestions/"
	urlCreateValidated       = "api/v1/accounts/create_validated/"
	urlCreate                = "api/v1/accounts/create/"
	urlCheckUsername         = "api/v1/users/check_username/"
	urlLauncherSync          = "api/v1/launcher/sync/"
	urlCheckAgeEligibility   = "api/v1/consent/check_age_eligibility/"
	urlNewUserFlowBegins     = "api/v1/consent/new_user_flow_begins/"
	urlGetSteps              = "api/v1/dynamic_onboarding/get_steps/"
	urlGetNamePrefill        = "api/v1/accounts/get_name_prefill/"
	urlLookup                = "api/v1/users/lookup/"
	urlNewAccountNuxSeen     = "api/v1/nux/new_account_nux_seen/"

	// account
	urlCurrentUser          = "api/v1/accounts/current_user/"
	urlChangePass           = "api/v1/accounts/change_password/"
	urlSetPrivate           = "api/v1/accounts/set_private/"
	urlSetPublic            = "api/v1/accounts/set_public/"
	urlRemoveProfPic        = "api/v1/accounts/remove_profile_picture/"
	urlFeedSaved            = "api/v1/feed/saved/"
	urlFeedLiked            = "api/v1/feed/liked/"
	urlChangeProfilePicture = "api/v1/accounts/change_profile_picture/"
	urlEditProfile          = "api/v1/accounts/edit_profile/"
	// account and profile
	urlFollowers                         = "api/v1/friendships/%d/followers/"
	urlFollowing                         = "api/v1/friendships/%d/following/"
	urlGetSignupConfig                   = "api/v1/consent/get_signup_config/"
	urlGetCommonEmailDomains             = "api/v1/accounts/get_common_email_domains/"
	urlPrecheckCloudId                   = "api/v1/accounts/precheck_cloud_id/"
	urlIgUser                            = "api/v1/fb/ig_user/"
	urlCheckEmail                        = "api/v1/users/check_email/"
	urlSendVerifyEmail                   = "api/v1/accounts/send_verify_email/"
	urlCheckConfirmationCode             = "api/v1/accounts/check_confirmation_code/"
	urlAddressBookLink                   = "api/v1/address_book/link/"
	urlSetGender                         = "api/v1/accounts/set_gender/"
	urlSendConfirmEmail                  = "api/v1/accounts/send_confirm_email/"
	urlSetBiography                      = "api/v1/accounts/set_biography/"
	urlWriteDeviceCapabilities           = "api/v1/creatives/write_device_capabilities/"
	urlUpdatePronouns                    = "api/v1/bloks/apps/com.instagram.equity.pronouns.update_pronouns.action/"
	urlGraphql                           = "api/v1/ads/graphql/"
	urlLocationSearch                    = "api/v1/location_search/"
	urlInvalidatePrivacyViolatingMediaV2 = "api/v1/feed/invalidate_privacy_violating_media_v2/"
	urlGetSsoAccounts                    = "api/v1/fxcal/get_sso_accounts/"

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
	urlMediaInfo                  = "api/v1/media/%s/info/"
	urlMediaDelete                = "api/v1/media/%s/delete/"
	urlMediaLike                  = "api/v1/media/%s/like/"
	urlMediaUnlike                = "api/v1/media/%s/unlike/"
	urlMediaSave                  = "api/v1/media/%s/save/"
	urlMediaUnsave                = "api/v1/media/%s/unsave/"
	urlMediaSeen                  = "api/v1/media/seen/"
	urlMediaLikers                = "api/v1/media/%s/likers/"
	urlConfigureToStory           = "api/v1/media/configure_to_story/"
	urlCreateReel                 = "api/v1/highlights/create_reel/"
	urlSetReelSettings            = "api/v1/users/set_reel_settings/"
	urlConfigureToClips           = "api/v1/media/configure_to_clips/"
	urlClipsInfoForCreation       = "api/v1/clips/clips_info_for_creation/"
	urlClipsAssets                = "api/v1/creatives/clips_assets/"
	urlVerifyOriginalAudioTitle   = "api/v1/music/verify_original_audio_title/"
	urlUpdateVideoWithQualityInfo = "api/v1/media/update_video_with_quality_info/"
	urlConfigure                  = "api/v1/media/configure/"

	// comments
	urlCommentAdd     = "api/v1/media/%s/comment/"
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
	urlUploadPhone  = "rupload_igphoto/"
	urlUploadVideo  = "rupload_igvideo/"
	urlUploadFinish = "api/v1/media/upload_finish/"

	// msg
	urlSendText          = "api/v1/direct_v2/threads/broadcast/text/"
	urlShareVoice        = "api/v1/direct_v2/threads/broadcast/share_voice/"
	urlSendImage         = "api/v1/direct_v2/threads/broadcast/configure_photo/"
	urlCreateGroupThread = "api/v1/direct_v2/create_group_thread/"
	urlSendLink          = "api/v1/direct_v2/threads/broadcast/link/"
	//	graph
	urlLoggingClientEvents = "logging_client_events"

	//feed
	urlDiscoverVideosFeed = "api/v1/discover/videos_feed/"

	urlCheckOffensiveComment = "api/v1/media/comment/check_offensive_comment/"
	urlShareMedia            = "api/v1/media/%s/permalink/"
)
