package repository

type User struct {
	// gorm.Model
	Id              int64  `gorm:"column:user_id; primary_key;"`
	Name            string `gorm:"column:user_name"`
	Password        string `gorm:"column:password"`
	Follow          int64  `gorm:"column:follow_count"`
	Follower        int64  `gorm:"column:follower_count"`
	Avatar          string `gorm:"column:avatar"`
	BackgroundImage string `gorm:"column:background_image"`
	Signature       string `gorm:"column:signature"`
	TotalFav        int64  `gorm:"column:total_favorited"`
	FavCount        int64  `gorm:"column:favorite_count"`
}
