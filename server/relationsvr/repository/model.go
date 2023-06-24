package repository

type Relation struct {
	Id       int64 `gorm:"column:relation_id; primary_key;"`
	Follow   int64 `gorm:"column:follow_id"`
	Follower int64 `gorm:"column:follower_id"`
}

func (r *Relation) TableName() string {
	return "t_relations"
}
