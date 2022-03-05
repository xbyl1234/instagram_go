package goinsta

import (
	"fmt"
	"time"
)

type Account struct {
	inst    *Instagram
	HadSync bool
	ID      int64

	Detail *UserDetail
}

type curUser struct {
	User struct {
		Pk                         int64  `json:"pk"`
		Username                   string `json:"username"`
		FullName                   string `json:"full_name"`
		IsPrivate                  bool   `json:"is_private"`
		ProfilePicUrl              string `json:"profile_pic_url"`
		ProfilePicId               string `json:"profile_pic_id"`
		IsVerified                 bool   `json:"is_verified"`
		FollowFrictionType         int    `json:"follow_friction_type"`
		HasAnonymousProfilePicture bool   `json:"has_anonymous_profile_picture"`
		Biography                  string `json:"biography"`
		CanLinkEntitiesInBio       bool   `json:"can_link_entities_in_bio"`
		BiographyWithEntities      struct {
			RawText  string        `json:"raw_text"`
			Entities []interface{} `json:"entities"`
		} `json:"biography_with_entities"`
		ExternalUrl            string `json:"external_url"`
		ExternalLynxUrl        string `json:"external_lynx_url"`
		ShowFbLinkOnProfile    bool   `json:"show_fb_link_on_profile"`
		PrimaryProfileLinkType int    `json:"primary_profile_link_type"`
		ReelAutoArchive        string `json:"reel_auto_archive"`
		HdProfilePicVersions   []struct {
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Url    string `json:"url"`
		} `json:"hd_profile_pic_versions"`
		HdProfilePicUrlInfo struct {
			Url    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"hd_profile_pic_url_info"`
		ShowConversionEditEntry                    bool          `json:"show_conversion_edit_entry"`
		AllowedCommenterType                       string        `json:"allowed_commenter_type"`
		HasHighlightReels                          bool          `json:"has_highlight_reels"`
		IsBusiness                                 bool          `json:"is_business"`
		ProfessionalConversionSuggestedAccountType int           `json:"professional_conversion_suggested_account_type"`
		AccountType                                int           `json:"account_type"`
		IsCallToActionEnabled                      interface{}   `json:"is_call_to_action_enabled"`
		InteropMessagingUserFbid                   int64         `json:"interop_messaging_user_fbid"`
		IsHideMoreCommentEnabled                   bool          `json:"is_hide_more_comment_enabled"`
		PersonalAccountAdsPageName                 string        `json:"personal_account_ads_page_name"`
		PersonalAccountAdsPageId                   string        `json:"personal_account_ads_page_id"`
		FbidV2                                     int64         `json:"fbid_v2"`
		IsMutedWordsGlobalEnabled                  bool          `json:"is_muted_words_global_enabled"`
		IsMutedWordsCustomEnabled                  bool          `json:"is_muted_words_custom_enabled"`
		Birthday                                   string        `json:"birthday"`
		BiographyProductMentions                   []interface{} `json:"biography_product_mentions"`
		PhoneNumber                                string        `json:"phone_number"`
		CountryCode                                int           `json:"country_code"`
		NationalNumber                             int64         `json:"national_number"`
		Gender                                     int           `json:"gender"`
		Email                                      string        `json:"email"`
		Pronouns                                   []string      `json:"pronouns"`
		CustomGender                               string        `json:"custom_gender"`
		ProfileEditParams                          struct {
			Username struct {
				ShouldShowConfirmationDialog bool   `json:"should_show_confirmation_dialog"`
				IsPendingReview              bool   `json:"is_pending_review"`
				ConfirmationDialogText       string `json:"confirmation_dialog_text"`
				DisclaimerText               string `json:"disclaimer_text"`
			} `json:"username"`
			FullName struct {
				ShouldShowConfirmationDialog bool   `json:"should_show_confirmation_dialog"`
				IsPendingReview              bool   `json:"is_pending_review"`
				ConfirmationDialogText       string `json:"confirmation_dialog_text"`
				DisclaimerText               string `json:"disclaimer_text"`
			} `json:"full_name"`
		} `json:"profile_edit_params"`
	} `json:"user"`
	Status string `json:"status"`
}

func (this *Account) Sync() error {
	var resp RespUserInfo
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlCurrentUser,
		Query: map[string]interface{}{
			"edit": true,
		},
		IsPost: false}, &resp)

	if err == nil {
		this.Detail = &resp.User
		this.HadSync = true
	}
	err = resp.CheckError(err)
	return err
}

type RespChangeProfilePicture struct {
	BaseApiResp
	User User `json:"user"`
}

type UserProfile struct {
	ExternalUrl string `json:"external_url"`
	Biography   string `json:"biography"`
	FirstName   string `json:"first_name"`
	UploadId    string `json:"upload_id"`
	//Username    string `json:"username"`
	Email string `json:"email"`
}

func (this *Account) EditProfile(profile *UserProfile) error {
	if profile.ExternalUrl == "" {
		profile.ExternalUrl = this.Detail.ExternalURL
	}
	if profile.Biography == "" {
		profile.Biography = this.Detail.Biography
	}
	if profile.FirstName == "" {
		profile.FirstName = this.Detail.FullName
	}
	//if profile.Username == "" {
	//	profile.Username = this.Detail.Username
	//}
	if profile.Email == "" {
		profile.Email = this.Detail.Email
	}

	params := map[string]interface{}{
		//"client_timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		//"timezone_offset":  this.inst.AccountInfo.Location.Timezone,.
		"_uuid":        this.inst.AccountInfo.Device.DeviceID,
		"_uid":         this.inst.ID,
		"phone_number": this.Detail.PhoneNumber,
		"external_url": profile.ExternalUrl,
		"biography":    profile.Biography,
		"first_name":   profile.FirstName,
		"username":     this.inst.User,
		"device_id":    this.inst.AccountInfo.Device.DeviceID,
		"email":        profile.Email,
	}

	if profile.UploadId != "" {
		params["upload_id"] = profile.UploadId
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlEditProfile,
		Query:   params,
		Signed:  true,
		IsPost:  true},
		&resp)
	err = resp.CheckError(err)
	return err
}

func (this *Account) ChangeProfilePicture(uploadID string) error {
	params := map[string]interface{}{
		"waterfall_id": "",
		//"share_to_feed":    "true",
		"_uuid":            this.inst.AccountInfo.Device.DeviceID,
		"_uid":             this.inst.ID,
		"device_id":        this.inst.AccountInfo.Device.DeviceID,
		"client_timestamp": time.Now().Unix(),
		"upload_id":        uploadID,
		"timezone_offset":  this.inst.AccountInfo.Location.Timezone,
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlChangeProfilePicture,
		Query:   params,
		Signed:  true,
		IsPost:  true,
	}, &resp)

	err = resp.CheckError(err)
	return err
}

func (this *Account) UpdatePronouns(pronouns string) error {
	params := map[string]interface{}{
		"_uuid":                      this.inst.AccountInfo.Device.DeviceID,
		"_uid":                       fmt.Sprintf("%d", this.inst.ID),
		"is_pronouns_followers_only": "off",
		"bloks_versioning_id":        this.inst.AccountInfo.Device.BloksVersionID,
		"pronouns":                   "[\n  \"" + pronouns + "\"\n]",
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUpdatePronouns,
		Query:   params,
		Signed:  true,
		IsPost:  true,
	}, &resp)

	err = resp.CheckError(err)
	return err
}

func (this *Account) SetBiography(bio string) error {
	params := map[string]interface{}{
		"_uuid":     this.inst.AccountInfo.Device.DeviceID,
		"_uid":      fmt.Sprintf("%d", this.inst.ID),
		"device_id": this.inst.AccountInfo.Device.DeviceID,
		"raw_text":  bio,
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSetBiography,
		Query:   params,
		Signed:  true,
		IsPost:  true,
	}, &resp)

	err = resp.CheckError(err)
	return err
}

func (this *Account) SetGender() error {
	params := map[string]interface{}{
		"_uuid":         this.inst.AccountInfo.Device.DeviceID,
		"custom_gender": "",
		"gender":        "1",
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSetGender,
		Query:   params,
		Signed:  false,
		IsPost:  true,
	}, &resp)

	err = resp.CheckError(err)
	return err
}

func (this *Account) SendConfirmEmail(email string) error {
	params := map[string]interface{}{
		"_uuid":       this.inst.AccountInfo.Device.DeviceID,
		"_uid":        fmt.Sprintf("%d", this.inst.ID),
		"device_id":   this.inst.AccountInfo.Device.DeviceID,
		"send_source": "personal_information",
		"email":       email,
	}

	var resp RespChangeProfilePicture
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlSendConfirmEmail,
		Query:   params,
		Signed:  true,
		IsPost:  true,
	}, &resp)

	err = resp.CheckError(err)
	return err
}
