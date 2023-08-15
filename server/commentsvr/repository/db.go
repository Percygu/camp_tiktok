package repository

import (
	"commentsvr/middleware/db"
	"time"
)

func DbCommentAdd(userId, videoId int64, comment_text string) (*Comment, error) {
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

	return &comment, nil
}

func DbCommentDelete(commentID int64) error {
	db := db.GetDB()
	comment := Comment{}
	err := db.Model(&Comment{}).Where("id = ?", commentID).Delete(&comment).Error
	if err != nil {
		return err
	}
	return nil
}

func DbCommentList(videoId int64) ([]*Comment, error) {
	db := db.GetDB()
	var comments []*Comment
	err := db.Where("video_id = ?", videoId).Order("id DESC").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
