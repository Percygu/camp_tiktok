package repository

import (
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
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Find(&Favorite{}).Error
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("you have liked this video")
	}
	err = db.Create(&favorite).Error

	if err != nil {
		return err
	}

	// TODO: update user TotalFav
	// 得这个video的作者id

	// authorid, _ := CacheGetAuthor(vid) // todo videosvr
	// todo usercountcache change usersvr
	// go CacheChangeUserCount(uid, add, "like")
	// go CacheChangeUserCount(authorid, add, "liked")
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
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Find(&Favorite{}).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
