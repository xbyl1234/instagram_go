package main

import (
	"fmt"
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"strings"
	"time"
)

func SendShortVideo(inst *goinsta.Instagram, video *goinsta.RawVideoMedia) {
	opt := inst.GetUserOperate()
	upload := inst.GetUpload()

	err := opt.ClipsInfoForCreation()
	if err != nil {
		return
	}
	_, err = opt.ClipsAssets(video.Latitude, video.Longitude)
	if err != nil {
		return
	}

	uploadVideo, waterfallVideo, err := upload.UploadVideo(video.VideoData, &goinsta.VideoUploadParams{
		UploadParamsBase: goinsta.UploadParamsBase{
			ContentTags:     "source-type-library,landscape",
			XsharingUserIds: []string{},
			MediaType:       goinsta.UploadImageMediaTypeVideo,
			IsClipsVideo:    "1",
			UploadId:        fmt.Sprintf("%d", time.Now().UnixMicro()),
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
		if float64(idx)*0.9 > video.Duration {
			break
		}
		frames[idx].Ssim = 0.95175731182098389
		frames[idx].Timestamp = float64(idx) * 0.9
	}

	err = opt.UpdateVideoWithQualityInfo(uploadVideo, &goinsta.QualityInfo{
		OriginalVideoCodec: video.VideoCodec,
		EncodedVideoCodec:  video.VideoCodec,
		//OriginalColorPrimaries: video.YcbcrMatrix,
		OriginalWidth:     video.Width,
		OriginalFrameRate: video.FrameRate,
		//OriginalTransferFunction: video.YcbcrMatrix,
		EncodedHeight:           video.High,
		OriginalBitRate:         video.BitRate,
		EncodedColorPrimaries:   video.YcbcrMatrix,
		OriginalHeight:          video.High,
		EncodedBitRate:          video.BitRate,
		EncodedFrameRate:        video.FrameRate,
		EncodedYcbcrMatrix:      video.YcbcrMatrix,
		OriginalYcbcrMatrix:     video.YcbcrMatrix,
		EncodedWidth:            video.Width,
		MeasuredFrames:          frames,
		EncodedTransferFunction: video.YcbcrMatrix,
	})
	if err != nil {
		log.Error("upload video with quality error: %v", err)
		//return
	}

	clips, err := opt.ConfigureToClips(video)
	if err != nil {
		log.Error("configure to clips error: %v", err)
		return
	}
	print(clips)
}

var RawMedias []*goinsta.RawVideoMedia

func LoadVideo(path string) {
	dir, err := ioutil.ReadDir(path + "video/")
	if err != nil {
		log.Error("load video error: %v", err)
		return
	}
	RawMedias = make([]*goinsta.RawVideoMedia, len(dir))
	index := 0
	for _, item := range dir {
		if !item.IsDir() && strings.Contains(item.Name(), ".mp4") {
			rawMedia := &goinsta.RawVideoMedia{
				Caption:    "can you give me a start? #fashion #followme #like4like #love#test5555555555555555555 " + common.GenString(common.CharSet_All, 5),
				AudioTitle: "like and follow please~",
				From:       goinsta.FromCamera,
			}
			name := strings.ReplaceAll(item.Name(), ".mp4", "")
			err = rawMedia.LoadVideo(path+"video/"+item.Name(),
				path+"cover/"+name+".jpg")
			if err != nil {
				log.Error("load video %s error: %v", item.Name(), err)
				continue
			}
			RawMedias[index] = rawMedia
			index++
		}
	}
	RawMedias = RawMedias[:index]
}

func ShortVideoTask() {
	//inst := routine.ReqAccount(goinsta.OperNameSendMsg, config.AccountTag)
	//rawMedia := &goinsta.RawVideoMedia{
	//	Caption:    "Are you ready for the boys of summer",
	//	AudioTitle: "Like and follow",
	//	From:       goinsta.FromCamera,
	//}
	//rawMedia.LoadVideo("C:\\Users\\Administrator\\Desktop\\mn\\test.mp4",
	//	"C:\\Users\\Administrator\\Desktop\\mn\\暴风截图2022310297691375.jpg")
	//SendShortVideo(inst, rawMedia)
	LoadVideo("C:/Users/Administrator/Desktop/mn/刘二/")

	for idx, inst := range goinsta.AccountPool.Accounts {
		//rawMedia := &goinsta.RawVideoMedia{
		//	Caption:    "can you give me a start? #fashion #followme #like4like #love#test5555555555555555555",
		//	AudioTitle: "Like and follow",
		//	From:       goinsta.FromLibrary,
		//}
		//
		//rawMedia.LoadVideo("C:\\Users\\Administrator\\Desktop\\mn\\刘二\\video\\"+video[idx]+".mp4",
		//	"C:\\Users\\Administrator\\Desktop\\mn\\刘二\\cover\\"+video[idx]+".jpg")
		rawMedias := RawMedias[idx]
		SendShortVideo(inst, rawMedias)
	}
}
