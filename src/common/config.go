package common

var IsDebug = true
var UseCharles = true
var UseTruncation = true

var DefaultHttpProxy = &Proxy{
	Ip:        "127.0.0.1",
	Port:      "8080",
	Username:  "",
	Passwd:    "",
	ProxyType: ProxyHttp,
	NeedAuth:  false,
}

var DefaultSocksProxy = &Proxy{
	Ip:        "127.0.0.1",
	Port:      "8889",
	Username:  "",
	Passwd:    "",
	ProxyType: ProxySocket,
	NeedAuth:  false,
}
