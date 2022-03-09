package main

import "makemoney/goinsta"

func SendShortVideo(inst *goinsta.Instagram, video *RawVideoMedia) {

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

	uploadVideo, waterfall, err := upload.UploadVideo(video.VideoData, &goinsta.VideoUploadParams{
		UploadParamsBase: goinsta.UploadParamsBase{
			ContentTags:     "source-type-library,landscape",
			XsharingUserIds: []string{},
			MediaType:       2,
			IsClipsVideo:    "1",
		},
		UploadMediaHeight:     video.High,
		UploadMediaWidth:      video.Width,
		UploadMediaDurationMs: video.Duration,
	})

	if err != nil {
		return
	}
	title, err := opt.VerifyOriginalAudioTitle(video.AudioTitle)
	if err != nil {
		return
	}
	upload.UploadPhoto(video.ImageData, &goinsta.ImageUploadParams{
		UploadParamsBase: goinsta.UploadParamsBase{
			IsClipsVideo:    "",
			UploadId:        uploadVideo,
			XsharingUserIds: nil,
			MediaType:       0,
			ContentTags:     "",
		},
		ImageCompression: "",
	})
	//ClipsAssets
	upload.UploadFinish(uploadVideo)
	opt.UpdateVideoWithQualityInfo(uploadVideo, &goinsta.QualityInfo{
		OriginalVideoCodec:       "",
		EncodedVideoCodec:        "",
		OriginalColorPrimaries:   "",
		OriginalWidth:            0,
		OriginalFrameRate:        0,
		OriginalTransferFunction: "",
		EncodedHeight:            0,
		OriginalBitRate:          0,
		EncodedColorPrimaries:    "",
		OriginalHeight:           0,
		EncodedBitRate:           0,
		EncodedFrameRate:         0,
		EncodedYcbcrMatrix:       "",
		OriginalYcbcrMatrix:      "",
		EncodedWidth:             0,
		MeasuredFrames:           nil,
		EncodedTransferFunction:  "",
	})

	opt.ConfigureToClips()
	configure_to_clips
}
