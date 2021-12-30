package goinsta

import (
	"fmt"
	"makemoney/common"
)

var (
	InstagramHost       = "https://i.instagram.com/"
	InstagramHost_B     = "https://b.i.instagram.com/"
	InstagramHost_Graph = "https://graph.instagram.com/"
	InstagramUserAgent  = "Instagram %s (iPhone7,2; iOS 12_5_5; en_US; en-US; scale=2.00; 750x1334; %s) AppleWebKit/420+"
	InstagramVersion    = "190.0.0.26.119"
	InstagramLocation   = "en_US"
	//InstagramDeviceList = []string{"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G935F; hero2lte; samsungexynos8890",
	//	"24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom",
	//	"23/6.0.1; 640dpi; 1440x2392; LGE/lge; RS988; h1; h1",
	//	"24/7.0; 640dpi; 1440x2560; HUAWEI; LON-L29; HWLON; hi3660",
	//	"23/6.0.1; 640dpi; 1440x2560; ZTE; ZTE A2017U; ailsa_ii; qcom",
	//	"23/6.0.1; 640dpi; 1440x2560; samsung; SM-G930F; herolte; samsungexynos8890"}
	InstagramBloksVersionID = "5538d18f11cad2fa88efa530f8a717c5d5339d1d53fc5140af9125216d1f7a89"
	InstagramAppID          = "124024574287414"
	InstagramAccessToken    = "124024574287414|84a456d620314b6e92a16d8ff1c792dc"
)

func GenUserAgent() string {
	//return InstagramUserAgent
	//return fmt.Sprintf("Instagram %s Android (%s; %s)", InstagramVersion, InstagramDeviceList[rand.Intn(len(InstagramDeviceList))], InstagramLocation)
	return fmt.Sprintf(InstagramUserAgent, InstagramVersion, common.GenString(common.CharSet_123, 9))
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
	urlGetNamePrefill      = "api/v1/accounts/get_name_prefill/"
	urlLookup              = "api/v1/users/lookup/"
	urlNewAccountNuxSeen   = "api/v1/nux/new_account_nux_seen/"
	// account
	urlCurrentUser          = "api/v1/accounts/current_user/"
	urlChangePass           = "api/v1/accounts/change_password/"
	urlSetPrivate           = "api/v1/accounts/set_private/"
	urlSetPublic            = "api/v1/accounts/set_public/"
	urlRemoveProfPic        = "api/v1/accounts/remove_profile_picture/"
	urlFeedSaved            = "api/v1/feed/saved/"
	urlSetBiography         = "api/v1/accounts/set_biography/"
	urlFeedLiked            = "api/v1/feed/liked/"
	urlChangeProfilePicture = "api/v1/accounts/change_profile_picture/"
	urlEditProfile          = "api/v1/accounts/edit_profile/"
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
