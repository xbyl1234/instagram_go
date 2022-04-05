package main

import (
	"fmt"
	"makemoney/common/log"
	"net/http"
	"strings"
)

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

func encryptionCode(link string) string {
	var ret string
	for _, c := range link {
		ret += fmt.Sprintf("%d,", c)
	}
	return ret[:len(ret)-1]
}
