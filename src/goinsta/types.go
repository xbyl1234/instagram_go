package goinsta

import "encoding/json"

//
//import (
//	"encoding/json"
//)
//
// School is void structure (yet).
type School struct {
}

// PicURLInfo repre
type PicURLInfo struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// Nametag is part of the account information.
type Nametag struct {
	Mode          int64       `json:"mode"`
	Gradient      json.Number `json:"gradient,Number"`
	Emoji         string      `json:"emoji"`
	SelfieSticker json.Number `json:"selfie_sticker,Number"`
}

//
//type friendResp struct {
//	Status     string     `json:"status"`
//	Friendship Friendship `json:"friendship_status"`
//}
//
// Location stores media location information.
type Location struct {
	Pk               int64   `json:"pk"`
	Name             string  `json:"name"`
	Address          string  `json:"address"`
	City             string  `json:"city"`
	ShortName        string  `json:"short_name"`
	Lng              float64 `json:"lng"`
	Lat              float64 `json:"lat"`
	ExternalSource   string  `json:"external_source"`
	FacebookPlacesID int64   `json:"facebook_places_id"`
}

//// SuggestedUsers stores the information about user suggestions.
type SuggestedUsers struct {
	Type        int `json:"type"`
	Suggestions []struct {
		User            User          `json:"user"`
		Algorithm       string        `json:"algorithm"`
		SocialContext   string        `json:"social_context"`
		Icon            string        `json:"icon"`
		Caption         string        `json:"caption"`
		MediaIds        []interface{} `json:"media_ids"`
		ThumbnailUrls   []interface{} `json:"thumbnail_urls"`
		LargeUrls       []interface{} `json:"large_urls"`
		MediaInfos      []interface{} `json:"media_infos"`
		Value           float64       `json:"value"`
		IsNewSuggestion bool          `json:"is_new_suggestion"`
	} `json:"suggestions"`
	LandingSiteType  string `json:"landing_site_type"`
	Title            string `json:"title"`
	ViewAllText      string `json:"view_all_text"`
	LandingSiteTitle string `json:"landing_site_title"`
	NetegoType       string `json:"netego_type"`
	UpsellFbPos      string `json:"upsell_fb_pos"`
	AutoDvance       string `json:"auto_dvance"`
	ID               string `json:"id"`
	TrackingToken    string `json:"tracking_token"`
}

//
// Friendship stores the details of the relationship between two users.
type Friendship struct {
	IncomingRequest bool `json:"incoming_request"`
	FollowedBy      bool `json:"followed_by"`
	OutgoingRequest bool `json:"outgoing_request"`
	Following       bool `json:"following"`
	Blocking        bool `json:"blocking"`
	IsPrivate       bool `json:"is_private"`
	Muting          bool `json:"muting"`
	IsMutingReel    bool `json:"is_muting_reel"`
	IsRestricted    bool `json:"is_restricted"`
	IsFeedFavorite  bool `json:"is_feed_favorite"`
	IsBestie        bool `json:"is_bestie"`
}

// Images are different quality images
type Images struct {
	Versions []Candidate `json:"candidates"`
}

//// GetBest returns the URL of the image with the best quality.
//func (img Images) GetBest() string {
//	best := ""
//	var mh, mw int
//	for _, v := range img.Versions {
//		if v.Width > mw || v.Height > mh {
//			best = v.URL
//			mh, mw = v.Height, v.Width
//		}
//	}
//	return best
//}
//
// Candidate is something that I really have no idea what it is.
type Candidate struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

// Caption is media caption
type Caption struct {
	ID                 int64  `json:"pk"`
	UserID             int64  `json:"user_id"`
	Text               string `json:"text"`
	Type               int    `json:"type"`
	CreatedAt          int64  `json:"created_at"`
	CreatedAtUtc       int64  `json:"created_at_utc"`
	ContentType        string `json:"content_type"`
	Status             string `json:"status"`
	BitFlags           int    `json:"bit_flags"`
	User               User   `json:"user"`
	DidReportAsSpam    bool   `json:"did_report_as_spam"`
	MediaID            int64  `json:"media_id"`
	HasTranslation     bool   `json:"has_translation"`
	IsCovered          bool   `json:"is_covered"`
	PrivateReplyStatus int    `json:"private_reply_status"`
	ShareEnabled       bool   `json:"share_enabled"`
}

// Mentions is a user being mentioned on media.
type Mentions struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int64   `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	User     User    `json:"user"`
}

// Video are different quality videos
type Video struct {
	Type   int    `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
	ID     string `json:"id"`
}

