package repository

import (
	"favoritesvr/log"
	db "favoritesvr/middleware/db"
	"fmt"
	"gorm.io/gorm"
)

// Favorite 点赞
func LikeAction(uid, vid int64) error {
	db := db.GetDB()
	favorite := Favorite{
		UserId:  uid,
		VideoId: vid,
	}
	// 这个video 是否有被当前userid 点赞过
	err := db.Where("user_id = ? and video_id = ?", uid, vid).First(&Favorite{}).Error
	// 沒有err，说明有这条记录，说明已经点赞过了
	if err == nil {
		return fmt.Errorf("you have Like this video")
	}
	// 有err，但不是gorm.ErrRecordNotFound，说明有其他错误，返回
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return err
	}
	// 沒有這條記錄，則插入
	err = db.Create(&favorite).Error

	if err != nil {
		return err
	}

	return nil
}

// UnFavorite 取消点赞
func UnLikeAction(uid, vid int64) error {
	db := db.GetDB()
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Delete(&Favorite{}).Error
	if err != nil {
		return err
	}
	// authorid, _ := CacheGetAuthor(vid)
	// go func() {
	// go CacheChangeUserCount(uid, sub, "like")
	// go CacheChangeUserCount(authorid, sub, "liked")
	// }()
	return nil
}

// GetFavoriteList 获取我点赞的视频
func GetFavoriteIdList(uid int64) ([]int64, error) {
	db := db.GetDB()
	var favoriteList []*Favorite
	err := db.Model(&Favorite{}).Where("user_id= ?", uid).Find(&favoriteList).Error

	if err == gorm.ErrRecordNotFound {
		return []int64{}, nil
	} else if err != nil {
		return nil, err
	}

	var videoIDList []int64
	for _, favorite := range favoriteList {
		videoIDList = append(videoIDList, favorite.VideoId)
	}

	return videoIDList, nil
}

// IsFavoriteVideo 判断是否点赞这个视频
func IsFavoriteVideo(uid, vid int64) (bool, error) {
	db := db.GetDB()
	err := db.Where("user_id = ? and video_id = ?", uid, vid).First(&Favorite{}).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			log.Info("111")
			return false, nil
		}
		log.Info("222")
		return false, err
	}

	log.Infof("is favorite video", uid, vid)
	return true, nil
}
