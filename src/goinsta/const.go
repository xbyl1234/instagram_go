package goinsta

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"strconv"
	"strings"
)

var (
	InstagramHost        = "https://i.instagram.com/"
	InstagramHost_B      = "https://b.i.instagram.com/"
	InstagramHost_Graph  = "https://graph.instagram.com/"
	InstagramUserAgent2  = "Instagram %s (iPhone7,2; iOS 12_5_5; en_US; en-US; scale=2.00; 750x1334; %s) AppleWebKit/420+"
	InstagramUserAgent   = "Instagram %s (%s; iOS %s; en_US; en-US; %s; %s; %s) AppleWebKit/%s"
	InstagramLocation    = "en_US"
	InstagramVersions    []*InstDeviceInfo
	InstagramAppID       = "124024574287414"
	InstagramAccessToken = "124024574287414|84a456d620314b6e92a16d8ff1c792dc"

	InstagramDeviceList2 = []string{"iPhone7,2 12_5_5 scale=2.00 750x1334"}

	InstagramDeviceList = []string{
		"iPhone7,1 12_5_5 scale=2.61 1080x1920",
		//"iPhone7,2 12_5_5 scale=2.00 750x1334",
		//"iPhone7,2 11_2_2 scale=2.00 750x1334",
		"iPhone8,1 13_7 scale=2.00 750x1334",
		"iPhone8,1 14_6 scale=2.00 750x1334",
		"iPhone8,1 12_1_4 scale=2.34 750x1331",
		"iPhone8,1 14_7_1 scale=2.00 750x1334",
		"iPhone9,1 15_1 scale=2.00 750x1334",
		"iPhone9,3 15_1 scale=2.00 750x1334",
		"iPhone9,3 14_8_1 scale=2.00 750x1334",
		"iPhone9,3 14_1 scale=2.00 750x1334",
		"iPhone9,4 14_6 scale=2.61 1080x1920",
		"iPhone9,4 15_1 scale=2.61 1080x1920",
		"iPhone9,4 14_8_1 scale=2.61 1080x1920",
		"iPhone10,1 15_2 scale=2.00 750x1334",
		"iPhone10,1 14_8_1 scale=2.00 750x1334",
		//"iPhone10,3 15_1 scale=3.00 1125x2436",
		//"iPhone10,4 15_0 scale=2.00 750x1334",
		//"iPhone10,4 14_8_1 scale=2.00 750x1334",
		//"iPhone10,4 15_1 scale=2.00 750x1334",
		//"iPhone10,4 14_7_1 scale=2.00 750x1334",
		//"iPhone10,4 15_0_2 scale=2.00 750x1334",
		//"iPhone10,5 15_1 scale=2.88 1080x1920",
		//"iPhone10,5 15_0_1 scale=2.61 1080x1920",
		//"iPhone10,5 13_5_1 scale=2.61 1080x1920",
		//"iPhone10,6 15_2 scale=3.00 1125x2436",
		//"iPhone10,6 15_1 scale=3.00 1125x2436",
		"iPhone11,2 14_7_1 scale=3.00 1125x2436",
		"iPhone11,2 15_1 scale=3.00 1125x2436",
		"iPhone11,2 14_6 scale=3.00 1125x2436",
		//"iPhone11,6 15_1 scale=3.00 1242x2688",
		//"iPhone11,6 14_8_1 scale=3.00 1242x2688",
		//"iPhone11,8 14_8_1 scale=2.00 828x1792",
		//"iPhone11,8 14_6 scale=2.00 828x1792",
		//"iPhone11,8 15_0_2 scale=2.00 828x1792",
		//"iPhone11,8 15_2 scale=2.00 828x1792",
		//"iPhone11,8 15_1 scale=2.00 828x1792",
		"iPhone12,1 14_7_1 scale=2.00 828x1792",
		"iPhone12,1 14_8_1 scale=2.00 828x1792",
		"iPhone12,1 15_2 scale=2.00 828x1792",
		"iPhone12,1 15_0 scale=2.00 828x1792",
		"iPhone12,1 13_7 scale=2.00 828x1792",
		"iPhone12,1 15_0_2 scale=2.00 828x1792",
		"iPhone12,3 15_0_2 scale=3.00 1125x2436",
		"iPhone12,3 14_6 scale=3.00 1125x2436",
		"iPhone12,3 14_2 scale=3.00 1125x2436",
		"iPhone12,3 15_1 scale=3.00 1125x2436",
		//"iPhone12,5 14_3 scale=3.00 1242x2688",
		//"iPhone12,5 15_1 scale=3.31 1242x2689",
		//"iPhone12,5 14_6 scale=3.00 1242x2688",
		//"iPhone12,5 13_3_1 scale=3.31 1242x2689",
		//"iPhone12,5 14_7_1 scale=3.00 1242x2688",
		//"iPhone12,5 14_8_1 scale=3.00 1242x2688",
		"iPhone12,8 15_2 scale=2.00 750x1334",
		"iPhone12,8 15_1 scale=2.00 750x1334",
		"iPhone12,8 14_8_1 scale=2.00 750x1334",
		"iPhone12,8 14_7_1 scale=2.00 750x1334",
		"iPhone13,1 14_7_1 scale=2.88 1080x2338",
		"iPhone13,2 15_2 scale=3.00 1170x2532",
		"iPhone13,2 15_0_2 scale=3.00 1170x2532",
		"iPhone13,2 15_1_1 scale=3.00 1170x2532",
		"iPhone13,2 14_6 scale=3.00 1170x2532",
		"iPhone13,2 14_7_1 scale=3.00 1170x2532",
		//"iPhone13,3 15_1_1 scale=3.00 1170x2532",
		//"iPhone13,3 15_2 scale=3.00 1170x2532",
		//"iPhone13,3 14_8 scale=3.00 1170x2532",
		//"iPhone13,3 14_8_1 scale=3.00 1170x2532",
		"iPhone13,4 15_1_1 scale=3.00 1284x2778",
		"iPhone13,4 15_2 scale=3.00 1284x2778",
		//"iPhone14,2 15_1_1 scale=3.00 1170x2532",
		//"iPhone14,3 15_2 scale=3.00 1284x2778",
		//"iPhone14,5 15_0 scale=3.00 1170x2532",
		//"iPhone14,5 15_2 scale=3.00 1170x2532",
		//"iPhone14,5 15_0 scale=3.66 1170x2533",
		//"iPhone14,5 15_1_1 scale=3.00 1170x2532",
	}

	LensModel = map[string]string{
		"7,1":  "iPhone,6,back camera,4.15mm,f/2.2",
		"8,1":  "iPhone,6s,back camera,4.15mm,f/2.2",
		"9,1":  "iPhone,7,back camera,3.99mm,f/1.8",
		"9,2":  "iPhone,7,Plus back dual camera,3.99mm,f/1.8",
		"9,3":  "iPhone,7,back camera,3.99mm,f/1.8",
		"9,4":  "iPhone,7,Plus back dual camera,3.99mm,f/1.8",
		"10,1": "iPhone,8,back camera,3.99mm,f/1.8",
		"11,2": "iPhone,XS,Max back dual camera,4.25mm,f/1.8",
		"12,1": "iPhone,11,back dual wide camera,4.25mm,f/1.8",
		"12,3": "iPhone,11,Pro back triple camera,4.25mm,f/1.8",
		"12,8": "iPhone,SE,(2nd generation) back camera,3.99mm,f/1.8",
		"13,1": "iPhone,12,mini back dual wide camera,4.2mm,f/1.6",
		"13,2": "iPhone,12,back dual wide camera,4.2mm,f/1.6",
		"13,4": "iPhone,12,Pro Max back triple camera,5.1mm,f/1.6",
	}

	InstagramVersionData = []string{
		"190.0.0.26.119 294609445 5538d18f11cad2fa88efa530f8a717c5d5339d1d53fc5140af9125216d1f7a89",
		"191.0.0.25.122 296543649 bf3e79f2304601044c85a6f9c44dab59a72558ca9f9a821b96882a4a54ca3c3a",
		"192.0.0.37.119 298025452 9fd8ac08308424f3385019b6c63fc3eb52f3d9d1314f33b78b5db21716d3bf7a",
		"193.0.0.29.121 299401192 02872c8277b5f0ccc5275f61be6e600aeda09ad6927f75ed89e4177ba297bd4f",
		"194.0.0.33.171 301014278 f487d54d4da25bc844e5abc7af030f30a3f83d27bf24d8765274c62268d17dab",
		"195.0.0.24.119 302211069 64f0fd5651d93707c9c3f98da6e5af03e46140770801f4d6eb3c2799911ec21f",
		"196.0.0.21.120 303649428 d5507f7c0ad817ba90a28666727e7cfbe294b33203ab934be3e68e8823c27605",
		"197.0.0.20.119 305020938 29d3248efc1cfc10a0dbafd84eb58fd7eebc6b99e626c5bfa0fe615b8ff784d9",
		"198.0.0.27.119 306495430 4e3e06f8c5ab8a9ab19a536615e0cd79a8eb4742068f2f63af62a4511a69187c",
		"199.0.0.27.120 307960803 4591eb23633140d3162279972d8981f471b81a93fe8dcc4f28d7c24014470eee",
		"200.0.0.19.118 309467717 c8855ee853644de5034628e4716bd7c915e25b17ba524f999d904099b41e910d",
		"200.1.0.20.118 310109037 c8855ee853644de5034628e4716bd7c915e25b17ba524f999d904099b41e910d",
		"202.0.0.23.119 312612729 212380a0da87f76cbe923f35a5e31cd753ec95b4d6d584729e8054420f7822ca",
		"203.0.0.26.117 314122909 0d09d64dfdf69cb3c47fd6f30f0948fbdb48a267aff5aaf8e450a5122dd0a68f",
		"204.0.0.16.119 315460786 f646f3986f57a81c3eef6b8fde357c324737ea272a2f59045ad73fd2275519e5",
		"205.0.0.20.115 317250287 d28e6e26fb0da1c17e9f33a6165eb94acc93bea287ea880be847d1a6796acbb9",
		"206.0.0.30.118 318760365 965488f430b5bb53716f1a026db174f1c0a069c657b3ff9cc11afa555acd2797",
		"206.1.0.30.118 320107871 965488f430b5bb53716f1a026db174f1c0a069c657b3ff9cc11afa555acd2797",
		"207.0.0.28.118 320361397 3ced1a8d80642ea42beaafe4d2769cd5afd3b763d091b0eddb60f30850562184",
		"208.0.0.26.131 322184766 c9e961cf88c88c69f97f01c5d4342d299bdb880c11894424c622d5b97e8a6a76",
		"209.0.0.23.113 324180638 a94a49452ea2139f4109a16413cf473fe4bbe021b27e5f6e86625923ff3969d9",
		"209.1.0.25.113 324511477 a94a49452ea2139f4109a16413cf473fe4bbe021b27e5f6e86625923ff3969d9",
		"210.0.0.16.67 325544617 857f81be49bc66e64a152220da016951113eec808e531a9e235d55e342cf43c4",
		"211.0.0.21.118 327311214 93202078f561251d153a675632166c73377a4c41e127ace6f25ee9484a77c7ce",
		"212.0.0.22.118 328988229 0d38efe9f67cf51962782e8aae19001881099884d8d86c683d374fc1b89ffad1",
		"212.1.0.25.118 329643252 0d38efe9f67cf51962782e8aae19001881099884d8d86c683d374fc1b89ffad1",
		"213.0.0.19.117 330663239 5bf3152c14a8e8651b2ec5a689994b294f4e0a74b86b5652da331aa7035d1c62",
		"213.1.0.22.117 332048479 5bf3152c14a8e8651b2ec5a689994b294f4e0a74b86b5652da331aa7035d1c62",
	}

	InstagramReqSpeed = []string{
		"35kbps",
		"29kbps",
		"42kbps",
		"45kbps",
		"58kbps",
		"78kbps",
		"9kbps",
		"46kbps",
	}
	NoLoginHeaderMap map[string]*HeaderSequence
	LoginHeaderMap   map[string]*HeaderSequence
	HeaderMD5Map     map[string]*HeaderSequence
	ReqHeaderJson    reqHeaderJson
)

