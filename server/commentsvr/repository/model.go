package repository

type Comment struct {
	Id          int64  `gorm:"column:id; primary_key;"`
	UserId      int64  `gorm:"column:user_id"`
	VideoId     int64  `gorm:"column:video_id"`
	CommentText string `gorm:"column:comment_text"`
	CreateTime  string `gorm:"column:create_time"`
}

func (c *Comment) TableName() string {
	return "t_comments"
}
