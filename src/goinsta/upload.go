package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type RespUpload struct {
	BaseApiResp
	Offset         int         `json:"offset"`
	UploadId       string      `json:"upload_id"`
	XsharingNonces interface{} `json:"xsharing_nonces"`
}

func (this *Upload) RuploadPhotoFromPath(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return this.RuploadPhoto(data)
}

func (this *Upload) RuploadVoiceFromPath(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return this.RuploadVoice(data)
}

func (this *Upload) RuploadPhoto(data []byte) (string, error) {
	upId := fmt.Sprintf("%d", time.Now().UnixMicro())
	path := common.GenString(common.CharSet_16_Num, 32)

	imageCompression, _ := json.Marshal(map[string]interface{}{
		"lib_name":    "uikit",
		"lib_version": "1575.230000",
		"quality":     64,
		"colorspace":  "kCGColorSpaceDeviceRGB",
		"ssim":        0.9,
	})

	params, _ := json.Marshal(map[string]interface{}{
		"upload_id":         upId,
		"media_type":        1,
		"image_compression": string(imageCompression),
	})
	paramsStr := string(params)

	body := bytes.NewBuffer(data)

	var resp = &RespUpload{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        urlUploadPhone + path,
		HeaderSequence: LoginHeaderMap[urlUploadPhone],
		IsPost:         true,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    common.GenString(common.CharSet_16_Num, 32),
			"x-instagram-rupload-params": paramsStr,
			"content-type":               "application/octet-stream",
			"x-entity-type":              "image/jpeg",
			"x-entity-name":              "image.jpeg",
			"x-entity-length":            strconv.Itoa(len(data)),
			"offset":                     "0",
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

	var resp = &RespUpload{}
	err = this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        urlUploadVideo + entityName,
		HeaderSequence: LoginHeaderMap[urlUploadVideo],
		IsPost:         true,
		Header: map[string]string{
			"x_fb_video_waterfall_id":    this.inst.Device.WaterID,
			"x-instagram-rupload-params": paramsStr,
			"content-type":               "application/octet-stream",
			"offset":                     "0",
			"segment-start-offset":       "0",
			"x-entity-name":              entityName,
			"x-entity-type":              "video/mp4",
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	return upId, err
}

func (this *Upload) RuploadVoice(data []byte) (string, error) {
	upId := fmt.Sprintf("%d", time.Now().UnixMicro())
	path := common.GenString(common.CharSet_16_Num, 32)

	uploadParams, _ := json.Marshal(map[string]interface{}{
		"upload_id":         "1644744220215219",
		"xsharing_user_ids": "[]",
		"media_type":        11,
		"is_direct_voice":   true,
	})

	uploadParamsStr := string(uploadParams)
	body := bytes.NewBuffer(data)

	var resp = &RespUpload{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        urlUploadVideo + path,
		HeaderSequence: LoginHeaderMap[urlUploadVideo],
		IsPost:         true,
		Header: map[string]string{
			"Media_hash":                 "",
			"X-Entity-Name":              "audio.m4a",
			"X-Entity-Type":              "audio/aac",
			"Offset":                     "0",
			"Content-Type":               "application/octet-stream",
			"X-Entity-Length":            fmt.Sprintf("%d", len(data)),
			"X_fb_video_waterfall_id":    common.GenString(common.CharSet_16_Num, 32),
			"X-Instagram-Rupload-Params": uploadParamsStr,
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	return upId, err
}