type InstDeviceInfo struct {
	Version        string  `json:"version" bson:"version"`
	VersionCode    string  `json:"version_code" bson:"version_code"`
	BloksVersionID string  `json:"bloks_version_id" bson:"bloks_version_id"`
	UserAgent      string  `json:"user_agent" bson:"user_agent"`
	IDFA           string  `json:"idfa" bson:"idfa"`
	AppLocale      string  `json:"app_locale" bson:"app_locale"`
	TimezoneOffset string  `json:"timezone_offset" bson:"timezone_offset"`
	StartupCountry string  `json:"startup_country" bson:"startup_country"`
	AcceptLanguage string  `json:"accept_language" bson:"accept_language"`
	NetWorkType    string  `json:"net_work_type" bson:"net_work_type"`
	DeviceID       string  `json:"device_id" bson:"device_id"`
	FamilyID       string  `json:"family_id" bson:"family_id"`
	WaterID        string  `json:"water_id" bson:"water_id"`
	DeviceToken    string  `json:"device_token" bson:"device_token"`
	SystemVersion  string  `json:"system_version" bson:"system_version"`
	LensModel      string  `json:"lens_model" bson:"lens_model"`
	FocalLength    float64 `json:"focal_length" bson:"focal_length"`
	Aperture       float64 `json:"aperture" bson:"aperture"`
}

