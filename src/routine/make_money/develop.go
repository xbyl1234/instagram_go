package main

import (
	"makemoney/goinsta"
	"time"
)

type DevelopMeta struct {
	inst *goinsta.Instagram
	Feed goinsta.Feed

	feed      *goinsta.VideoFeed
	comments  *goinsta.Comments
	followSet map[int64]bool

	curVideoList    *goinsta.VideosFeedResp
	curComments     *goinsta.RespComments
	nextVideoIdx    int
	nextCommentIdx  int
	subCommentCount int

	addSubCommentFinish bool
	hadShareMedia       bool
	hadCheckMedia       bool

	lastFeedBackTime time.Time
	isRunning        bool
}

//func NewDevelopMeta(inst *goinsta.Instagram, feed goinsta.Feed) *DevelopMeta {
//	d := &DevelopMeta{}
//	d.inst = inst
//	d.Feed = feed
//	return d
//}
//
//type FuncPostDeal func(inst *goinsta.Instagram, post *goinsta.Media) error
//
//func (this *DevelopMeta) ForEachPost(deal FuncPostDeal) error {
//	for true {
//		post, err := this.Feed.NextPost()
//		if err != nil {
//			return err
//		}
//
//		err = deal(this.inst, post)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//type FuncCommentDeal func(inst *goinsta.Instagram, fromPost *goinsta.VideosFeedResp, comment *goinsta.RespComments) error
//
//func (this *DevelopMeta) ForEachComments(post *goinsta.Media, deal FuncCommentDeal) error {
//	if post.CommentingDisabledForViewer {
//		return nil
//	}
//	if post.CommentCount > 0 {
//		comment := this.inst.NewComments(post.Id)
//
//	}
//}
