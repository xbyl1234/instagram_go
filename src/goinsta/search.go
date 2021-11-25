package goinsta

import (
	"makemoney/common"
)

type SearchType int

var (
	SearchType_User     SearchType
	SearchType_Tags     SearchType
	SearchType_Location SearchType
	SearchType_Top      SearchType
)

type Search struct {
	inst      *Instagram
	q         string
	rankToken string
	pageToken string
	hasMore   bool
	Type      SearchType
}

func newSearch(inst *Instagram, q string) *Search {
	search := &Search{
		inst:    inst,
		q:       q,
		hasMore: true,
	}
	return search
}

type SearchResult struct {
	search *Search

	BaseApiResp
	HasMore   bool   `json:"has_more"`
	RankToken string `json:"rank_token"`
	PageToken string `json:"page_token"`

	Tags []struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		MediaCount int    `json:"media_count"`
	} `json:"results"`

	NumResults int64 `json:"num_results"`
	// User search results
	Users []User `json:"users"`

	// Tag search results
	InformModule interface{} `json:"inform_module"`
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

func (this *SearchResult) GetTags() []*Tags {
	if this.Tags == nil {
		return nil
	}
	ret := make([]*Tags, len(this.Tags))
	for index := range this.Tags {
		ret[index] = newTags(this.Tags[index].Name, this.search.inst)
	}
	return ret
}

// User search by username, you can use count optional parameter to get more than 50 items.
//func (this *Search) User(user string, countParam ...int) (*SearchResult, error) {
//	count := 50
//	if len(countParam) > 0 {
//		count = countParam[0]
//	}
//	insta := this.inst
//	res := &SearchResult{}
//
//	err := insta.HttpRequestJson(
//		&reqOptions{
//			ApiPath: urlSearchUser,
//			Query: map[string]interface{}{
//				"ig_sig_key_version": goInstaSigKeyVersion,
//				"is_typeahead":       "true",
//				"q":                  user,
//				"count":              fmt.Sprintf("%d", count),
//				//"rank_token":         insta.rankToken,
//			}}, res)
//
//	if err != nil {
//		return nil, err
//	}
//
//	for id := range res.Users {
//		res.Users[id].inst = insta
//	}
//	return res, err
//}

func (this *Search) NextTags() (*SearchResult, error) {
	this.Type = SearchType_Tags
	if !this.hasMore {
		return nil, common.MakeMoneyError_NoMore
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

	if err == nil {
		this.rankToken = res.RankToken
		this.pageToken = res.PageToken
	}
	this.hasMore = res.HasMore
	return res, err
}

func (this *Search) NextLocation() (*SearchResult, error) {
	this.Type = SearchType_Location
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

	if err == nil {
		this.rankToken = res.RankToken
		this.pageToken = res.PageToken
	}
	this.hasMore = res.HasMore
	return res, err
}
