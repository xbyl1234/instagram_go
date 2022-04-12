package goinsta

import (
	"fmt"
	"makemoney/common"
)

type UserOperate struct {
	inst *Instagram
	//for post
	cameraEntryPoint int
	//for comment
	commentSessionId string
}

func newUserOperate(inst *Instagram) *UserOperate {
	return &UserOperate{
		inst:             inst,
		cameraEntryPoint: 0,
		commentSessionId: common.GenString(common.CharSet_16_Num, 32),
	}
}

type RespLikeUser struct {
	BaseApiResp
	PreviousFollowing bool       `json:"previous_following,omitempty"`
	FriendshipStatus  Friendship `json:"friendship_status"`
}

func (this *UserOperate) LikeUser(userID int64) error {
	this.inst.Increase(OperNameLikeUser)
	params := map[string]interface{}{
		"_uuid":            this.inst.AccountInfo.Device.DeviceID,
		"_uid":             this.inst.ID,
		"user_id":          userID,
		"device_id":        this.inst.AccountInfo.Device.DeviceID,
		"container_module": "profile",
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath:        fmt.Sprintf(urlUserFollow, userID),
		HeaderSequence: LoginHeaderMap[urlUserFollow],
		IsPost:         true,
		Signed:         true,
		Query:          params,
	}, resp)

	err = resp.CheckError(err)
	if err == nil {
		this.inst.IncreaseSuccess(OperNameLikeUser)
	}
	return err
}

type AddCommentParams struct {
	ParentCommentId string
	UserName        string

	LoggingInfoToken string
	MediaId          string
	CommentText      string
}

var AddCommentNavChain = []string{"IGSundialHomeViewController:clips_viewer_clips_tab:2,IGCommentThreadV2ViewController:comments_v2_clips_viewer_clips_tab:3"}
var AddSubCommentNavChain = []string{"IGSundialHomeViewController:clips_viewer_clips_tab:2"}

type RespCheckOffensiveComment struct {
	BaseApiResp
	IsOffensive     bool    `json:"is_offensive"`
	TextLanguage    string  `json:"text_language"`
	BullyClassifier float64 `json:"bully_classifier"`
}

type RespAddComment struct {
	BaseApiResp
	Comment struct {
		ContentType  string             `json:"content_type"`
		User         AddCommentRespUser `json:"user"`
		Pk           int64              `json:"pk"`
		Text         string             `json:"text"`
		Type         int                `json:"type"`
		CreatedAt    int                `json:"created_at"`
		CreatedAtUtc int                `json:"created_at_utc"`
		MediaId      int64              `json:"media_id"`
		Status       string             `json:"status"`
		ShareEnabled bool               `json:"share_enabled"`
	} `json:"comment"`
}

func (this *UserOperate) AddComment(params *AddCommentParams) error {
	this.inst.Increase(OperNameComment)

	checkParams := map[string]interface{}{
		"media_id":           params.MediaId,
		"_uuid":              this.inst.AccountInfo.Device,
		"_uid":               fmt.Sprintf("%d", this.inst.ID),
		"comment_session_id": this.commentSessionId,
	}

	commentParam := map[string]interface{}{
		"delivery_class":          "organic",
		"logging_info_token":      params.LoggingInfoToken,
		"_uuid":                   this.inst.AccountInfo.Device,
		"_uid":                    fmt.Sprintf("%d", this.inst.ID),
		"idempotence_token":       common.GenString(common.CharSet_16_Num, 32),
		"is_carousel_bumped_post": "false",
		"carousel_index":          "0",
		"container_module":        "comments_v2_clips_viewer_clips_tab",
	}
	if params.ParentCommentId != "" {
		text := "@" + params.UserName + " " + params.CommentText
		checkParams["comment_text"] = text
		commentParam["comment_text"] = text
		commentParam["parent_comment_id"] = params.ParentCommentId
		commentParam["replied_to_comment_id"] = params.ParentCommentId
		commentParam["nav_chain"] = AddSubCommentNavChain[common.GenNumber(0, len(AddSubCommentNavChain))]
	} else {
		checkParams["comment_text"] = params.CommentText
		commentParam["comment_text"] = params.CommentText
		commentParam["nav_chain"] = AddCommentNavChain[common.GenNumber(0, len(AddCommentNavChain))]
	}
	checkResp := &RespCheckOffensiveComment{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:         true,
		Signed:         true,
		ApiPath:        urlCheckOffensiveComment,
		HeaderSequence: LoginHeaderMap[urlCheckOffensiveComment],
		Query:          checkParams,
	}, checkResp)

	err = checkResp.CheckError(err)
	if err != nil {
		return err
	}

	addResp := &RespAddComment{}
	err = this.inst.HttpRequestJson(&reqOptions{
		IsPost:         true,
		Signed:         true,
		ApiPath:        fmt.Sprintf(urlCommentAdd, params.MediaId),
		HeaderSequence: LoginHeaderMap[urlCommentAdd],
		Query:          commentParam,
	}, addResp)

	err = addResp.CheckError(err)
	if err == nil {
		this.inst.IncreaseSuccess(OperNameComment)
	}
	return err
}

type RespShareMedia struct {
	BaseApiResp
	Permalink string `json:"permalink"`
}

func (this *UserOperate) ShareMedia(mediaId string) (string, error) {
	resp := &RespShareMedia{}
	err := this.inst.HttpRequestJson(&reqOptions{
		IsPost:         false,
		ApiPath:        fmt.Sprintf(urlShareMedia, mediaId),
		HeaderSequence: LoginHeaderMap[urlShareMedia],
		Query: map[string]interface{}{
			"share_to_app": "copy_link",
		},
	}, resp)
	err = resp.CheckError(err)
	return resp.Permalink, err
}
