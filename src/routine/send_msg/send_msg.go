package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	config2 "makemoney/config"
	"makemoney/goinsta"
	"makemoney/routine"
	"os"
	"sync"
	"sync/atomic"
)

type CrawConfig struct {
	TaskName             string   `json:"task_name"`
	CoroCount            int      `json:"coro_count"`
	ProxyPath            string   `json:"proxy_path"`
	TargetUserDB         string   `json:"target_user_db"`
	TargetUserCollection string   `json:"target_user_collection"`
	TargetUserSource     string   `json:"target_user_source"`
	TextMsg              []string `json:"text_msg"`
	ImageMsg             []string `json:"image_msg"`
	IntervalTimeAccount  int      `json:"interval_time_account"`
	IntervalTimeMsg      int      `json:"interval_time_msg"`
}

var config CrawConfig
var ImageData [][]byte
var ImageMd5 []string
var WaitAll sync.WaitGroup
var UserChan = make(chan *routine.UserComb, 100)

var SendSuccessCount int32
var SendErrorCount int32

func AutoReleaseErrorAccount(inst *goinsta.Instagram, err error) {

}

func SendMsg(inst *goinsta.Instagram, user *routine.UserComb) error {
	if inst.MatePoint == nil {
		mate := make(map[string]string)
		for index := range ImageMd5 {
			recordID, _ := goinsta.FindUploadID(inst.User, ImageMd5[index])
			if recordID != nil {
				mate[ImageMd5[index]] = recordID.UploadID
			}
		}
		inst.MatePoint = mate
	}
	if user.User.ID == 0 {
		goinsta.AccountPool.ReleaseOne(inst)
		return nil
	}

	var err error
	var imageIndex = 0
	for index := range config.TextMsg {
		message := inst.GetMessage()
		if config.TextMsg[index] == "image" {
			uploadID := inst.MatePoint.(map[string]string)[ImageMd5[imageIndex]]
			if uploadID == "" {
				uploadID, err = inst.GetUpload().RuploadPhoto(ImageData[imageIndex])
				if err != nil {
					AutoReleaseErrorAccount(inst, err)
					log.Error("account: %s, send to %d, error: %v", inst.User, user.User.ID, err)
					return err
				}
				_ = goinsta.SaveUploadID(&goinsta.UploadIDRecord{
					FileMd5:  ImageMd5[imageIndex],
					Username: user.User.Username,
					FileType: "img",
					FileName: config.ImageMsg[imageIndex],
					UploadID: uploadID,
				})
				inst.MatePoint.(map[string]string)[ImageMd5[imageIndex]] = uploadID
			}
			err = message.SendImgMessage(fmt.Sprintf("%d", user.User.ID), uploadID)
			if err != nil {
				AutoReleaseErrorAccount(inst, err)
				log.Error("account: %s, send img to %d, error: %v", inst.User, user.User.ID, err)
				return err
			}
			imageIndex++
		} else {
			err = message.SendTextMessage(fmt.Sprintf("%d", user.User.ID), config.TextMsg[index])
			if err != nil {
				AutoReleaseErrorAccount(inst, err)
				log.Error("account: %s, send to %d, error: %v", inst.User, user.User.ID, err)
				return err
			}
		}
	}

	return nil
}

func SendTask() {
	defer WaitAll.Done()

	for user := range UserChan {
		inst, err := routine.ReqAccount(true)
		if err != nil {
			log.Error("req account unknow error!")
			return
		}

		err = SendMsg(inst, user)
		if err != nil {
			atomic.AddInt32(&SendErrorCount, 1)
		} else {
			atomic.AddInt32(&SendSuccessCount, 1)
		}
		goinsta.AccountPool.CoolingOne(inst)
	}
}

func SendUser() {
	defer WaitAll.Done()

	for true {
		users, err := routine.LoadUser(config.TargetUserSource, config.TaskName, 100)
		if err != nil {
			log.Error("load user error: %v", err)
			break
		}
		if len(users) == 0 {
			log.Info("no more user to load")
			break
		}

		for index := range users {
			UserChan <- &users[index]
		}
	}

	close(UserChan)
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
	if config.TargetUserSource == "" {
		log.Error("TargetUserSource is null")
		os.Exit(0)
	}
	if len(config.TextMsg) == 0 {
		log.Error("TextMsg is null")
		os.Exit(0)
	}
	if len(config.ImageMsg) == 0 {
		log.Error("ImageMsg is null")
		os.Exit(0)
	}

	for index := range config.ImageMsg {
		ImageData[index], err = os.ReadFile(config.ImageMsg[index])
		if err != nil {
			log.Error("load image %s error: %v", config.ImageMsg[index], err)
			os.Exit(0)
		}
		h := md5.New()
		h.Write(ImageData[index])
		ImageMd5[index] = hex.EncodeToString(h.Sum(nil))
	}

	//if len(config.ImageMsgUploadID) == 0 {
	//	id := time.Now().Unix()
	//	config.ImageMsgUploadID = make([]string, len(config.ImageMsg))
	//	for index := range config.ImageMsgUploadID {
	//		config.ImageMsgUploadID[index] = strconv.FormatInt(id+int64(index), 10)
	//	}
	//	common.Dumps(*TaskConfigPath, config)
	//}
}

func main() {
	config2.UseCharles = false
	initParams()
	routine.InitRoutine(config.ProxyPath)
	routine.InitSendMsgDB(config.TargetUserDB, config.TargetUserCollection)
	WaitAll.Add(config.CoroCount + 1)
	go SendUser()
	for index := 0; index < config.CoroCount; index++ {
		go SendTask()
	}
	WaitAll.Wait()
}
