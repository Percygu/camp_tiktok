package repository

type Relation struct {
	Id       int64 `gorm:"column:id; primary_key;"`
	Follow   int64 `gorm:"column:follow_id"`   // 博主
	Follower int64 `gorm:"column:follower_id"` // 粉絲
}

func (r *Relation) TableName() string {
	return "t_relation"
}
