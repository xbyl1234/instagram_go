package goinsta

import (
	"container/list"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
)

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
		inst:          inst,
		name:          name,
		moreAvailable: true,
		rankToken:     common.GenUUID(),
	}
}

type RespHashtag struct {
	inst *Instagram
	BaseApiResp
	Sections []struct {
		LayoutType    string `json:"layout_type"`
		LayoutContent struct {
			FillItems []struct {
				Media Item `json:"media"`
			} `json:"fill_items"`
			OneByTwoItem struct {
				Clips struct {
					Items []struct {
						Media Item `json:"media"`
					} `json:"items"`
					Id            string `json:"id"`
					Tag           string `json:"tag"`
					MaxId         string `json:"max_id"`
					MoreAvailable bool   `json:"more_available"`
					Design        string `json:"design"`
					Label         string `json:"label"`
				} `json:"clips"`
			} `json:"one_by_two_item"`
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
	buff := list.New()

	for sectionIndex := range this.Sections {
		section := this.Sections[sectionIndex]
		if section.LayoutType == "two_by_two_right" {
			buff.PushBack(&section.LayoutContent.TwoByTwoItem.Channel.Media)
			for itemIndex := range section.LayoutContent.FillItems {
				buff.PushBack(&section.LayoutContent.FillItems[itemIndex].Media)
			}
		} else if section.LayoutType == "media_grid" {
			for itemIndex := range section.LayoutContent.Medias {
				buff.PushBack(&section.LayoutContent.Medias[itemIndex].Item)
			}
		} else if section.LayoutType == "one_by_two_item" || section.LayoutType == "one_by_two_left" {
			for itemIndex := range section.LayoutContent.OneByTwoItem.Clips.Items {
				buff.PushBack(&section.LayoutContent.OneByTwoItem.Clips.Items[itemIndex].Media)
			}
			for itemIndex := range section.LayoutContent.FillItems {
				buff.PushBack(&section.LayoutContent.FillItems[itemIndex].Media)
			}
		} else {
			log.Error("unknow LayoutType: %s", section.LayoutType)
		}
	}

	ret := make([]*Item, buff.Len())
	var index = 0
	for item := buff.Front(); item != nil; item = item.Next() {
		ret[index] = item.Value.(*Item)
		ret[index].inst = this.inst
		index++
	}
	return ret
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
			ApiPath: fmt.Sprintf(urlTagSections, this.name),
			IsPost:  true,
		}, ht,
	)

	err = ht.CheckError(err)
	if err == nil {
		this.nextID = ht.NextID
		this.nextPage = ht.NextPage
		this.nextMediaIds = ht.NextMediaIds
		this.moreAvailable = ht.MoreAvailable
		ht.inst = this.inst
	}

	return ht, err
}
