package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"os"
	"strings"
	"sync"
)

//function encryptionCode(str){
//   var len=str.length;
//   var rs="";
//   for(var i=0;i<len;i++){
//          var k=str.substring(i,i+1);
//          rs+= (i==0?"":",")+str.charCodeAt(i);
//   }
//   return rs;
//}

type Config struct {
	ShortLink []struct {
		Key  string `json:"key"`
		Link string `json:"link"`
	} `json:"short_link"`
	FakeHtmlPath     string   `json:"fake_html_path"`
	RedirectHtmlPath string   `json:"redirect_html_path"`
	MogoUri          string   `json:"mogo_uri"`
	Black            []Black  `json:"black"`
	Hosts            []string `json:"hosts"`
	ShortLinkMap     map[string]string
}

type ShortLinkApp struct {
	Router *mux.Router
}

type Black struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type BlackHistory struct {
	IP     string `bson:"ip"`
	Path   string `bson:"path"`
	UA     string `bson:"ua"`
	Host   string `bson:"host"`
	Reason string `bson:"reason"`
}

type ShortLinkLogDB struct {
	TimeTick    int64               `bson:"time_tick"`
	Time        string              `bson:"time"`
	UserID      string              `bson:"user_id"`
	ShortLink   string              `bson:"short_link"`
	Url         string              `bson:"url"`
	UA          string              `bson:"ua"`
	IP          string              `bson:"ip"`
	Host        string              `bson:"host"`
	VisitorType string              `bson:"visitor_type"`
	ReqHeader   map[string][]string `bson:"req_header"`
}

var App ShortLinkApp
var config Config
var FakeHtmlData []byte
var RedirectHtmlData map[string][]byte
var blackHistory map[string]*BlackHistory
var historyLock sync.Mutex

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

func doHttpLog(vars map[string]string, visitorType string, req *http.Request) {
	err := ShortLinkLog2DB(&ShortLinkLogDB{
		UserID:      vars["user_id"],
		ShortLink:   vars["short_link"],
		Url:         req.RequestURI,
		UA:          req.UserAgent(),
		IP:          req.RemoteAddr,
		Host:        req.Host,
		VisitorType: visitorType,
		ReqHeader:   req.Header,
	})

	if err != nil {
		log.Error("save log error: %v", err)
	}
}

func AddBlack(ip string, reason string, req *http.Request) {
	black := &BlackHistory{
		IP:     ip,
		Path:   req.RequestURI,
		UA:     req.UserAgent(),
		Reason: reason,
		Host:   req.Host,
	}
	SaveBlackHistory(black)
	historyLock.Lock()
	blackHistory[ip] = black
	historyLock.Unlock()
}

func CheckBlack(req *http.Request, params map[string]string) bool {
	var IP string
	sp := strings.Split(req.RemoteAddr, ":")
	if len(sp) == 2 {
		IP = sp[0]
	}

	var hasHost = false
	for _, item := range config.Hosts {
		if item == req.Host {
			hasHost = true
		}
	}
	if !hasHost {
		log.Warn("ip %s host is black", IP)
		//AddBlack(IP, "host", req)
		//return true
	}

	if req.Method != "GET" {
		log.Warn("ip %s method is black", IP)
		AddBlack(IP, "method", req)
		return true
	}

	if req.UserAgent() == "" {
		log.Warn("ip %s ua is black", IP)
		AddBlack(IP, "ua", req)
		return true
	}

	for _, item := range config.Black {
		if item.Type == "ua" {
			if strings.Contains(req.UserAgent(), item.Data) {
				log.Warn("ip %s ua is black", IP)
				AddBlack(IP, "ua", req)
				return true
			}
		} else {
			if req.RequestURI == item.Data {
				log.Warn("ip %s path is black", IP)
				AddBlack(IP, "path", req)
				return true
			}
		}
	}

	for _, item := range blackHistory {
		if item.IP == IP {
			log.Warn("ip %s ip is black", IP)
			return true
		}
	}

	return false
}

func CheckFB(req *http.Request, params map[string]string) bool {
	if req.Header.Get("X-Fb-Crawlerbot") != "" {
		return true
	}

	url := strings.ToLower(req.RequestURI)
	if strings.Contains(url, "fbclid") {
		return true
	}

	ua := strings.ToLower(req.UserAgent())
	if strings.Contains(ua, "facebookexternalhit") {
		return true
	}
	if strings.Contains(ua, "www.facebook.com") {
		return true
	}
	return false
}

func IsNumberStr(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func CheckIsInstagram(req *http.Request, params map[string]string) bool {
	if params["user_id"] == "" {
		return false
	}
	if !IsNumberStr(params["user_id"]) {
		return false
	}

	find := false
	for _, item := range config.ShortLink {
		if item.Key == params["short_link"] {
			find = true
			break
		}
	}
	if !find {
		return false
	}

	if !strings.Contains(req.UserAgent(), "Instagram") {
		return false
	}
	return true
}

func CheckVisitor(req *http.Request, params map[string]string) string {
	if CheckBlack(req, params) {
		return "black"
	}
	if CheckFB(req, params) {
		return "fb"
	}
	if CheckIsInstagram(req, params) {
		return "ins"
	}
	return "other"
}

func (this *ShortLinkApp) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	write.WriteHeader(200)
	vars := mux.Vars(req)
	visitor := CheckVisitor(req, vars)
	switch visitor {
	case "black":
		visitor = "black"
		write.Write([]byte("fuck your mather?"))
		break
	case "fb":
		visitor = "fb"
		write.Write(FakeHtmlData)
		break
	case "ins":
		visitor = "ins"
		data := RedirectHtmlData[vars["short_link"]]
		if data == nil {
			log.Error("write ins html is null!")
			break
		}
		write.Write(data)
		break
	case "other":
		visitor = "other"
		var key string
		for k := range config.ShortLinkMap {
			key = k
			break
		}
		data := RedirectHtmlData[key]
		if data == nil {
			log.Error("write ins html is null!")
			break
		}
		write.Write(data)
		break
	}
	log.Info("visitor: %s", visitor)
	doHttpLog(vars, visitor, req)
}

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

	RedirectHtmlData = make(map[string][]byte)
	config.ShortLinkMap = make(map[string]string)
	for _, item := range config.ShortLink {
		config.ShortLinkMap[item.Key] = item.Link
		RedirectHtmlData[item.Key] = bytes.ReplaceAll(data, []byte("flag1"), []byte(item.Link))
	}

	App.Router = mux.NewRouter()
	App.Router.Handle("/{short_link}/{user_id}", &App)
	App.Router.Handle("/{short_link}", &App)
	App.Router.PathPrefix("/").Handler(&App)

	err = http.ListenAndServe(":80", App.Router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
