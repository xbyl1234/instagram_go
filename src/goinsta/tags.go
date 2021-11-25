package goinsta

import (
	"fmt"
	"makemoney/common"
)

// Tags is used for getting the media that matches a hashtag on instagram.
type Tags struct {
	inst      *Instagram
	name      string
	rankToken string

	moreAvailable bool
	nextID        string
	nextPage      int
	nextMediaIds  []int64
}

func newTags(name string, inst *Instagram) *Tags {
	return &Tags{
		inst:      inst,
		name:      name,
		rankToken: common.GenUUID(),
	}
}

type RespHashtag struct {
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

func (this *RespHashtag) GetAllMedias() []*Item {
	var allCount int = 0
	for sectionIndex := range this.Sections {
		allCount += this.Sections[sectionIndex].ExploreItemInfo.TotalNumColumns
	}
	ret := make([]*Item, allCount)

	var index = 0
	for sectionIndex := range this.Sections {
		allCount += this.Sections[sectionIndex].ExploreItemInfo.TotalNumColumns
		section := this.Sections[sectionIndex]
		if section.LayoutType == "two_by_two_right" {
			ret[index] = &section.LayoutContent.TwoByTwoItem.Channel.Media
			index++
			for itemIndex := range section.LayoutContent.FillItems {
				ret[index] = &section.LayoutContent.FillItems[itemIndex].Media
				index++
			}
		} else if section.LayoutType == "media_grid" {
			for itemIndex := range section.LayoutContent.Medias {
				ret[index] = &section.LayoutContent.Medias[itemIndex].Item
				index++
			}
		}
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

// Sync updates Tags information preparing it to Next call.
func (this *Tags) Sync() error {
	resp := &RespTagsInfo{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: fmt.Sprintf(urlTagSync, this.name),
	}, resp)

	return err
}

// Stories returns hashtag stories.
func (this *Tags) Stories() (*StoryMedia, error) {
	var resp struct {
		Story  StoryMedia `json:"story"`
		Status string     `json:"status"`
	}

	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: fmt.Sprintf(urlTagStories, this.name),
	}, &resp)

	return &resp.Story, err
}

// Next paginates over hashtag pages (xd).
func (this *Tags) Next() (*RespHashtag, error) {
	if !this.moreAvailable {
		return nil, common.MakeMoneyError_NoMore
	}

	var params = map[string]interface{}{
		"_uuid":      this.inst.uuid,
		"rank_token": this.rankToken,
	}

	if this.nextID == "" {
		params["supported_tabs"] = []string{"top", "recent"}
		params["include_persistent"] = true
		params["rank_token"] = this.rankToken
	} else {
		params["max_id"] = this.nextID
		params["tab"] = "top"
		params["page"] = this.nextPage
		params["include_persistent"] = false
		params["next_media_ids"] = this.nextMediaIds
	}

	ht := &RespHashtag{}
	err := this.inst.HttpRequestJson(
		&reqOptions{
			Query: map[string]interface{}{
				"max_id":     this.nextID,
				"rank_token": "",
				"page":       fmt.Sprintf("%d", this.nextPage),
			},
			ApiPath: fmt.Sprintf(urlTagContent, this.name),
			IsPost:  false,
		}, ht,
	)

	err = ht.CheckError(err)
	if err == nil {
		this.nextID = ht.NextID
		this.nextPage = ht.NextPage
		this.nextMediaIds = ht.NextMediaIds
		this.moreAvailable = ht.MoreAvailable
	}

	return ht, err
}
