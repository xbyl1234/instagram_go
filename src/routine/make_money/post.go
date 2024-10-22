package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
)

func Post(inst *goinsta.Instagram, rawMedia *goinsta.RawImgMedia) {
	upload := inst.GetUpload()
	uploadID, waterfall, err := upload.UploadPhoto(rawMedia.GetImage(), nil)
	if err != nil {
		log.Error("account: %s, error: %v", inst.User, err)
		return
	}

	userOpt := inst.GetUserOperate()
	search, err := userOpt.LocationSearch(-118.25185775756837, 34.047112447489795)
	if err != nil {
		log.Error("account: %s, error: %v", inst.User, err)
		return
	}
	rawMedia.Loc = &search.Venues[common.GenNumber(0, len(search.Venues))]

	mediaInfo := &goinsta.UploadMediaInfo{
		UploadID:  uploadID,
		Waterfall: waterfall,
		High:      rawMedia.High,
		Width:     rawMedia.Width,
	}

	post, err := userOpt.ConfigurePost(rawMedia.Caption, mediaInfo, uploadID, rawMedia.Loc)
	if err != nil {
		log.Error("account: %s, error: %v", inst.User, err)
		return
	}
	print(post)
}

func PostTask() {
	//inst := routine.ReqAccount(goinsta.OperNameSendMsg, config.AccountTag)
	for _, inst := range goinsta.AccountPool.Accounts {
		rawMedia := &goinsta.RawImgMedia{
			Caption: "hi boy,this is me, do you love me?#test5555555555555555555",
			Loc:     nil,
		}
		rawMedia.LoadImage("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture\\1e60dd4625408b8aac26f7169a1e5deb.jpg")
		Post(inst, rawMedia)
	}
}
