package goinsta

import (
	"fmt"
)

//// Hashtag is used for getting the media that matches a hashtag on instagram.
type Hashtag struct {
	inst      *Instagram
	rankToken string
	name      string

	BaseApiResp
	Sections []struct {
		LayoutType    string `json:"layout_type"`
		LayoutContent struct {
			FillItems []struct {
				Media Item `json:"media"`
			} `json:"fill_items"`
			TwoByTwoItem struct {
				Channel struct {
					Media       Item   `json:"media"`
					ChannelId   string `json:"channel_id"`
					ChannelType string `json:"channel_type"`
					Context     string `json:"context"`
					Header      string `json:"header"`
					MediaCount  int    `json:"media_count"`
					Title       string `json:"title"`
				} `json:"channel"`
			} `json:"two_by_two_item"`
			Medias []struct {
				Item Item `json:"media"`
			} `json:"medias"`
		} `json:"layout_content"`
		FeedType        string `json:"feed_type"`
		ExploreItemInfo struct {
			NumColumns      int     `json:"num_columns"`
			TotalNumColumns int     `json:"total_num_columns"`
			AspectRatio     float32 `json:"aspect_ratio"`
			Autoplay        bool    `json:"autoplay"`
		} `json:"explore_item_info"`
	} `json:"sections"`
	MediaCount          int     `json:"media_count"`
	MoreAvailable       bool    `json:"more_available"`
	NextID              string  `json:"next_max_id"`
	NextPage            int     `json:"next_page"`
	NextMediaIds        []int64 `json:"next_media_ids"`
	AutoLoadMoreEnabled bool    `json:"auto_load_more_enabled"`
}

func (this *Hashtag) setValues() {
	for i := range this.Sections {
		for j := range this.Sections[i].LayoutContent.Medias {
			m := &FeedMedia{
				inst: this.inst,
			}
			setToItem(&this.Sections[i].LayoutContent.Medias[j].Item, m)
		}
	}
}

// NewHashtag returns initialised hashtag structure
// Name parameter is hashtag name
func (inst *Instagram) NewHashtag(name string) *Hashtag {
	return &Hashtag{
		inst: inst,
		name: name,
	}
}

type RespTagsInfo struct {
	BaseApiResp
	AllowFollowing             int           `json:"allow_following"`
	AllowMutingStory           bool          `json:"allow_muting_story"`
	ChallengeId                interface{}   `json:"challenge_id"`
	DebugInfo                  interface{}   `json:"debug_info"`
	Description                string        `json:"description"`
	DestinationInfo            interface{}   `json:"destination_info"`
	FollowButtonText           string        `json:"follow_button_text"`
	FollowStatus               int           `json:"follow_status"`
	Following                  int           `json:"following"`
	FormattedMediaCount        string        `json:"formatted_media_count"`
	FreshTopicMetadata         interface{}   `json:"fresh_topic_metadata"`
	Id                         uint64        `json:"id"`
	MediaCount                 int           `json:"media_count"`
	Name                       string        `json:"name"`
	NonViolating               int           `json:"non_violating"`
	ProfilePicUrl              interface{}   `json:"profile_pic_url"`
	PromoBanner                interface{}   `json:"promo_banner"`
	RelatedTags                interface{}   `json:"related_tags"`
	ShowFollowDropDown         bool          `json:"show_follow_drop_down"`
	SocialContext              string        `json:"social_context"`
	SocialContextFacepileUsers []interface{} `json:"social_context_facepile_users"`
	SocialContextProfileLinks  []interface{} `json:"social_context_profile_links"`
	Subtitle                   string        `json:"subtitle"`
}

// Sync updates Hashtag information preparing it to Next call.
func (this *Hashtag) Sync() error {
	resp := &RespTagsInfo{}
	err := this.inst.HttpRequestJson(&reqOptions{
		Endpoint: fmt.Sprintf(urlTagSync, this.name),
	}, resp)

	return err
}

// Stories returns hashtag stories.
func (this *Hashtag) Stories() (*StoryMedia, error) {
	var resp struct {
		Story  StoryMedia `json:"story"`
		Status string     `json:"status"`
	}

	err := this.inst.HttpRequestJson(&reqOptions{
		Endpoint: fmt.Sprintf(urlTagStories, this.name),
	}, &resp)

	return nil, err
}

// Next paginates over hashtag pages (xd).
func (this *Hashtag) Next() (*Hashtag, error) {
	var params = map[string]interface{}{
		"_uuid":      this.inst.uuid,
		"rank_token": this.rankToken,
	}

	if this.NextID == "" {
		params["supported_tabs"] = []string{"top", "recent"}
		params["include_persistent"] = true
		params["rank_token"] = this.rankToken
	} else {
		params["max_id"] = this.NextID
		params["tab"] = "top"
		params["page"] = this.NextPage
		params["include_persistent"] = false
		params["next_media_ids"] = this.NextMediaIds
	}

	ht := &Hashtag{
		inst:      this.inst,
		name:      this.name,
		rankToken: this.rankToken,
	}
	err := this.inst.HttpRequestJson(
		&reqOptions{
			Query: map[string]interface{}{
				"max_id":     this.NextID,
				"rank_token": "",
				"page":       fmt.Sprintf("%d", this.NextPage),
			},
			Endpoint: fmt.Sprintf(urlTagContent, this.name),
			IsPost:   false,
		}, ht,
	)

	return ht, err
}
