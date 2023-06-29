package service

import (
	"context"
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"go.uber.org/zap"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"videosvr/config"
	"videosvr/log"
	"videosvr/middleware/minioStore"
	"videosvr/repository"
	"videosvr/utils"
)

type VideoService struct {
	pb.UnimplementedVideoServiceServer
}

func (v VideoService) GetPublishVideoList(ctx context.Context, req *pb.GetPublishVideoListRequest) (*pb.GetPublishVideoListResponse, error) {
	videos, err := repository.GetVideoListByAuthorId(req.UserID)
	if err != nil {
		return nil, err
	}
	list := &pb.GetPublishVideoListResponse{
		VideoList: VideoInfo(videos, req.TokenUserId),
	}

	return list, nil
}
func (v VideoService) PublishVideo(ctx context.Context, req *pb.PublishVideoRequest) (*pb.PublishVideoResponse, error) {
	client := minioStore.GetMinio()
	videoUrl, err := client.UploadFile("video", req.SaveFile, strconv.FormatInt(req.UserId, 10))
	log.Infof("save file: %v", req.SaveFile)
	if err != nil {
		return nil, err
	}

	// 生成视频封面（截取第一桢）
	imageFile, err := GetImageFile(req.SaveFile)

	if err != nil {
		return nil, err
	}

	log.Infof("imageFile", zap.String("imageFile", imageFile))

	picUrl, err := client.UploadFile("pic", imageFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		picUrl = "https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7909abe413ec4a1e82032d2beb810157~tplv-k3u1fbpfcp-zoom-in-crop-mark:1304:0:0:0.awebp?"
	}

	err = repository.InsertVideo(req.UserId, videoUrl, picUrl, req.Title)
	if err != nil {
		return nil, err
	}
	return &pb.PublishVideoResponse{}, nil
}

func (v VideoService) GetFeedList(ctx context.Context, req *pb.GetFeedListRequest) (*pb.GetFeedListResponse, error) {
	videoList, err := repository.GetVideoListByFeed(req.CurrentTime)
	if err != nil {
		return nil, err
	}

	feed := &pb.GetFeedListResponse{
		VideoList: VideoInfo(videoList, req.TokenUserId),
	}

	nextTime := time.Now().UnixNano() / 1e6
	if len(videoList) == 20 {
		nextTime = videoList[len(videoList)-1].PublishTime
	}
	feed.NextTime = nextTime
	return feed, nil
}

func (v VideoService) GetFavoriteVideoList(ctx context.Context, req *pb.GetFavoriteVideoListReq) (*pb.GetFavoriteVideoListRsp, error) {
	resp, err := utils.GetFavoriteSvrClient().GetFavoriteVideoIdList(ctx, &pb.GetFavoriteVideoIdListReq{
		UserId: req.UserId,
	})

	videoInfoListRsp, err := v.GetVideoInfoList()
	if videoInfoListRsp == nil {
		return nil, fmt.Errorf("videoInfoList is nil")
	}
	return videoInfoListRsp.VideoInfoList, nil
}

func (v VideoService) GetVideoInfoList(ctx context.Context, req *pb.GetVideoInfoListReq) (*pb.GetVideoInfoListRsp, error) {
	videoList, err := repository.GetVideoListByVideoIdList(req.VideoId)
	if err != nil {
		return nil, err
	}

	return &pb.GetVideoInfoListRsp{
		VideoInfoList: videoList,
	}, nil
}

func VideoInfo(videoList []repository.Video, userId int64) []*pb.VideoInfo {
	var err error
	FollowList := make(map[int64]struct{})
	favList := make(map[int64]struct{})
	if userId != int64(0) {
		FollowList, err = tokenFollowList(userId)
		if err != nil {
			return nil
		}
		favList, err = tokenFavList(userId)
		if err != nil {
			return nil
		}
	}
	lists := make([]*pb.VideoInfo, len(videoList))
	for i, video := range videoList {
		v := &pb.VideoInfo{
			Id:            video.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false,
			// Author:        messageUserInfo(video.Author),
			Title: video.Title,
		}
		if _, ok := FollowList[video.AuthorId]; ok {
			v.Author.IsFollow = true
		}
		if _, ok := favList[video.Id]; ok {
			v.IsFavorite = true
		}
		lists[i] = v
	}
	return lists
}

func tokenFollowList(userId int64) (map[int64]struct{}, error) {
	m := make(map[int64]struct{})
	reply, err := utils.NewRelationSvrClient(config.GetGlobalConfig().SvrConfig.RelationSvrName).GetRelationFollowList(context.Background(), &pb.GetRelationFollowListReq{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	list := reply.UserInfoList
	for _, u := range list {
		m[u.Id] = struct{}{}
	}
	return m, nil
}

func tokenFavList(tokenUserId int64) (map[int64]struct{}, error) {
	m := make(map[int64]struct{})

	reply, err := utils.NewFavoriteSvrClient(config.GetGlobalConfig().SvrConfig.FavoriteSvrName).GetFavoriteVideoList(context.Background(), &pb.GetFavoriteVideoListReq{
		UserId: tokenUserId,
	})
	if err != nil {
		return nil, err
	}

	list := reply.VideoInfoList
	for _, v := range list {
		m[v.Id] = struct{}{}
	}
	return m, nil
}

// func messageUserInfo(info repository.User) *pb.UserInfo {
// 	return &pb.UserInfo{
// 		Id:              info.Id,
// 		Name:            info.Name,
// 		FollowCount:     info.Follow,
// 		FollowerCount:   info.Follower,
// 		IsFollow:        false,
// 		Avatar:          info.Avatar,
// 		BackgroundImage: info.BackgroundImage,
// 		Signature:       info.Signature,
// 		TotalFavorited:  info.TotalFav,
// 		FavoriteCount:   info.FavCount,
// 	}
// }

func GetImageFile(videoPath string) (string, error) {
	temp := strings.Split(videoPath, "/")
	videoName := temp[len(temp)-1]
	b := []byte(videoName)
	videoName = string(b[:len(b)-3]) + "jpg"
	picPath := config.GetGlobalConfig().MinioConfig.PicPath
	picName := filepath.Join(picPath, videoName)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "1", "-f", "image2", "-t", "0.01", "-y", picName)
	err := cmd.Run()
	if err != nil {
		log.Errorf("cmd.Run() failed with %s\n", err)
		return "", err
	}
	return picName, nil
}