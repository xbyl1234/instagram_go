package goinsta

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"
)

type Upload struct {
	inst *Instagram
}

func NewUpload(inst *Instagram) *Upload {
	return &Upload{inst: inst}
}

type RespRuploadIgPhoto struct {
	Offset         int         `json:"offset"`
	UploadId       string      `json:"upload_id"`
	XsharingNonces interface{} `json:"xsharing_nonces"`
}

func (this *Upload) RuploadIgPhoto(path string) (string, error) {
	upId := strconv.FormatInt(time.Now().Unix(), 10)

	var resp = &RespRuploadIgPhoto{}
	_, err := this.inst.HttpSend(&sendOptions{
		Url: goInstaHost + urlUploadStory,
		Header: map[string]string{
			"upload_id":               upId,
			"media_type":              "1",
			"x_fb_photo_waterfall_id": "",
		},
	}, &resp)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	_, err = this.inst.HttpSend(&sendOptions{
		Url:    goInstaHost + urlUploadStory,
		IsPost: true,
		Header: map[string]string{
			"upload_id":               upId,
			"media_type":              "1",
			"x_fb_photo_waterfall_id": "",
		},
		Body: body,
	}, &resp)

	return upId, err
}
