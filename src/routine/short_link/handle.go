package main

import (
	"github.com/gorilla/mux"
	"makemoney/common/log"
	"net/http"
)

func Redirect(vars map[string]string) []byte {
	data := RedirectHtmlData[vars["short_link"]]
	if data != nil {
		return data
	}

	var key string
	for k := range config.ShortLinkMap {
		key = k
		break
	}
	return RedirectHtmlData[key]
}

func doHttpLog(vars map[string]string, visitorType string, req *http.Request) {
	err := DoShortLinkLog2DB(&ShortLinkLogDB{
		UserID:      vars["user_id"],
		ShortLink:   vars["short_link"],
		Url:         req.RequestURI,
		UA:          req.UserAgent(),
		RemoteIP:    req.RemoteAddr,
		RequestHost: req.Host,
		VisitorType: visitorType,
		ReqHeader:   req.Header,
	})

	if err != nil {
		log.Error("save log error: %v", err)
	}
}

func HttpHandleShortLink(write http.ResponseWriter, req *http.Request) {
	write.WriteHeader(200)
	vars := mux.Vars(req)
	visitor := CheckVisitor(req, vars)
	if req.RequestURI != "/favicon.ico" {
		switch visitor {
		case "black":
			write.Write([]byte("fuck your mather?"))
			break
		case "fb":
			write.Write(FakeHtmlData)
			break
		case "ins":
			write.Write(Redirect(vars))
			break
		case "other":
			write.Write(Redirect(vars))
			break
		}
	}
	log.Info("ip %s url %s visitor: %s", req.RemoteAddr, req.RequestURI, visitor)
	doHttpLog(vars, visitor, req)
}

func HttpHandleLog(write http.ResponseWriter, req *http.Request) {
	log.Info("ip %s url %s send log", req.RemoteAddr, req.RequestURI)
	slog := &RedirectLog{
		UA:        req.UserAgent(),
		RemoteIP:  req.RemoteAddr,
		Url:       req.RequestURI,
		ReqHeader: req.Header,
	}
	DoRedirectLog(slog)
}
