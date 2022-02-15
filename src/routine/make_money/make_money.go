package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type MakeMoneyConfig struct {
	TaskName             string       `json:"task_name"`
	CoroCount            int          `json:"coro_count"`
	ProxyPath            string       `json:"proxy_path"`
	TargetUserDB         string       `json:"target_user_db"`
	TargetUserCollection string       `json:"target_user_collection"`
	AccountTag           string       `json:"account_tag"`
	Msgs                 [][]*Message `json:"msgs"`
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Data    []byte
	MD5     string
}

const (
	TexeMsg  = "text"
	VoiceMsg = "voice"
	ImgMsg   = "img"
)

var config MakeMoneyConfig
var WaitAll sync.WaitGroup
var UserChan = make(chan *routine.UserComb, 100)

var SendSuccessCount int32
var SendErrorCount int32

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
		uploadID, err = inst.GetUpload().UploadPhoto(msg.Data)
	} else if msg.Type == VoiceMsg {
		uploadID, err = inst.GetUpload().UploadVoice(msg.Data)
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

func SendTask() {
	defer WaitAll.Done()
	var err error
	var uploadID string

	for user := range UserChan {
		inst := routine.ReqAccount(goinsta.OperNameSendMsg, config.AccountTag)
		message := config.Msgs[common.GenNumber(0, len(config.Msgs))]
		err = nil
		for _, item := range message {
			switch item.Type {
			case TexeMsg:
				err = inst.GetMessage().SendTextMessage(user.User.ID, item.Content)
				break
			case ImgMsg:
				uploadID, err = UploadRes(inst, item)
				if err == nil {
					//err = inst.GetMessage().SendImgMessage(user.User.ID, uploadID)
					err = inst.GetAccount().ChangeProfilePicture(uploadID)
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

		routine.SaveSendFlag(routine.SendTargeUserColl, user, config.TaskName)

		if err != nil {
			atomic.AddInt32(&SendErrorCount, 1)
			log.Info("%s send to %d finish with error!", inst.User, user.User.ID)
		} else {
			atomic.AddInt32(&SendSuccessCount, 1)
			log.Info("%s send to %d finish!", inst.User, user.User.ID)
		}
		goinsta.AccountPool.ReleaseOne(inst)
	}
}

func SendUser() {
	defer WaitAll.Done()
	//51082952034
	for true {
		user := &routine.UserComb{
			User:        &goinsta.User{ID: 51082952034},
			Source:      "",
			Followes:    nil,
			SendHistory: nil,
			Black:       false,
		}
		UserChan <- user
	}

	for true {
		users, err := routine.LoadUser(config.TaskName, 100)
		if err != nil {
			log.Error("load user error: %v", err)
			time.Sleep(time.Minute)
			continue
		}
		if len(users) == 0 {
			log.Info("no more user to load")
			time.Sleep(time.Minute)
			continue
		}

		for index := range users {
			UserChan <- &users[index]
		}
	}

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

func initParams() {
	var err error
	var TaskConfigPath = flag.String("task", "", "task")
	log.InitDefaultLog("send_msg", true, true)
	flag.Parse()
	if *TaskConfigPath == "" {
		log.Error("task config path is null!")
		os.Exit(0)
	}
	err = common.LoadJsonFile(*TaskConfigPath, &config)
	if err != nil {
		log.Error("load task config error: %v", err)
		os.Exit(0)
	}
	if config.TaskName == "" {
		log.Error("task name is null")
		os.Exit(0)
	}
	if config.CoroCount == 0 {
		config.CoroCount = 1
	}
	if config.TargetUserDB == "" {
		log.Error("TargetUserDB name is null")
		os.Exit(0)
	}
	if config.TargetUserCollection == "" {
		log.Error("TargetUserCollection is null")
		os.Exit(0)
	}
	if config.Msgs == nil || len(config.Msgs) == 0 {
		log.Error("Msgs is null")
		os.Exit(0)
	}
	err = LoadScreenplay()
	if err != nil {
		log.Error("LoadScreenplay error: %v", err)
		os.Exit(0)
	}
}

func main() {
	config2.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitSendMsgDB(config.TargetUserDB, config.TargetUserCollection)
	intas := goinsta.LoadAccountByTags([]string{config.AccountTag, config.AccountTag})
	if len(intas) == 0 {
		log.Warn("there have no account!")
	} else {
		goinsta.InitAccountPool(intas)
	}

	WaitAll.Add(config.CoroCount + 1)
	go SendUser()
	for index := 0; index < config.CoroCount; index++ {
		go SendTask()
	}
	WaitAll.Wait()
}
