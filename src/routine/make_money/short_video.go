package main

import (
	"makemoney/common/log"
	"makemoney/goinsta"
)

func SendShortVideo(inst *goinsta.Instagram, video *goinsta.RawVideoMedia) {
	opt := inst.GetUserOperate()
	upload := inst.GetUpload()

	err := opt.ClipsInfoForCreation()
	if err != nil {
		return
	}
	assets, err := opt.ClipsAssets(video.Latitude, video.Longitude)
	if err != nil {
		return
	}
	print(assets)

	uploadVideo, waterfallVideo, err := upload.UploadVideo(video.VideoData, &goinsta.VideoUploadParams{
		UploadParamsBase: goinsta.UploadParamsBase{
			ContentTags:     "source-type-library,landscape",
			XsharingUserIds: []string{},
			MediaType:       goinsta.UploadImageMediaTypeVideo,
			IsClipsVideo:    "1",
		},
		UploadMediaHeight:     video.High,
		UploadMediaWidth:      video.Width,
		UploadMediaDurationMs: video.Duration,
	})
	if err != nil {
		log.Error("upload video error: %v", err)
		return
	}
	video.UploadId = uploadVideo
	video.Waterfall = waterfallVideo

	verifyTitle, err := opt.VerifyOriginalAudioTitle(video.AudioTitle)
	if err != nil || !verifyTitle.IsValid {
		log.Error("verify title error: %v", err)
		return
	}

	_, _, err = upload.UploadPhoto(video.ImageData, &goinsta.ImageUploadParams{
		UploadParamsBase: goinsta.UploadParamsBase{
			IsClipsVideo:    "1",
			UploadId:        uploadVideo,
			WaterfallId:     waterfallVideo,
			XsharingUserIds: nil,
			MediaType:       goinsta.UploadImageMediaTypeVideo,
			ContentTags:     "portrait,source-type-library",
		},
		ImageCompression: "",
	})
	if err != nil {
		log.Error("upload cover error: %v", err)
		return
	}

	//ClipsAssets
	err = upload.UploadFinish(uploadVideo)
	if err != nil {
		log.Error("upload finish error: %v", err)
		return
	}
	frames := make([]goinsta.MeasuredFrames, int(video.Duration/1000/0.9))
	for idx := range frames {
		frames[idx].Ssim = 0.95175731182098389
		frames[idx].Timestamp = float64(idx) * 0.9
	}

	//err = opt.UpdateVideoWithQualityInfo(uploadVideo, &goinsta.QualityInfo{
	//	OriginalVideoCodec:       video.VideoCodec,
	//	EncodedVideoCodec:        video.VideoCodec,
	//	OriginalColorPrimaries:   video.YcbcrMatrix,
	//	OriginalWidth:            video.Width,
	//	OriginalFrameRate:        video.FrameRate,
	//	OriginalTransferFunction: video.YcbcrMatrix,
	//	EncodedHeight:            video.High,
	//	OriginalBitRate:          int(video.BitRate),
	//	EncodedColorPrimaries:    video.YcbcrMatrix,
	//	OriginalHeight:           video.High,
	//	EncodedBitRate:           video.BitRate,
	//	EncodedFrameRate:         video.FrameRate,
	//	EncodedYcbcrMatrix:       video.YcbcrMatrix,
	//	OriginalYcbcrMatrix:      video.YcbcrMatrix,
	//	EncodedWidth:             video.Width,
	//	MeasuredFrames:           frames,
	//	EncodedTransferFunction:  video.YcbcrMatrix,
	//})
	//if err != nil {
	//	log.Error("upload video with quality error: %v", err)
	//	return
	//}

	clips, err := opt.ConfigureToClips(video)
	if err != nil {
		log.Error("configure to clips error: %v", err)
		return
	}
	print(clips)
}

func ShortVideoTask() {
	//inst := routine.ReqAccount(goinsta.OperNameSendMsg, config.AccountTag)
	//rawMedia := &goinsta.RawVideoMedia{
	//	Caption:    "Are you ready for the boys of summer",
	//	AudioTitle: "Like and follow",
	//	From:       goinsta.FromCamera,
	//}
	//
	//rawMedia.LoadVideo("C:\\Users\\Administrator\\Desktop\\mn\\test.mp4",
	//	"C:\\Users\\Administrator\\Desktop\\mn\\暴风截图2022310297691375.jpg")
	//SendShortVideo(inst, rawMedia)

	for _, inst := range goinsta.AccountPool.Accounts {
		rawMedia := &goinsta.RawVideoMedia{
			Caption:    "Are you ready for the boys of summer #test5555555555555555555",
			AudioTitle: "Like and follow",
			From:       goinsta.FromCamera,
		}

		rawMedia.LoadVideo("C:\\Users\\Administrator\\Desktop\\mn\\test.mp4",
			"C:\\Users\\Administrator\\Desktop\\mn\\暴风截图2022310297691375.jpg")
		SendShortVideo(inst, rawMedia)
	}
}
