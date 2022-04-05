package main

import (
	"io/ioutil"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"strings"
)

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

	//for idx, inst := range goinsta.AccountPool.Accounts {
	//rawMedia := &goinsta.RawVideoMedia{
	//	Caption:    "can you give me a start? #fashion #followme #like4like #love#test5555555555555555555",
	//	AudioTitle: "Like and follow",
	//	From:       goinsta.FromLibrary,
	//}
	//
	//rawMedia.LoadVideo("C:\\Users\\Administrator\\Desktop\\mn\\刘二\\video\\"+video[idx]+".mp4",
	//	"C:\\Users\\Administrator\\Desktop\\mn\\刘二\\cover\\"+video[idx]+".jpg")
	//rawMedias := RawMedias[idx]
	//SendShortVideo(inst, rawMedias)
	//}
}
