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

func GenUploadID() string {
	upId := fmt.Sprintf("%d", time.Now().UnixMicro())
	return upId[:len(upId)-1]
}

func (this *Upload) UploadPhotoFromPath(path string, params *ImageUploadParams) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	return this.UploadPhoto(data, params)
}

func (this *Upload) UploadVoiceFromPath(path string) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	return this.UploadVoice(data)
}

func (this *Upload) UploadVideoFromPath(path string, params *VideoUploadParams) (string, string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	return this.UploadVideo(data, params)
}

type UploadParamsBase struct {
	IsClipsVideo    string   `json:"is_clips_video,omitempty"`
	UploadId        string   `json:"upload_id,omitempty"`
	XsharingUserIds []string `json:"xsharing_user_ids,omitempty"`
	MediaType       int      `json:"media_type,omitempty"`
	ContentTags     string   `json:"content_tags,omitempty"`
	WaterfallId     string   `json:"-"`
}

type ImageCompression struct {
	LibName    string  `json:"lib_name,omitempty"`
	LibVersion string  `json:"lib_version,omitempty"`
	Quality    int     `json:"quality,omitempty"`
	Colorspace string  `json:"colorspace,omitempty"`
	Ssim       float64 `json:"ssim,omitempty"`
}

type ImageUploadParams struct {
	UploadParamsBase
	ImageCompression string `json:"image_compression,omitempty"`
}

type VideoUploadParams struct {
	UploadParamsBase
	UploadMediaHeight     int     `json:"upload_media_height,omitempty"`
	UploadMediaWidth      int     `json:"upload_media_width,omitempty"`
	UploadMediaDurationMs float64 `json:"upload_media_duration_ms,omitempty"`
}

const UploadImageMediaTypeImage = 1
const UploadImageMediaTypeVideo = 2

func (this *Upload) UploadPhoto(data []byte, params *ImageUploadParams) (string, string, error) {
	if params == nil {
		params = &ImageUploadParams{
			UploadParamsBase: UploadParamsBase{
				IsClipsVideo: "",
				MediaType:    UploadImageMediaTypeImage,
				ContentTags:  "",
			},
		}
	}
	if params.UploadId == "" {
		params.UploadId = GenUploadID()
	}
	if params.WaterfallId == "" {
		params.WaterfallId = common.GenString(common.CharSet_16_Num, 32)
	}
	if params.ImageCompression == "" {
		imageCompression, _ := json.Marshal(&ImageCompression{
			LibName:    "uikit",
			LibVersion: "1575.230000",
			Quality:    64,
			Colorspace: "kCGColorSpaceDeviceRGB",
			Ssim:       0.96855711936950684,
		})
		params.ImageCompression = string(imageCompression)
	}
	if params.XsharingUserIds == nil {
		params.XsharingUserIds = []string{}
	}

	paramsStr, _ := json.Marshal(params)

	body := bytes.NewBuffer(data)
	var resp = &RespUpload{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        urlUploadPhone + common.GenString(common.CharSet_16_Num, 32),
		HeaderSequence: LoginHeaderMap[urlUploadPhone],
		IsPost:         true,
		Header: map[string]string{
			"x_fb_photo_waterfall_id":    params.WaterfallId,
			"x-instagram-rupload-params": string(paramsStr),
			"content-type":               "application/octet-stream",
			"x-entity-type":              "image/jpeg",
			"x-entity-name":              "image.jpeg",
			"x-entity-length":            strconv.Itoa(len(data)),
			"offset":                     "0",
			"Media_hash":                 "",
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	return params.UploadId, params.WaterfallId, err
}

func (this *Upload) UploadVideo(data []byte, params *VideoUploadParams) (string, string, error) {
	//retryContext, _ := json.Marshal(map[string]interface{}{
	//	"num_reupload":          0,
	//	"num_step_auto_retry":   0,
	//	"num_step_manual_retry": 0,
	//})
	//params, _ := json.Marshal(map[string]interface{}{
	//	"upload_media_height":      1280,
	//	"upload_media_width":       720,
	//	"direct_v2":                1,
	//	"rotate":                   3,
	//	"xsharing_user_ids":        "[]",
	//	"hflip":                    false,
	//	"upload_media_duration_ms": common.GenNumber(1000, 5000),
	//	"upload_id":                upId,
	//	"retry_context":            retryContext,
	//	"media_type":               2})

	waterfall := common.GenString(common.CharSet_16_Num, 32)
	path := common.GenUUID()
	params.UploadId = GenUploadID()
	paramsStr, _ := json.Marshal(params)

	body := bytes.NewBuffer(data)

	var resp = &RespUpload{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        urlUploadVideo + path,
		HeaderSequence: LoginHeaderMap[urlUploadVideo],
		IsPost:         true,
		Header: map[string]string{
			"x_fb_video_waterfall_id":    waterfall,
			"x-instagram-rupload-params": string(paramsStr),
			"content-type":               "application/octet-stream",
			"offset":                     "0",
			"segment-start-offset":       "0",
			"x-entity-name":              "video.mp4",
			"x-entity-type":              "video/mp4",
			"x-entity-length":            strconv.Itoa(len(data)),
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	return params.UploadId, waterfall, err
}

func (this *Upload) UploadFinish(uploadID string) error {
	var resp = &BaseApiResp{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: urlUploadFinish,
		IsPost:  true,
		Signed:  true,
		Query: map[string]interface{}{
			"upload_id": uploadID,
			"_uuid":     this.inst.AccountInfo.Device.DeviceID,
			"_uid":      fmt.Sprintf("%d", this.inst.ID),
		},
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *Upload) UploadVoice(data []byte) (string, string, error) {
	upId := GenUploadID()
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
	waterfall := common.GenString(common.CharSet_16_Num, 32)
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
			"X_fb_video_waterfall_id":    waterfall,
			"X-Instagram-Rupload-Params": uploadParamsStr,
		},
		Body: body,
	}, resp)
	err = resp.CheckError(err)
	if err != nil {
		return "", waterfall, err
	}
	err = this.UploadFinish(upId)
	return upId, waterfall, err
}
