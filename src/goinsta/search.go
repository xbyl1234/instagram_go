package goinsta

import (
	"fmt"
	"makemoney/common"
)

// Search is the object for all searches like Facebook, Location or Tag search.
type Search struct {
	inst      *Instagram
	q         string
	rankToken string
	pageToken string
	hasMore   bool
}

// SearchResult handles the data for the results given by each type of Search.
type SearchResult struct {
	BaseApiResp
	HasMore    bool   `json:"has_more"`
	RankToken  string `json:"rank_token"`
	PageToken  string `json:"page_token"`
	Status     string `json:"status"`
	NumResults int64  `json:"num_results"`
	// User search results
	Users []User `json:"users"`

	// Tag search results
	InformModule interface{} `json:"inform_module"`
	Tags         []struct {
		ID               int64       `json:"id"`
		Name             string      `json:"name"`
		MediaCount       int         `json:"media_count"`
		FollowStatus     interface{} `json:"follow_status"`
		Following        interface{} `json:"following"`
		AllowFollowing   interface{} `json:"allow_following"`
		AllowMutingStory interface{} `json:"allow_muting_story"`
		ProfilePicURL    interface{} `json:"profile_pic_url"`
		NonViolating     interface{} `json:"non_violating"`
		RelatedTags      interface{} `json:"related_tags"`
		DebugInfo        interface{} `json:"debug_info"`
	} `json:"results"`

	// Location search result
	RequestID string `json:"request_id"`
	Venues    []struct {
		ExternalIDSource string  `json:"external_id_source"`
		ExternalID       string  `json:"external_id"`
		Lat              float64 `json:"lat"`
		Lng              float64 `json:"lng"`
		Address          string  `json:"address"`
		Name             string  `json:"name"`
	} `json:"venues"`

	// Facebook
	// Facebook also uses `Users`
	Places   []interface{} `json:"places"`
	Hashtags []struct {
		Position int `json:"position"`
		Hashtag  struct {
			Name       string `json:"name"`
			ID         int64  `json:"id"`
			MediaCount int    `json:"media_count"`
		} `json:"hashtag"`
	} `json:"hashtags"`

	ClearClientCache bool `json:"clear_client_cache"`
}

// newSearch creates new Search structure
func newSearch(inst *Instagram, q string) *Search {
	search := &Search{
		inst:    inst,
		q:       q,
		hasMore: true,
	}
	return search
}

// User search by username, you can use count optional parameter to get more than 50 items.
func (this *Search) User(user string, countParam ...int) (*SearchResult, error) {
	count := 50
	if len(countParam) > 0 {
		count = countParam[0]
	}
	insta := this.inst
	res := &SearchResult{}

	err := insta.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSearchUser,
			Query: map[string]interface{}{
				"ig_sig_key_version": goInstaSigKeyVersion,
				"is_typeahead":       "true",
				"q":                  user,
				"count":              fmt.Sprintf("%d", count),
				//"rank_token":         insta.rankToken,
			}}, res)

	if err != nil {
		return nil, err
	}

	for id := range res.Users {
		res.Users[id].inst = insta
	}
	return res, err
}

// Tags search by tag
func (this *Search) NextTags() (*SearchResult, error) {
	if !this.hasMore {
		return nil, &common.MakeMoneyError{"no more", 0}
	}

	res := &SearchResult{}
	var params = map[string]interface{}{
		"search_surface":  "hashtag_search_page",
		"timezone_offset": -18000,
		"count":           30,
		"q":               this.q,
	}

	if this.pageToken != "" {
		params["rank_token"] = this.rankToken
		params["page_token"] = this.pageToken
	}

	err := this.inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSearchTag,
			Query:   params,
		}, res)

	if err != nil {
		this.rankToken = res.RankToken
		this.pageToken = res.PageToken
	}
	this.hasMore = res.HasMore
	return res, err
}

func (this *Search) NextLocation() (*SearchResult, error) {
	if !this.hasMore {
		return nil, &common.MakeMoneyError{"no more", 0}
	}

	insta := this.inst
	params := map[string]interface{}{
		"places_search_page": "places_search_page",
		"timezone_offset":    -18000,
		"lat":                nil,
		"lng":                nil,
		"count":              30,
		"query":              this.q,
	}
	if this.pageToken != "" {
		params["rank_token"] = this.rankToken
		params["page_token"] = this.pageToken
	}

	res := &SearchResult{}
	err := insta.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSearchLocation,
			Query:   params,
		}, res)

	if err != nil {
		this.rankToken = res.RankToken
		this.pageToken = res.PageToken
	}
	this.hasMore = res.HasMore
	return res, err
}

//// Facebook search by facebook user.
//func (this *Search) Facebook(user string) (*SearchResult, error) {
//	insta := this.inst
//
//	res := &SearchResult{}
//	err := insta.HttpRequestJson(
//		&reqOptions{
//			ApiPath: urlSearchFacebook,
//			Query: map[string]interface{}{
//				"query":      user,
//				"rank_token": insta.rankToken,
//			},
//		}, res)
//	if err != nil {
//		return nil, err
//	}
//	return res, err
//}
