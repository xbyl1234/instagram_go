package goinsta

import (
	"container/list"
	"encoding/json"
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
)

var TabRecent = "recent"
var TabTop = "top"

type Tags struct {
	Id                   int64  `json:"id"`
	Name                 string `json:"name"`
	MediaCount           int    `json:"media_count"`
	FormattedMediaCount  string `json:"formatted_media_count"`
	SearchResultSubtitle string `json:"search_result_subtitle"`
	ProfilePicUrl        string `json:"profile_pic_url"`
	UseDefaultAvatar     bool   `json:"use_default_avatar"`
}

type TagsFeed struct {
	Feed
	Name          string
	MediaCount    int
	Session       string
	Tab           string
	MoreAvailable bool
	NextID        string
	NextPage      int
	NextMediaIds  []int64
}

func NewTagsFeed(inst *Instagram, tagName string, tab string) *TagsFeed {
	t := &TagsFeed{
		Feed:          Feed{Inst: inst},
		Name:          tagName,
		Session:       "0_" + common.GenUUID(),
		Tab:           tab,
		MoreAvailable: true,
	}
	return t
}

type RespHashtag struct {
	inst *Instagram
	BaseApiResp
	Sections []struct {
		LayoutType    string `json:"layout_type"`
		LayoutContent struct {
			FillItems []struct {
				Media Media `json:"media"`
			} `json:"fill_items"`
			OneByTwoItem struct {
				Clips struct {
					Items []struct {
						Media Media `json:"media"`
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
					Media       Media  `json:"media"`
					ChannelId   string `json:"channel_id"`
					ChannelType string `json:"channel_type"`
					Context     string `json:"context"`
					Header      string `json:"header"`
					MediaCount  int    `json:"media_count"`
					Title       string `json:"title"`
				} `json:"channel"`
			} `json:"two_by_two_item"`
			Medias []struct {
				Item Media `json:"media"`
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

func (this *RespHashtag) GetAllMedias() []*Media {
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

	ret := make([]*Media, buff.Len())
	var index = 0
	for item := buff.Front(); item != nil; item = item.Next() {
		ret[index] = item.Value.(*Media)
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

func (this *TagsFeed) Sync(tab string) error {
	this.Tab = tab

	resp := &RespTagsInfo{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath:        fmt.Sprintf(urlTagSync, this.Name),
		HeaderSequence: LoginHeaderMap[urlTagSync],
	}, resp)
	err = resp.CheckError(err)
	return err
}

func (this *TagsFeed) Stories() (*StoryMedia, error) {
	var resp struct {
		BaseApiResp
		Story StoryMedia `json:"story"`
	}

	err := this.Inst.HttpRequestJson(&reqOptions{
		ApiPath:        fmt.Sprintf(urlTagStories, this.Name),
		HeaderSequence: LoginHeaderMap[urlTagStories],
	}, &resp)

	err = resp.CheckError(err)
	return &resp.Story, err
}

func (this *TagsFeed) Next() (*RespHashtag, error) {
	if !this.MoreAvailable {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}

	this.Inst.Increase(OperNameCrawMedia)
	var params = map[string]interface{}{
		"_uuid":              this.Inst.AccountInfo.Device.DeviceID,
		"include_persistent": 0,
		"supported_tabs":     "[\"recent\",\"top\",\"igtv\",\"places\",\"shopping\"]",
		"tab":                TabTop,
		"surface":            "grid",
		"seen_media_ids":     "",
		"session_id":         this.Session,
	}

	if this.NextID != "" {
		params["max_id"] = this.NextID
		params["page"] = this.NextPage
		tmp, _ := json.Marshal(this.NextMediaIds)
		params["next_media_ids"] = common.B2s(tmp)
	}

	ht := &RespHashtag{}
	err := this.Inst.HttpRequestJson(
		&reqOptions{
			Query:          params,
			ApiPath:        fmt.Sprintf(urlTagSections, this.Name),
			HeaderSequence: LoginHeaderMap[urlTagSections],
			IsPost:         true,
		}, ht,
	)

	err = ht.CheckError(err)
	if err == nil {
		this.Inst.IncreaseSuccess(OperNameCrawMedia)
		this.NextID = ht.NextID
		this.NextPage = ht.NextPage
		this.NextMediaIds = ht.NextMediaIds
		this.MoreAvailable = ht.MoreAvailable
		ht.inst = this.Inst
	}

	return ht, err
}

func (this *TagsFeed) NextPost() (*Media, error) {
	if this.media == nil || this.indexMedia >= len(this.media) {
		next, err := this.Next()
		if err != nil {
			return nil, err
		}
		this.media = next.GetAllMedias()
		this.indexMedia = 0
	}
	ret := this.media[this.indexMedia]
	this.indexMedia++
	return ret, nil
}
