package goinsta

// StoryReelMention represent story reel mention
type StoryReelMention struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int     `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	User     User
}

// StoryCTA represent story cta
type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

type StoryMedia struct {
	inst     *Instagram
	endpoint string
	uid      int64

	err error

	Pk              interface{} `json:"id"`
	LatestReelMedia int64       `json:"latest_reel_media"`
	ExpiringAt      float64     `json:"expiring_at"`
	HaveBeenSeen    float64     `json:"seen"`
	CanReply        bool        `json:"can_reply"`
	Title           string      `json:"title"`
	CanReshare      bool        `json:"can_reshare"`
	ReelType        string      `json:"reel_type"`
	User            User        `json:"user"`
	Items           []Item      `json:"items"`
	ReelMentions    []string    `json:"reel_mentions"`
	PrefetchCount   int         `json:"prefetch_count"`
	// this field can be int or bool
	HasBestiesMedia      interface{} `json:"has_besties_media"`
	StoryRankingToken    string      `json:"story_ranking_token"`
	Broadcasts           []Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
	HasNewNuxStory       bool        `json:"has_new_nux_story"`
	Status               string      `json:"status"`
}
