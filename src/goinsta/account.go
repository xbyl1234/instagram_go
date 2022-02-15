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
		"client_timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		"phone_number":     this.Detail.PhoneNumber,
		"timezone_offset":  "-28800",
		"external_url":     profile.ExternalUrl,
		"waterfall_id":     this.inst.Device.WaterID,
		"biography":        profile.Biography,
		"first_name":       profile.FirstName,
		"username":         this.inst.User,
		"device_id":        this.inst.Device.DeviceID,
		"email":            profile.Email,
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
		"_uuid":            this.inst.Device.DeviceID,
		"_uid":             this.inst.ID,
		"device_id":        this.inst.Device.DeviceID,
		"client_timestamp": time.Now().Unix(),
		"upload_id":        uploadID,
		"timezone_offset":  this.inst.Device.TimezoneOffset,
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