type AutoSetHeaderFun func(inst *Instagram, opt *reqOptions, req *http.Request)

type HeaderSequence struct {
	HeaderFun []AutoSetHeaderFun
	HeaderSeq []string
}

type pathsMap struct {
	Path string `json:"path"`
	Md5  string `json:"md5"`
}

type reqHeaderJson struct {
	PathsNoLogin []pathsMap `json:"paths_no_login"`
	PathsLogin   []pathsMap `json:"paths_login"`
	Md5S         []*struct {
		Desp      string `json:"desp,omitempty"`
		Md5       string `json:"md5"`
		Header    string `json:"header"`
		headerSeq HeaderSequence
	} `json:"md5s"`
}

func InitInstagramConst() error {
	InstagramVersions = make([]*InstDeviceInfo, len(InstagramVersionData))
	for index, item := range InstagramVersionData {
		sp := strings.Split(item, " ")
		InstagramVersions[index] = &InstDeviceInfo{
			Version:        sp[0],
			VersionCode:    sp[1],
			BloksVersionID: sp[2],
		}
	}

	err := common.LoadJsonFile("config/http_header_sequence.json", &ReqHeaderJson)
	if err != nil {
		return err
	}
	HeaderMD5Map = make(map[string]*HeaderSequence)
	for _, md5 := range ReqHeaderJson.Md5S {
		sp := strings.Split(md5.Header, ",")
		if len(sp) == 0 {
			return &common.MakeMoneyError{
				ErrStr: "header is null,md5: " + md5.Md5,
			}
		}
		sp = sp[:len(sp)-1]
		md5.headerSeq.HeaderFun = GetAutoHeaderFunc(sp)
		md5.headerSeq.HeaderSeq = sp
		HeaderMD5Map[md5.Md5] = &md5.headerSeq
	}

	var makeSeqMap = func(headerMap *map[string]*HeaderSequence, pahts []pathsMap) {
		*headerMap = make(map[string]*HeaderSequence)
		for _, path := range pahts {
			find := false
			for _, md5 := range ReqHeaderJson.Md5S {
				if path.Md5 == md5.Md5 {
					(*headerMap)[path.Path] = &md5.headerSeq
					find = true
					break
				}
			}
			if !find {
				log.Warn("not find header md5: %s", path.Md5)
			}
		}
	}

	makeSeqMap(&NoLoginHeaderMap, ReqHeaderJson.PathsNoLogin)
	makeSeqMap(&LoginHeaderMap, ReqHeaderJson.PathsLogin)
	return err
}

