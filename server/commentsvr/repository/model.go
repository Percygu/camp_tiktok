package repository

type Comment struct {
	CommentId int64  `gorm:"column:comment_id; primary_key;"`
	UserId    int64  `gorm:"column:user_id"`
	VideoId   int64  `gorm:"column:video_id"`
	Comment   string `gorm:"column:comment"`
	Time      string `gorm:"column:time"`
}

func (c *Comment) TableName() string {
	return "t_comments"
}
