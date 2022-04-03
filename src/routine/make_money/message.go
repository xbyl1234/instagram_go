package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"sync/atomic"
)

var UserChan = make(chan *routine.UserComb, 100)

var SendSuccessCount int32
var SendErrorCount int32

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Data    []byte
	MD5     string
}

const (
	TexeMsg  = "text"
	TexeLink = "link"
	VoiceMsg = "voice"
	ImgMsg   = "img"
)

func UploadRes(inst *goinsta.Instagram, msg *Message) (string, error) {
	var mateData map[string]string
	if inst.MatePoint == nil {
		mateData = make(map[string]string)
		inst.MatePoint = mateData
		uploads, err := goinsta.LoadUploadID(inst.ID)
		if err != nil {
			log.Error("load user: %s uploads error: %v", inst.User, err)
		} else {
			for _, item := range uploads {
				mateData[item.FileMd5] = item.UploadID
			}
		}
	}

	mateData = inst.MatePoint.(map[string]string)
	uploadID := mateData[msg.MD5]

	if uploadID != "" {
		return uploadID, nil
	}

	var err error
	if msg.Type == ImgMsg {
		uploadID, _, err = inst.GetUpload().UploadPhoto(msg.Data, nil)
	} else if msg.Type == VoiceMsg {
		uploadID, _, err = inst.GetUpload().UploadVoice(msg.Data)
	}

	if err != nil {
		log.Error("account: %s, upload file %s, error: %v", inst.User, msg.Content, err)
		return "", err
	}

	_ = goinsta.SaveUploadID(&goinsta.UploadIDRecord{
		FileMd5:  msg.MD5,
		UserID:   inst.ID,
		FileType: msg.Type,
		UploadID: uploadID,
	})
	mateData[msg.MD5] = uploadID
	return uploadID, nil
}

func SendTask(WaitAll *sync.WaitGroup) {
	defer WaitAll.Done()
	var err error
	var uploadID string

	for user := range UserChan {
		inst := routine.ReqAccount(goinsta.OperNameSendMsg, config.AccountTag)
		//236
		err = inst.GetUserOperate().LikeUser(user.User.ID)
		if err != nil {
			log.Error("account %s like %d error: %v", inst.User, user.User.ID, err)
			goinsta.AccountPool.ReleaseOne(inst)
			continue
		}

		_, err = inst.GetMessage().GetThreadId(user.User.ID)
		if err != nil {
			goinsta.AccountPool.ReleaseOne(inst)
			routine.SaveBlackUser(user)
			continue
		}

		message := config.Msgs[common.GenNumber(0, len(config.Msgs))]
		err = nil

		for _, item := range message {
			switch item.Type {
			case TexeMsg:
				err = inst.GetMessage().SendTextMessage(user.User.ID, item.Content)
				break
			case TexeLink:
				err = inst.GetMessage().SendLinkMessage(user.User.ID, item.Content+fmt.Sprintf("%d", user.User.ID))
				break
			case ImgMsg:
				uploadID, err = UploadRes(inst, item)
				if err == nil {
					err = inst.GetMessage().SendImgMessage(user.User.ID, uploadID)
				}
				break
			case VoiceMsg:
				uploadID, err = UploadRes(inst, item)
				if err == nil {
					err = inst.GetMessage().SendVoiceMessage(user.User.ID, uploadID)
				}
				break
			}
			if err != nil {
				log.Error("send %s msg %s to %d, error: %v", item.Type, inst.User, user.User.ID, err)
				break
			} else {
				log.Info("send %s msg %s to %d success!", item.Type, inst.User, user.User.ID)
			}
		}

		if err != nil {
			atomic.AddInt32(&SendErrorCount, 1)
			log.Info("%s send to %d finish with error!", inst.User, user.User.ID)
		} else {
			atomic.AddInt32(&SendSuccessCount, 1)
			log.Info("%s send to %d finish!", inst.User, user.User.ID)
			routine.SaveSendFlag(user, config.TaskName)
		}

		log.Info("count: %d", SendSuccessCount)
		goinsta.AccountPool.ReleaseOne(inst)
	}
}

func SendUser(WaitAll *sync.WaitGroup) {
	defer WaitAll.Done()
	//51082952034
	//for true {
	//	user := &routine.UserComb{
	//		User:        &goinsta.User{ID: 51082952034},
	//		Source:      "",
	//		Followes:    nil,
	//		SendHistory: nil,
	//		Black:       false,
	//	}
	//	UserChan <- user
	//}
	err := routine.LoadUser(config.TaskName, UserChan)
	if err != nil {
		log.Error("load user error: %v", err)
	}
	log.Info("load user finish!")
	close(UserChan)
}

func LoadScreenplay() error {
	var err error
	for _, items := range config.Msgs {
		for _, item := range items {
			if item.Type == ImgMsg || item.Type == VoiceMsg {
				item.Data, err = os.ReadFile(item.Content)
				if err != nil {
					log.Error("load image %s error: %v", item.Content, err)
					return err
				}
				h := md5.New()
				h.Write(item.Data)
				item.MD5 = hex.EncodeToString(h.Sum(nil))
			}
		}
	}
	return nil
}

func MessageTask() {
	var WaitAll sync.WaitGroup
	WaitAll.Add(config.CoroCount + 1)
	go SendUser(&WaitAll)
	for index := 0; index < config.CoroCount; index++ {
		go SendTask(&WaitAll)
	}
	WaitAll.Wait()
}
