package repository

import (
	"commentsvr/log"
	"commentsvr/middleware/db"
	"strconv"
	"time"
)

func CommentAdd(userId, videoId int64, comment_text string) (*Comment, error) {
	db := db.GetDB()

	nowTime := time.Now()
	comment := Comment{
		UserId:      userId,
		VideoId:     videoId,
		CommentText: comment_text,
		CreateTime:  nowTime,
	}
	result := db.Create(&comment)

	if result.Error != nil {
		return nil, result.Error
	}
	// 评论缓存起来
	if err := SetCommentCacheInfo(&comment); err != nil {
		log.Errorf("CommentAdd|SetCommentCacheInfo err:%v", err)
	}

	return &comment, nil
}

func CommentDelete(videoId, commentID int64) error {
	db := db.GetDB()
	comment := Comment{}

	err := db.Model(&Comment{}).Where("id = ?", commentID).Take(&comment).Error
	if err != nil {
		return err
	}
	commentIDStr := strconv.FormatInt(commentID, 10)
	DelCommentCacheInfo([]string{commentIDStr}, videoId)
	db.Delete(&comment)
	return nil
}

func CommentList(videoId int64) ([]*Comment, error) {
	var comments []*Comment
	db := db.GetDB()
	var err error
	// comments, err = GetCommentCacheList(videoId)
	log.Infof("comments-------------------------:%+v\n", comments)

	if len(comments) != 0 {
		return comments, nil
	}

	err = db.Where("video_id = ?", videoId).Order("id DESC").Find(&comments).Error
	if err != nil {
		log.Errorf("get video with %d comment list err:%v", videoId, err)
		return nil, err
	}
	//
	// for _, comment := range comments {
	// 	if err := SetCommentCacheInfo(comment); err != nil {
	// 		log.Errorf("CommentAdd|SetCommentCacheInfo err:%v", err)
	// 		DelCacheCommentAll(videoId)
	// 		return comments, nil
	// 	}
	// }
	// log.Infof("comments:%+v", comments)

	return comments, nil
}
