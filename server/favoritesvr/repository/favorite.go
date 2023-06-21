package repository

import (
	"favoritesvr/config"
	"favoritesvr/log"
	db "favoritesvr/middleware/db"
	"fmt"
	"gatewaysvr/utils/otgrpc"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func LikeAction(uid, vid int64) error {
	db := db.GetDB()
	favorite := Favorite{
		UserId:  uid,
		VideoId: vid,
	}
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Find(&Favorite{}).Error
	if err != gorm.ErrRecordNotFound {
		return fmt.Errorf("you have liked this video")
	}
	err = db.Create(&favorite).Error
	if err != nil {
		return err
	}
	//authorid, _ := CacheGetAuthor(vid) // todo videosvr
	// todo usercountcache change usersvr
	//go CacheChangeUserCount(uid, add, "like")
	//go CacheChangeUserCount(authorid, add, "liked")
	return nil
}

func UnLikeAction(uid, vid int64) error {
	db := db.GetDB()
	err := db.Where("user_id = ? and video_id = ?", uid, vid).Delete(&Favorite{}).Error
	if err != nil {
		return err
	}
	//authorid, _ := CacheGetAuthor(vid)
	// go func() {
	//go CacheChangeUserCount(uid, sub, "like")
	//go CacheChangeUserCount(authorid, sub, "liked")
	// }()
	return nil
}

func GetFavoriteList(uid int64) ([]*pb.VideoInfo, error) {
	var videos []*pb.VideoInfo
	db := db.GetDB()
	var favoriteList []*Favorite
	err := db.Model(&Favorite{}).Where("user_id= ?", uid).Find(&favoriteList).Error
	if err == gorm.ErrRecordNotFound {
		return videos, nil
	} else if err != nil {
		return nil, err
	}
	var videoIDList []int64
	for _, favorite := range favoriteList {
		videoIDList = append(videoIDList, favorite.VideoId)
	}
	videoSvrClient := NewVideoSvrClient(config.GetGlobalConfig().VideoSvrName)
	if videoSvrClient == nil {
		return nil, fmt.Errorf("videoSvrClient is nil")
	}
	videoInfoList := videoSvrClient.get

	consulInfo := config.GetGlobalConfig().ConsulConfig

	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port,
			config.GetGlobalConfig().VideoSvrName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, fmt.Errorf("连接用户服务失败")
	}

	err := db.Joins("left join favorites on videos.video_id = favorites.video_id").
		Where("favorites.user_id = ?", uid).Find(&videos).Error
	if err == gorm.ErrRecordNotFound {
		return videos, nil
	} else if err != nil {
		return nil, err
	}

	// global.UserSrvClient = proto.NewUserClient(userConn)

	for i, v := range videos {
		author, err := GetUserInfo(v.AuthorId)
		if err != nil {
			return videos, err
		}
		videos[i].Author = author
	}
	return videos, nil
}

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	return conn, nil
}

func NewVideoSvrClient(svrName string) pb.VideoClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewVideoClient(conn)
}

func NewUserSvrClient(svrName string) pb.VideoServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewVideoServiceClient(conn)
}
