package tools

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"makemoney/config"
	"net/http"
	neturl "net/url"
)

func LoadJsonFile(path string, ret interface{}) error {
	by, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(by, ret)
	return err
}

func DebugHttpClient(clinet *http.Client) {
	if config.UseCharles {
		uri, err := neturl.Parse("http://127.0.0.1:8888")
		if err == nil {
			clinet.Transport = &http.Transport{
				Proxy: http.ProxyURL(uri),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}
		}
	}
}
