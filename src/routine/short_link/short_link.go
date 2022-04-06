package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"os"
	"sync"
)

var Router *mux.Router
var config Config
var FakeHtmlData []byte
var RedirectHtmlData map[string][]byte
var blackHistory map[string]*BlackHistory
var historyLock sync.Mutex

func main() {
	log.InitDefaultLog("short_link", true, true)
	err := common.LoadJsonFile("./short_link.json", &config)
	if err != nil {
		log.Error("load config error: %v", err)
		return
	}
	if config.MogoUri == "" {
		log.Error("config MogoUri is null!")
		return
	}
	InitShortLinkDB(config.MogoUri)
	history, err := LoadBlackHistory()
	if err != nil {
		log.Error("load black history error: %v", err)
		return
	}
	blackHistory = make(map[string]*BlackHistory)
	for _, item := range history {
		blackHistory[item.IP] = item
	}

	if config.FakeHtmlPath == "" {
		config.FakeHtmlPath = "./fake.html"
	}
	if config.FakeHtmlPath == "" {
		config.FakeHtmlPath = "./redirect.html"
	}
	FakeHtmlData, err = os.ReadFile(config.FakeHtmlPath)
	if err != nil {
		log.Error("load %s error: %v", config.FakeHtmlPath, err)
		return
	}
	data, err := os.ReadFile(config.RedirectHtmlPath)
	if err != nil {
		log.Error("load %s error: %v", config.RedirectHtmlPath, err)
		return
	}
	data = bytes.ReplaceAll(data, []byte("flag2"), []byte(encryptionCode(config.LogUrl)))

	RedirectHtmlData = make(map[string][]byte)
	config.ShortLinkMap = make(map[string]string)
	for _, item := range config.ShortLink {
		config.ShortLinkMap[item.Key] = encryptionCode(item.Link)
		RedirectHtmlData[item.Key] = bytes.ReplaceAll(data, []byte("flag1"), []byte(encryptionCode(item.Link)))
	}

	Router = mux.NewRouter()
	Router.HandleFunc("/log", HttpHandleLog)
	Router.HandleFunc("/{short_link}/{user_id}", HttpHandleShortLink)
	Router.HandleFunc("/{short_link}", HttpHandleShortLink)
	Router.PathPrefix("/").HandlerFunc(HttpHandleShortLink)

	err = http.ListenAndServe(":80", Router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
