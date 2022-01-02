package goinsta

import (
	"fmt"
	"makemoney/common"
	"strings"
)

type Comments struct {
	media   *Item
	Inst    *Instagram
	MediaID string `json:"media_id"`
	//Next    struct {
	//	CachedCommentsCursor string `json:"cached_comments_cursor"`
	//	BifilterToken        string `json:"bifilter_token"`
	//} `json:"next"`
	Next    string `json:"next"`
	HasMore bool   `json:"has_more"`
}

func (this *Comments) SetAccount(inst *Instagram) {
	this.Inst = inst
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

func (this *RespComments) GetAllComments() []Comment {
	return this.Comments
}

func (this *Comments) NextComments() (*RespComments, error) {
	if !this.HasMore {
		return nil, &common.MakeMoneyError{
			ErrType: common.NoMoreError,
		}
	}

	params := map[string]interface{}{
		"can_support_threading": true,
	}

	if this.Next != "" {
		params["min_id"] = this.Next
	}

	ret := &RespComments{}
	err := this.Inst.HttpRequestJson(&reqOptions{
		IsPost:  false,
		ApiPath: fmt.Sprintf(urlComment, this.MediaID),
		Query:   params,
	}, ret)

	err = ret.CheckError(err)
	if err == nil {
		this.HasMore = ret.HasMoreHeadloadComments
		if this.HasMore {
			this.Next = strings.ReplaceAll(ret.NextMinId, "\\", "\"")
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
