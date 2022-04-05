package goinsta

import (
	"fmt"
	"makemoney/common"
	"strings"
)

type Comments struct {
	Inst *Instagram `bson:"-"`

	MediaID string `json:"media_id" bson:"media_id"`
	//Next    struct {
	//	CachedCommentsCursor string `json:"cached_comments_cursor"`
	//	BifilterToken        string `json:"bifilter_token"`
	//} `json:"next"`
	Next    string `json:"next" bson:"next"`
	HasMore bool   `json:"has_more" bson:"has_more"`
}

func newComments(inst *Instagram, mediaID string) *Comments {
	return &Comments{
		Inst:    inst,
		MediaID: mediaID,
		Next:    "",
		HasMore: true,
	}
}

func (this *Comments) SetAccount(inst *Instagram) {
	this.Inst = inst
}

func (this *Comments) NextComments() (*RespComments, error) {
	if !this.HasMore {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}
	this.Inst.Increase(OperNameCrawComment)

	params := map[string]interface{}{
		"can_support_threading": true,
	}

	if this.Next != "" {
		params["min_id"] = this.Next
	}

	ret := &RespComments{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		IsPost:         false,
		ApiPath:        fmt.Sprintf(urlComment, this.MediaID),
		HeaderSequence: LoginHeaderMap[urlComment],
		Query:          params,
	}, ret)

	err = ret.CheckError(err)
	if err == nil {
		this.HasMore = ret.HasMoreComments
		if this.HasMore {
			this.Next = strings.ReplaceAll(ret.NextMinId, "\\", "\"")
		}
	}
	return ret, err
}
