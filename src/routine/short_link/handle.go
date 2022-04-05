package main

import (
	"github.com/gorilla/mux"
	"makemoney/common/log"
	"net/http"
)

func HttpHandleShortLink(write http.ResponseWriter, req *http.Request) {
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

func HttpHandleLog(write http.ResponseWriter, req *http.Request) {
	slog := &ShortLinkJsLogDB{
		UA:        req.UserAgent(),
		IP:        req.RemoteAddr,
		ReqHeader: req.Header,
	}
	DoShortLinkJsLogDB(slog)
}
