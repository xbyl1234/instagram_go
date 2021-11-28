package goinsta

import (
	"fmt"
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
		ApiPath: fmt.Sprintf(urlUserInfo, this.ID),
		Query: map[string]interface{}{
			"edit": true,
		}},
		resp)

	if err == nil {
		this.Detail = &resp.User
		this.HadSync = true
	}
	return err
}

type RespChangeProfilePicture struct {
	BaseApiResp
	User User `json:"user"`
}

func (this *Account) ChangeProfilePicture(path string) error {
	upID, err := this.inst.GetUpload().RuploadIgPhoto(path)
	if err != nil {
		return err
	}
	var resp RespChangeProfilePicture
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlChangeProfilePicture,
		Query: map[string]interface{}{
			"_uuid":          this.inst.uuid,
			"upload_id":      upID,
			"use_fbuploader": "true",
		},
		Signed: true,
		IsPost: true},
		&resp)
	err = resp.CheckError(err)
	return err
}
