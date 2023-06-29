package repository

import (
	"context"
	db "favoritesvr/middleware/db"
	"favoritesvr/utils"
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"gorm.io/gorm"
)

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

func GetFavoriteList(uid int64) ([]*pb.VideoInfo, error) {
	db := db.GetDB()
	var favoriteList []*Favorite
	err := db.Model(&Favorite{}).Where("user_id= ?", uid).Find(&favoriteList).Error
	if err == gorm.ErrRecordNotFound {
		return []*pb.VideoInfo{}, nil
	} else if err != nil {
		return nil, err
	}
	var videoIDList []int64
	for _, favorite := range favoriteList {
		videoIDList = append(videoIDList, favorite.VideoId)
	}
	videoSvrClient := utils.GetVideoSvrClient()
	if videoSvrClient == nil {
		return nil, fmt.Errorf("videoSvrClient is nil")
	}
	getVideoInfoListReq := &pb.GetVideoInfoListReq{VideoId: videoIDList}
	videoInfoListRsp, err := videoSvrClient.GetVideoInfoList(context.Background(), getVideoInfoListReq)
	if videoInfoListRsp == nil {
		return nil, fmt.Errorf("videoInfoList is nil")
	}
	return videoInfoListRsp.VideoInfoList, nil
}
