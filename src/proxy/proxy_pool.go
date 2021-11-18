package proxy

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Proxy struct {
	Country   string `json:"country"`
	Area      string `json:"area"`
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Passwd    string `json:"passwd"`
	Rip       string `json:"rip"`
	ProxyType int    `json:"proxy_type"`
	NeedAuth  bool   `json:"need_auth"`
}

func LoadProxy(path string) {
	file, err := os.Open(path)
	if err != nil {
		Error("读代理文件错误!")
		os.Exit(0)
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')

		line = strings.Trim(line, "\n\r ")
		sp := strings.Split(line, "\t")
		if len(sp) == 2 {
			ProxyList.PushBack(&Proxy{
				sp[0],
				sp[1],
				"",
				"",
				ProxyType,
				false})
		} else if len(sp) == 4 {
			ProxyList.PushBack(&Proxy{
				sp[0],
				sp[1],
				sp[2],
				sp[3],
				ProxyType,
				true})
		} else {
			Warn("行错误! %s", line)
			continue
		}

		if err != nil || io.EOF == err {
			break
		}
	}

	if ProxyList.Len() == 0 {
		Error("无代理!")
		os.Exit(0)
	}
	ProxyEmt = ProxyList.Front()
}
