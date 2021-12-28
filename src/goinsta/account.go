package goinsta

import (
	"fmt"
	"io/ioutil"
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
			//	from_module=self_profile
		}},
		&resp)

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

func (this *Account) ChangeProfilePicture(path string) error {
	data, err := ioutil.ReadFile(path)
	upID, err := this.inst.GetUpload().RuploadPhoto(data)
	if err != nil {
		return err
	}
	var resp RespChangeProfilePicture
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlChangeProfilePicture,
		Query: map[string]interface{}{
			"_uuid":          this.inst.deviceID,
			"upload_id":      upID,
			"use_fbuploader": "true",
		},
		Signed: true,
		IsPost: true},
		&resp)
	err = resp.CheckError(err)
	return err
}
