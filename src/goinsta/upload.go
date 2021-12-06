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

func (this *Upload) RuploadPhoto(path string) (string, error) {
	upId := strconv.FormatInt(time.Now().Unix(), 10)
	entityName := upId + "_0_" + common.GenString(common.CharSet_123, 10)

	imageCompression, _ := json.Marshal(map[string]string{
		"lib_name":    "moz",
		"lib_version": "3.1.m",
		"quality":     "70",
	})

	params, _ := json.Marshal(map[string]string{
		"upload_id":         upId,
		"media_type":        "1",
		"image_compression": string(imageCompression),
	})
	paramsStr := string(params)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	var resp = &RespRuploadIgPhoto{}
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUploadPhone + entityName,
		IsPost:  true,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": paramsStr,
			"content-type":               "application/octet-stream",
			"x-entity-type":              "image/jpeg",
			"offset":                     "0",
			"x-entity-name":              entityName,
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, &resp)
	err = resp.CheckError(err)
	return upId, err
}

func (this *Upload) RuploadIgPhoto(path string) (string, error) {
	upId := strconv.FormatInt(time.Now().Unix(), 10)
	entityName := upId + "_0_" + common.GenString(common.CharSet_123, 10)

	imageCompression, _ := json.Marshal(map[string]string{
		"lib_name":    "moz",
		"lib_version": "3.1.m",
		"quality":     "70",
	})

	params, _ := json.Marshal(map[string]string{
		"upload_id":         upId,
		"media_type":        "1",
		"image_compression": string(imageCompression),
	})
	paramsStr := string(params)
	var resp = &RespRuploadIgPhoto{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUploadPhone + entityName,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": paramsStr,
		},
	}, resp)

	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUploadPhone + entityName,
		IsPost:  true,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": paramsStr,
			"content-type":               "application/octet-stream",
			"x-entity-type":              "image/jpeg",
			"offset":                     "0",
			"x-entity-name":              entityName,
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	return upId, err
}

func (this *Upload) RuploadVideo(path string) (string, error) {
	upId := common.GenString(common.CharSet_123, 15)
	timeTick := strconv.FormatInt(time.Now().Unix(), 15) + "000"
	entityName := common.GenString(common.CharSet_16_Num, 32) + "-0-" +
		common.GenString(common.CharSet_123, 7) + "-" +
		timeTick + "-" + timeTick

	retryContext, _ := json.Marshal(map[string]interface{}{
		"num_reupload":          0,
		"num_step_auto_retry":   0,
		"num_step_manual_retry": 0,
	})

	params, _ := json.Marshal(map[string]interface{}{
		"upload_media_height":      1280,
		"upload_media_width":       720,
		"direct_v2":                1,
		"rotate":                   3,
		"xsharing_user_ids":        "[]",
		"hflip":                    false,
		"upload_media_duration_ms": common.GenNumber(1000, 5000),
		"upload_id":                upId,
		"retry_context":            retryContext,
		"media_type":               2})

	paramsStr := string(params)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(data)

	var resp = &RespRuploadIgPhoto{}
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUploadVideo + entityName,
		IsPost:  true,
		Header: map[string]string{
			"x_fb_video_waterfall_id":    this.inst.wid,
			"x-instagram-rupload-params": paramsStr,
			"content-type":               "application/octet-stream",
			"offset":                     "0",
			"segment-start-offset":       "0",
			"x-entity-name":              entityName,
			"x-entity-type":              "video/mp4",
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, &resp)
	err = resp.CheckError(err)
	return upId, err
}
