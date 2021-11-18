package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
)

type Proxy struct {
	Country  string `json:"country"`
	Area     string `json:"area"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	Rip      string `json:"rip"`
}

var ProxyMap = map[string][]Proxy{}

func loadIP(path string) *list.List {
	ProxyList := list.New()

	file, err := os.Open(path)
	if err != nil {
		os.Exit(0)
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')

		line = strings.Trim(line, "\n\r ")
		sp := strings.Split(line, ":")
		rip := strings.ReplaceAll(sp[2], "lum-customer-hl_28871e6d-zone-zone2-ip-", "")
		if len(sp) == 4 {
			ProxyList.PushBack(&Proxy{"", "",
				sp[0],
				sp[1],
				sp[2],
				sp[3], rip})
		} else {
			continue
		}

		if err != nil || io.EOF == err {
			break
		}
	}

	if ProxyList.Len() == 0 {
		os.Exit(0)
	}
	return ProxyList
}

func CheckIp(ipList *list.List) {
	for item := ipList.Front(); item != nil; item = item.Next() {
		ip := item.Value.(*Proxy)
		country, area := findIP(ip.Rip)
		ip.Country = country
		ip.Area = area

		if ProxyMap[country] == nil {
			ProxyMap[country] = []Proxy{}
		}
		ProxyMap[country] = append(ProxyMap[country], *ip)
	}
}

func main() {
	IPData.FilePath = "./qqwry.dat"
	res := IPData.InitIPData()

	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	ips := loadIP("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\ips-zone2.txt")
	CheckIp(ips)

	f, _ := os.Create("C:\\Users\\Administrator\\Desktop\\project\\github\\instagram_project\\data\\ips-zone2.json")
	s, _ := json.Marshal(ProxyMap)
	f.Write(s)
}

// findIP 查找 IP 地址的接口
func findIP(ip string) (string, string) {
	qqWry := NewQQwry()
	ret := qqWry.Find(ip)
	return ret.Country, ret.Area
}