//type timeStoryResp struct {
//	Status string       `json:"status"`
//	Media  []StoryMedia `json:"tray"`
//}
//
//type trayResp struct {
//	Reels  map[string]StoryMedia `json:"reels"`
//	Status string                `json:"status"`
//}
//
//// Tray is a set of story media received from timeline calls.
//type Tray struct {
//	Stories []StoryMedia `json:"tray"`
//	Lives   struct {
//		LiveItems []LiveItems `json:"post_live_items"`
//	} `json:"post_live"`
//	StoryRankingToken    string      `json:"story_ranking_token"`
//	Broadcasts           []Broadcast `json:"broadcasts"`
//	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
//	HasNewNuxStory       bool        `json:"has_new_nux_story"`
//	Status               string      `json:"status"`
//}
//
//func (tray *Tray) set(inst *Instagram, url string) {
//	for i := range tray.Stories {
//		tray.Stories[i].inst = inst
//		tray.Stories[i].endpoint = url
//		tray.Stories[i].setValues()
//	}
//	for i := range tray.Lives.LiveItems {
//		tray.Lives.LiveItems[i].User.inst = inst
//		for j := range tray.Lives.LiveItems[i].Broadcasts {
//			tray.Lives.LiveItems[i].Broadcasts[j].BroadcastOwner.inst = inst
//		}
//	}
//	for i := range tray.Broadcasts {
//		tray.Broadcasts[i].BroadcastOwner.inst = inst
//	}
//}
//
//// LiveItems are Live media items
//type LiveItems struct {
//	ID                  string      `json:"pk"`
//	User                User        `json:"user"`
//	Broadcasts          []Broadcast `json:"broadcasts"`
//	LastSeenBroadcastTs float64     `json:"last_seen_broadcast_ts"`
//	RankedPosition      int64       `json:"ranked_position"`
//	SeenRankedPosition  int64       `json:"seen_ranked_position"`
//	Muted               bool        `json:"muted"`
//	CanReply            bool        `json:"can_reply"`
//	CanReshare          bool        `json:"can_reshare"`
//}
//
//// Broadcast is live videos.
type Broadcast struct {
	ID                   int64  `json:"id"`
	BroadcastStatus      string `json:"broadcast_status"`
	DashManifest         string `json:"dash_manifest"`
	ExpireAt             int64  `json:"expire_at"`
	EncodingTag          string `json:"encoding_tag"`
	InternalOnly         bool   `json:"internal_only"`
	NumberOfQualities    int    `json:"number_of_qualities"`
	CoverFrameURL        string `json:"cover_frame_url"`
	BroadcastOwner       User   `json:"broadcast_owner"`
	PublishedTime        int64  `json:"published_time"`
	MediaID              string `json:"media_id"`
	BroadcastMessage     string `json:"broadcast_message"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	UserId               int64  `json:"user_id"`
}

// BlockedUser stores information about a used that has been blocked before.
type BlockedUser struct {
	// TODO: Convert to user
	UserID        int64  `json:"user_id"`
	Username      string `json:"username"`
	FullName      string `json:"full_name"`
	ProfilePicURL string `json:"profile_pic_url"`
	BlockAt       int64  `json:"block_at"`
}

//// Unblock unblocks blocked user.
//func (b *BlockedUser) Unblock() error {
//	u := User{ID: b.UserID}
//	return u.Unblock()
//}
//
//type blockedListResp struct {
//	BlockedList []BlockedUser `json:"blocked_list"`
//	PageSize    int           `json:"page_size"`
//	Status      string        `json:"status"`
//}
//
//// InboxItemMedia is inbox media item
//type InboxItemMedia struct {
//	ClientContext              string `json:"client_context"`
//	ExpiringMediaActionSummary struct {
//		Count     int    `json:"count"`
//		Timestamp int64  `json:"timestamp"`
//		Type      string `json:"type"`
//	} `json:"expiring_media_action_summary"`
//	ItemID     string `json:"item_id"`
//	ItemType   string `json:"item_type"`
//	RavenMedia struct {
//		MediaType int64 `json:"media_type"`
//	} `json:"raven_media"`
//	ReplyChainCount int           `json:"reply_chain_count"`
//	SeenUserIds     []interface{} `json:"seen_user_ids"`
//	Timestamp       int64         `json:"timestamp"`
//	UserID          int64         `json:"user_id"`
//	ViewMode        string        `json:"view_mode"`
//}
//
////InboxItemLike is the heart sent during a conversation.
//type InboxItemLike struct {
//	ItemID    string `json:"item_id"`
//	ItemType  string `json:"item_type"`
//	Timestamp int64  `json:"timestamp"`
//	UserID    int64  `json:"user_id"`
//}
//
type respLikers struct {
	Users     []User `json:"users"`
	UserCount int64  `json:"user_count"`
	Status    string `json:"status"`
}

//type threadResp struct {
//	Conversation Conversation `json:"thread"`
//	Status       string       `json:"status"`
//}
//
//type ErrChallengeProcess struct {
//	StepName string
//}
//
//func (ec ErrChallengeProcess) Error() string {
//	return ec.StepName
//}
