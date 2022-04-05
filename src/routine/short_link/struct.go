package main

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

type ShortLinkJsLogDB struct {
	TimeTick  int64               `bson:"time_tick"`
	Time      string              `bson:"time"`
	UA        string              `bson:"ua"`
	IP        string              `bson:"ip"`
	ReqHeader map[string][]string `bson:"req_header"`
}
