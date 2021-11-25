package goinsta

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"makemoney/common"
	"strconv"
	"time"
)

type Upload struct {
	inst *Instagram
}

func newUpload(inst *Instagram) *Upload {
	return &Upload{inst: inst}
}

type RespRuploadIgPhoto struct {
	BaseApiResp
	Offset         int         `json:"offset"`
	UploadId       string      `json:"upload_id"`
	XsharingNonces interface{} `json:"xsharing_nonces"`
}

func (this *Upload) RuploadIgPhoto(path string) (string, error) {
	upId := strconv.FormatInt(time.Now().Unix(), 10)
	entity_name := upId + "_0_" + common.GenString(common.CharSet_123, 10)

	image_compression, _ := json.Marshal(map[string]string{
		"lib_name":    "moz",
		"lib_version": "3.1.m",
		"quality":     "70",
	})

	params, _ := json.Marshal(map[string]string{
		"upload_id":         upId,
		"media_type":        "1",
		"image_compression": string(image_compression),
	})
	params_str := string(params)
	var resp = &RespRuploadIgPhoto{}
	_, err := this.inst.HttpSend(&sendOptions{
		Url: goInstaHost + urlUploadStory + entity_name,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": params_str,
		},
	}, &resp)

	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	_, err = this.inst.HttpSend(&sendOptions{
		Url:    goInstaHost + urlUploadStory + entity_name,
		IsPost: true,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": params_str,
			"content-type":               "application/octet-stream",
			"x-entity-type":              "image/jpeg",
			"offset":                     "0",
			"x-entity-name":              entity_name,
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, &resp)
	err = resp.CheckError(err)
	return upId, err
}
