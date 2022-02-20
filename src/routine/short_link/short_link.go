package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Config struct {
	ShortLink []struct {
		Key  string `json:"key"`
		Link string `json:"link"`
	} `json:"short_link"`
	HtmlPath     string   `json:"html_path"`
	MogoUri      string   `json:"mogo_uri"`
	Black        []Black  `json:"black"`
	Hosts        []string `json:"hosts"`
	ShortLinkMap map[string]string
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

var App ShortLinkApp
var config Config
var htmlData []byte
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

func doHttpLog(vars map[string]string, isFb bool, isBlack bool, req *http.Request) {
	err := ShortLinkLog2DB(&ShortLinkLogDB{
		UserID:    vars["user_id"],
		ShortLink: vars["short_link"],
		Url:       req.RequestURI,
		UA:        req.UserAgent(),
		IP:        req.RemoteAddr,
		Host:      req.Host,
		IsFb:      isFb,
		IsBlack:   isBlack,
		ReqHeader: req.Header,
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
		AddBlack(IP, "host", req)
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

func CheckFB(req *http.Request) bool {
	if req.Header.Get("X-Fb-Crawlerbot") != "" {
		return true
	}
	if strings.Contains(req.RequestURI, "fbclid") {
		return true
	}
	if !strings.Contains(req.UserAgent(), "Instagram") {
		return true
	}
	return false
}

func (this *ShortLinkApp) ServeHTTP(write http.ResponseWriter, req *http.Request) {
	var isFB = false
	var isBlack = false
	vars := mux.Vars(req)
	if CheckBlack(req, vars) {
		write.Write([]byte("fuck your mather?"))
		isBlack = true
	} else {
		if CheckFB(req) {
			_, err := write.Write(htmlData)
			if err != nil {
				log.Error("write html body error: %v", err)
			}
			write.WriteHeader(200)
			isFB = true
		} else {
			url := config.ShortLinkMap[vars["short_link"]]
			if url == "" {
				for _, v := range config.ShortLinkMap {
					url = v
					break
				}
			}
			http.Redirect(write, req, url, http.StatusTemporaryRedirect)
		}
	}
	doHttpLog(vars, isFB, isBlack, req)
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
	App.Router.Handle("/{short_link}/{user_id}", &App)
	App.Router.Handle("/{short_link}", &App)
	App.Router.PathPrefix("/").Handler(&App)

	err = http.ListenAndServe(":80", App.Router)
	if err != nil {
		log.Error("listen error: %v", err)
	}
}
