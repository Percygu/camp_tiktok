package repository

import (
	"gorm.io/gorm"
	"time"
	"videosvr/middleware/db"
)

func GetVideoList(AuthorId int64) ([]Video, error) {
	var videos []Video
	// userSvrClient := utils.NewUserSvrClient(config.GetGlobalConfig().UserSvrName)
	// reply, err := userSvrClient.GetUserInfo(context.Background(), &pb.GetUserInfoRequest{Id: AuthorId})
	// if err != nil {
	// 	return videos, err
	// }
	db := db.GetDB()
	err := db.Where("author_id = ?", AuthorId).Order("video_id DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return videos, err
	}
	// userInfo := reply.UserInfo
	// for i := range videos {
	// 	videos[i].Author = User{
	// 		Id:              userInfo.Id,
	// 		Name:            userInfo.Name,
	// 		Follow:          userInfo.FollowCount,
	// 		Follower:        userInfo.FollowerCount,
	// 		Avatar:          userInfo.Avatar,
	// 		BackgroundImage: userInfo.BackgroundImage,
	// 		Signature:       userInfo.Signature,
	// 		TotalFav:        userInfo.TotalFavorited,
	// 		FavCount:        userInfo.FavoriteCount,
	// 	}
	// }
	return videos, nil
}

func InsertVideo(authorId int64, playUrl, coverUrl, title string) error {
	video := Video{
		AuthorId:      authorId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		PublishTime:   time.Now().UnixNano() / 1e6,
		Title:         title,
	}
	db := db.GetDB()
	err := db.Create(&video).Error
	if err != nil {
		return err
	}
	return nil
}

func GetVideoListByFeed(currentTime int64) ([]Video, error) {
	var videos []Video
	db := db.GetDB()
	err := db.Where("publish_time < ?", currentTime).Limit(20).Order("video_id DESC").Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return videos, err
	}

	// for i, v := range videos {
	// 	author, err := utils.NewUserSvrClient(config.GetGlobalConfig().UserSvrName).GetUserInfo(context.Background(), &pb.GetUserInfoRequest{Id: v.AuthorId})
	// 	err = CacheSetAuthor(v.Id, v.AuthorId)
	// 	if err != nil {
	// 		return videos, err
	// 	}
	// 	videos[i].Author = User{
	// 		Id:              author.UserInfo.Id,
	// 		Name:            author.UserInfo.Name,
	// 		Follow:          author.UserInfo.FollowCount,
	// 		Follower:        author.UserInfo.FollowerCount,
	// 		Avatar:          author.UserInfo.Avatar,
	// 		BackgroundImage: author.UserInfo.BackgroundImage,
	// 		Signature:       author.UserInfo.Signature,
	// 		TotalFav:        author.UserInfo.TotalFavorited,
	// 		FavCount:        author.UserInfo.FavoriteCount,
	// 	}
	// }
	return videos, nil
}