//纽约 -18000

func GenInstDeviceInfo() *InstDeviceInfo {
	version := InstagramVersions[common.GenNumber(0, len(InstagramVersions))]
	device := InstagramDeviceList[common.GenNumber(0, len(InstagramDeviceList))]
	sp := strings.Split(device, " ")
	SystemVersion := strings.ReplaceAll(sp[1], "_", ".")

	spLens := strings.Split(LensModel[sp[0][6:]], ",")
	Lens := ""
	for _, item := range spLens {
		Lens += item + " "
	}
	Lens = Lens[:len(Lens)-1]
	FocalLength, _ := strconv.ParseFloat(strings.ReplaceAll(spLens[4], "f/", ""), 64)
	Aperture, _ := strconv.ParseFloat(strings.ReplaceAll(spLens[4], "mm", ""), 64)

	instVersion := &InstDeviceInfo{
		IDFA:           strings.ToUpper(common.GenUUID()),
		Version:        version.Version,
		VersionCode:    version.VersionCode,
		BloksVersionID: version.BloksVersionID,
		UserAgent:      fmt.Sprintf(InstagramUserAgent, version.Version, sp[0], sp[1], sp[2], sp[3], version.VersionCode, "420+"),
		AppLocale:      "en-US",
		TimezoneOffset: "-18000",
		StartupCountry: "US",
		AcceptLanguage: "en-US;q=1.0",
		NetWorkType:    "WiFi",
		DeviceID:       strings.ToUpper(common.GenUUID()),
		WaterID:        common.GenString(common.CharSet_16_Num, 32),
		SystemVersion:  SystemVersion,
		LensModel:      Lens,
		FocalLength:    FocalLength,
		Aperture:       Aperture,
	}

	instVersion.FamilyID = instVersion.DeviceID
	return instVersion
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
)
