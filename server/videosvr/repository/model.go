package repository

type Video struct {
	Id            int64  `gorm:"column:id; primary_key;"` // video_id
	AuthorId      int64  `gorm:"column:author_id;"`       // 谁发布的
	PlayUrl       string `gorm:"column:play_url;"`        // videoURL
	CoverUrl      string `gorm:"column:cover_url;"`       // picURL
	FavoriteCount int64  `gorm:"column:favorite_count;"`  // 点赞数
	CommentCount  int64  `gorm:"column:comment_count;"`   // 评论数
	PublishTime   int64  `gorm:"column:publish_time;"`    // 发布时间
	Title         string `gorm:"column:title;"`           // 标题
	// Author        User   `gorm:"foreignkey:AuthorId"`           // 作者
}

func (v *Video) TableName() string {
	return "t_video"
}

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:id; primary_key;"` // ID = userID
	Name            string `gorm:"column:user_name"`        // Name = userName = nickName
	Password        string `gorm:"column:password"`         // Password = password
	Follow          int64  `gorm:"column:follow_count"`     // Follow = followCount 我关注的人数
	Follower        int64  `gorm:"column:follower_count"`   // Follower = followerCount 关注我的人数
	Avatar          string `gorm:"column:avatar"`           // Avatar = avatar 头像
	BackgroundImage string `gorm:"column:background_image"` // BackgroundImage = backgroundImage 背景图
	Signature       string `gorm:"column:signature"`        // Signature = signature 个性签名
	TotalFav        int64  `gorm:"column:total_favorited"`  // 收到的赞
	FavCount        int64  `gorm:"column:favorite_count"`   // 我喜欢的视频数
}
