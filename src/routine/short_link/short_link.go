package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	ShortLink []struct {
		Key  string `json:"key"`
		Link string `json:"link"`
	} `json:"short_link"`
	HtmlPath     string `json:"html_path"`
	MogoUri      string `json:"mogo_uri"`
	ShortLinkMap map[string]string
}

type ShortLinkApp struct {
	Router *mux.Router
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("http: %s\n", r.RemoteAddr)
	fmt.Printf("get: %s\n", r.URL.RequestURI())
	fmt.Printf("header:\n")
	for k, v := range r.Header {
		fmt.Printf("---\t%s\t\t:\t%s\n", k, v)
	}
	fmt.Printf("\n\n")
	w.Write([]byte("123"))
}

//time ip is_fb url ua

var logFmtStr = "%d\t%s\t%v\t%s\t%s"

func doHttpLog(ip string, isFb bool, url string, ua string) {
	if isFb {
		httpLog.Warn(logFmtStr, time.Now().Unix(), ip, isFb, url, ua)
	} else {
		httpLog.Info(logFmtStr, time.Now().Unix(), ip, isFb, url, ua)
	}
	//log.Info(logFmtStr, time.Now().Unix(), ip, isFb, url, ua)
}

func (this *ShortLinkApp) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if req.Header.Get("X-Fb-Crawlerbot") != "" {
		_, err := write.Write(htmlData)
		if err != nil {
			log.Error("write html body error: %v", err)
		}
		write.WriteHeader(200)
		doHttpLog(req.RemoteAddr, true, req.URL.RequestURI(), req.Header.Get("User-Agent"))
		return
	}
	url := config.ShortLinkMap[vars["short_link"]]
	if url == "" {
		for _, v := range config.ShortLinkMap {
			url = v
			break
		}
	}
	http.Redirect(write, req, url, http.StatusTemporaryRedirect)
	doHttpLog(req.RemoteAddr, false, req.URL.RequestURI(), req.Header.Get("User-Agent"))
	err := ShortLinkLog2DB(&ShortLinkLogDB{
		UserID:    vars["user_id"],
		ShortLink: vars["short_link"],
		Url:       req.URL.RequestURI(),
		UA:        req.Header.Get("User-Agent"),
		IP:        req.RemoteAddr,
	})
	if err != nil {
		log.Error("save log error: %v", err)
	}
}

var App ShortLinkApp
var config Config
var htmlData []byte
var httpLog *log.Log

func main() {
	log.InitDefaultLog("short_link", true, true)
	httpLog = log.NewFileLog("http")

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
	if config.HtmlPath == "" {
		config.HtmlPath = "./fake.html"
	}
	config.ShortLinkMap = make(map[string]string)
	for _, item := range config.ShortLink {
		config.ShortLinkMap[item.Key] = item.Link
	}

	htmlData, err = os.ReadFile(config.HtmlPath)
	if err != nil {
		log.Error("load %s error: %v", config.HtmlPath, err)
	}

	App.Router = mux.NewRouter()
	App.Router.Handle("/{short_link:[a-zA-Z0-9]{0,}}/{user_id}", &App)
	App.Router.Handle("/{short_link:[a-zA-Z0-9]{0,}}", &App)
	App.Router.PathPrefix("/").Handler(&App)

	err = http.ListenAndServe(":80", App.Router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
