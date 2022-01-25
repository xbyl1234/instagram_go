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
	Inst *Instagram `bson:"-"`

	Q         string     `json:"q" bson:"q"`
	RankToken string     `json:"rank_token" bson:"rank_token"`
	PageToken string     `json:"page_token" bson:"page_token"`
	HasMore   bool       `json:"has_more" bson:"has_more"`
	Type      SearchType `json:"type" bson:"type"`
}

func newSearch(inst *Instagram, q string) *Search {
	search := &Search{
		Inst:    inst,
		Q:       q,
		HasMore: true,
	}
	return search
}

func (this *Search) SetAccount(inst *Instagram) {
	this.Inst = inst
}

type SearchResult struct {
	search *Search

	BaseApiResp
	HasMore   bool   `json:"has_more"`
	RankToken string `json:"rank_token"`
	PageToken string `json:"page_token"`

	Tags []Tags `json:"results"`

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

func (this *SearchResult) GetTags() []Tags {
	if this.Tags == nil {
		return nil
	}

	for index := range this.Tags {
		this.Tags[index].Inst = this.search.Inst
		this.Tags[index].MoreAvailable = true
		this.Tags[index].Session = "0_" + common.GenUUID()
	}
	return this.Tags
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
	if !this.HasMore {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}

	res := &SearchResult{}
	var params = map[string]interface{}{
		"search_surface":  "hashtag_search_page",
		"timezone_offset": -28800,
		"count":           30,
		"q":               this.Q,
		"is_typeahead":    true,
	}

	if this.PageToken != "" {
		//params["rank_token"] = this.RankToken
		params["page_token"] = this.PageToken
	}

	err := this.Inst.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSearchTag,
			Query:   params,
		}, res)

	if err == nil {
		this.RankToken = res.RankToken
		this.PageToken = res.PageToken
		res.search = this
	}
	this.HasMore = res.HasMore
	return res, err
}

func (this *Search) NextLocation() (*SearchResult, error) {
	this.Type = SearchType_Location
	if !this.HasMore {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}

	insta := this.Inst
	params := map[string]interface{}{
		"places_search_page": "places_search_page",
		"timezone_offset":    -28800,
		"query":              this.Q,
		"is_typeahead":       true,
	}
	if this.PageToken != "" {
		//params["rank_token"] = this.RankToken
		params["page_token"] = this.PageToken
	}

	res := &SearchResult{}
	err := insta.HttpRequestJson(
		&reqOptions{
			ApiPath: urlSearchLocation,
			Query:   params,
		}, res)

	if err == nil {
		this.RankToken = res.RankToken
		this.PageToken = res.PageToken
		res.search = this
	}
	this.HasMore = res.HasMore
	return res, err
}
