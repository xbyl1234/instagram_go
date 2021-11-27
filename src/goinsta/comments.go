package goinsta

import (
	"encoding/json"
	"fmt"
	"makemoney/common"
	"strings"
)

type Comments struct {
	item *Item
	next struct {
		CachedCommentsCursor string `json:"cached_comments_cursor"`
		BifilterToken        string `json:"bifilter_token"`
	}
	hasMore bool
}

type RespComments struct {
	BaseApiResp
	ScrollBehavior             int       `json:"scroll_behavior"`
	ThreadingEnabled           bool      `json:"threading_enabled"`
	HasMoreComments            bool      `json:"has_more_comments"`
	HasMoreHeadloadComments    bool      `json:"has_more_headload_comments"`
	InitiateAtTop              bool      `json:"initiate_at_top"`
	InsertNewCommentToTop      bool      `json:"insert_new_comment_to_top"`
	IsRanked                   bool      `json:"is_ranked"`
	MediaHeaderDisplay         string    `json:"media_header_display"`
	NextMinId                  string    `json:"next_min_id"`
	CaptionIsEdited            bool      `json:"caption_is_edited"`
	CommentCount               int       `json:"comment_count"`
	CommentCoverPos            string    `json:"comment_cover_pos"`
	CommentLikesEnabled        bool      `json:"comment_likes_enabled"`
	CanViewMorePreviewComments bool      `json:"can_view_more_preview_comments"`
	Comments                   []Comment `json:"comments"`
}

func (this *Comments) NextComments() (*RespComments, error) {
	if !this.hasMore {
		return nil, common.MakeMoneyError_NoMore
	}

	params := map[string]interface{}{
		"inventory_source":        "media_or_ad",
		"analytics_module":        "comments_v2_feed_timeline",
		"can_support_threading":   true,
		"is_carousel_bumped_post": true,
		"feed_position":           0,
	}

	if this.next.CachedCommentsCursor != "" {
		minId, _ := json.Marshal(this.next)
		params["min_id"] = minId
	}

	ret := &RespComments{}
	err := this.item.inst.HttpRequestJson(&reqOptions{
		IsPost:  false,
		ApiPath: fmt.Sprintf(urlComment, this.item.ID),
		Query:   params,
	}, ret)

	err = ret.CheckError(err)
	if err == nil {
		this.hasMore = ret.HasMoreComments
		if this.hasMore {
			next := strings.ReplaceAll(ret.NextMinId, "\\", "\"")
			err = json.Unmarshal([]byte(next), &this.next)
		}
	}
	return ret, err
}

type Comment struct {
	BaseApiResp
	BitFlags                       int           `json:"bit_flags"`
	ChildCommentCount              int           `json:"child_comment_count"`
	CommentIndex                   int           `json:"comment_index"`
	CommentLikeCount               int           `json:"comment_like_count"`
	ContentType                    string        `json:"content_type"`
	CreatedAt                      int           `json:"created_at"`
	CreatedAtUtc                   int           `json:"created_at_utc"`
	DidReportAsSpam                bool          `json:"did_report_as_spam"`
	HasLikedComment                bool          `json:"has_liked_comment"`
	HasMoreHeadChildComments       bool          `json:"has_more_head_child_comments"`
	HasMoreTailChildComments       bool          `json:"has_more_tail_child_comments"`
	InlineComposerDisplayCondition string        `json:"inline_composer_display_condition"`
	IsCovered                      bool          `json:"is_covered"`
	NumTailChildComments           int           `json:"num_tail_child_comments"`
	Pk                             int64         `json:"pk"`
	PreviewChildComments           []interface{} `json:"preview_child_comments"`
	PrivateReplyStatus             int           `json:"private_reply_status"`
	ShareEnabled                   bool          `json:"share_enabled"`
	Text                           string        `json:"text"`
	Type                           int           `json:"type"`
	UserId                         int           `json:"user_id"`
	User                           User          `json:"user"`
}
