package main

import (
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
)

func Post(inst *goinsta.Instagram, rawMedia *RawImgMedia) {
	upload := inst.GetUpload()
	uploadID, waterfall, err := upload.UploadPhoto(rawMedia.GetImage())
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
		rawMedia := &RawImgMedia{
			Caption: "hi boy,this is me, do you love me?",
			Loc:     nil,
		}
		rawMedia.LoadImage("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\girl_picture\\0ae56c901f796d10f8c944b7d5612daf.jpg")
		Post(inst, rawMedia)
	}
}
