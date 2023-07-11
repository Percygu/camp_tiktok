package repository

import (
	"fmt"
	"gorm.io/gorm"

	"relationsvr/log"
	"relationsvr/middleware/db"
)

func FollowAction(selfUserId, toUserId int64) error {
	db := db.GetDB()
	relation := Relation{
		Follow:   toUserId,
		Follower: selfUserId,
	}
	err := db.Where("follow_id = ? and follower_id = ?", toUserId, selfUserId).First(&Relation{}).Error
	// 沒有err，说明有这条记录，说明已经关注过了
	if err == nil {
		return fmt.Errorf("you have followed this user")
	}
	// 有err，但不是gorm.ErrRecordNotFound，说明有其他错误，返回
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return err
	}
	// 沒有這條記錄，則插入
	err = db.Create(&relation).Error
	if err != nil {
		return err
	}
	// 更新缓存中对应用户follower的数量
	// go CacheChangeUserCount(userId, add, "follow")
	// go CacheChangeUserCount(toUserId, add, "follower")
	return nil
}

func UnFollowAction(selfUserId, toUserId int64) error {
	db := db.GetDB()
	err := db.Where("follow_id = ? and follower_id = ?", toUserId, selfUserId).Delete(&Relation{}).Error
	if err != nil {
		return err
	}
	log.Debug("unfollow update user cache")
	// go CacheChangeUserCount(userId, sub, "follow")
	// go CacheChangeUserCount(toUserId, sub, "follower")
	return nil
}

// GetFollowList 获取我关注的博主
func GetFollowList(userId int64) ([]*Relation, error) {
	db := db.GetDB()
	relationList := make([]*Relation, 0)
	err := db.Where("follower_id = ?", userId).Find(&relationList).Error
	if err == gorm.ErrRecordNotFound {
		return relationList, nil
	} else if err != nil {
		return nil, err
	}
	return relationList, nil
}

// GetFollowList 获取关注我的粉丝
func GetFollowerList(userId int64) ([]*Relation, error) {
	db := db.GetDB()
	relationList := make([]*Relation, 0)
	err := db.Where("follow_id = ?", userId).Find(&relationList).Error
	if err == gorm.ErrRecordNotFound {
		return relationList, nil
	} else if err != nil {
		return nil, err
	}
	return relationList, nil
}

func IsFollow(selfUserId, toUserId int64) (bool, error) {
	if selfUserId == toUserId {
		return true, nil
	}

	db := db.GetDB()
	err := db.Where("follow_id = ? and follower_id = ?", toUserId, selfUserId).First(&Relation{}).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
