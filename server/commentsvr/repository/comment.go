package repository

import (
	"commentsvr/log"
	"strconv"
)

func CommentAdd(userId, videoId int64, comment_text string) (*Comment, error) {

	// 1. 更新数据库
	comment, err := DbCommentAdd(userId, videoId, comment_text)
	if err != nil {
		log.Errorf("CommentAdd|DbCommentAdd err:%v", err)
		return nil, err
	}

	// 2. 更新缓存
	if err := CacheSetComment(comment); err != nil {
		log.Errorf("CommentAdd|CacheSetComment err:%v", err)
		return nil, err
	}

	return comment, nil
}

func CommentDelete(videoId, commentID int64) error {
	if err := DbCommentDelete(commentID); err != nil {
		log.Errorf("CommentDelete|DbCommentDelete err:%v", err)
		return err
	}
	if err := CacheDelComment([]string{strconv.FormatInt(commentID, 10)}, videoId); err != nil {
		log.Errorf("CommentDelete|CacheDelComment err:%v", err)
		return err
	}
	return nil
}

func CommentList(videoId int64) ([]*Comment, error) {

	// 1. 从缓存中获取评论列表
	comments, err := CacheGetCommentList(videoId)
	if err != nil {
		log.Errorf("CommentList|CacheGetCommentList err:%v", err)
		return nil, err
	}

	// 如果缓存中有，则直接返回
	if len(comments) != 0 {
		return comments, nil
	}

	// 2. 如果缓存中没有，则从数据库中获取评论列表
	comments, err = DbCommentList(videoId)
	if err != nil {
		log.Errorf("CommentList|DbCommentList err:%v", err)
		return nil, err
	}
	return comments, nil
}
