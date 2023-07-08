package service

import (
	"context"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"go.uber.org/zap"
	"strconv"
	"time"
	"videosvr/log"
	"videosvr/middleware/minioStore"
	"videosvr/repository"
	"videosvr/utils"
)

type VideoService struct {
	pb.UnimplementedVideoServiceServer
}

// UpdateFavoriteCount 更新点赞数
func (v VideoService) UpdateFavoriteCount(ctx context.Context, req *pb.UpdateFavoriteCountReq) (*pb.UpdateFavoriteCountRsp, error) {
	err := repository.UpdateFavoriteNum(req.VideoId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateFavoriteCount err", zap.Error(err))
		return nil, err
	}
	return &pb.UpdateFavoriteCountRsp{}, nil
}

// UpdateCommentCount 更新评论数
func (v VideoService) UpdateCommentCount(ctx context.Context, req *pb.UpdateCommentCountReq) (*pb.UpdateCommentCountRsp, error) {
	err := repository.UpdateCommentNum(req.VideoId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateCommentCount err", zap.Error(err))
		return nil, err
	}
	return &pb.UpdateCommentCountRsp{}, nil
}

func (v VideoService) GetPublishVideoList(ctx context.Context, req *pb.GetPublishVideoListRequest) (*pb.GetPublishVideoListResponse, error) {
	videos, err := repository.GetVideoListByAuthorId(req.UserID)
	if err != nil {
		log.Errorf("GetVideoListByAuthorId err", zap.Error(err))
		return nil, err
	}
	// 结构体转换
	videoList := make([]*pb.VideoInfo, 0)
	for _, video := range videos {
		videoList = append(videoList, &pb.VideoInfo{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
		})
	}

	return &pb.GetPublishVideoListResponse{
		VideoList: videoList,
	}, nil
}
func (v VideoService) PublishVideo(ctx context.Context, req *pb.PublishVideoRequest) (*pb.PublishVideoResponse, error) {
	client := minioStore.GetMinio()
	videoUrl, err := client.UploadFile("video", req.SaveFile, strconv.FormatInt(req.UserId, 10))
	log.Infof("save file: %v", req.SaveFile)
	if err != nil {
		log.Errorf("UploadFile err", zap.Error(err))
		return nil, err
	}

	// 生成视频封面（截取第一桢）
	imageFile, err := utils.GetImageFile(req.SaveFile)

	if err != nil {
		log.Errorf("GetImageFile err", zap.Error(err))
		return nil, err
	}

	picUrl, err := client.UploadFile("pic", imageFile, strconv.FormatInt(req.UserId, 10))
	if err != nil {
		log.Errorf("UploadFile err", zap.Error(err))
		picUrl = "https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7909abe413ec4a1e82032d2beb810157~tplv-k3u1fbpfcp-zoom-in-crop-mark:1304:0:0:0.awebp?"
	}

	err = repository.InsertVideo(req.UserId, videoUrl, picUrl, req.Title)
	if err != nil {
		log.Errorf("InsertVideo err", zap.Error(err))
		return nil, err
	}
	return &pb.PublishVideoResponse{}, nil
}

func (v VideoService) GetFeedList(ctx context.Context, req *pb.GetFeedListRequest) (*pb.GetFeedListResponse, error) {
	// 拿出一批视频
	videoList, err := repository.GetVideoListByFeed(req.CurrentTime)
	if err != nil {
		log.Errorf("GetFeedList|GetVideoListByFeed err:%v", err)
		return nil, err
	}
	nextTime := time.Now().UnixNano() / 1e6
	if len(videoList) == 20 {
		nextTime = videoList[len(videoList)-1].PublishTime
	}

	var VideoInfoList []*pb.VideoInfo
	for _, video := range videoList {
		VideoInfoList = append(VideoInfoList, &pb.VideoInfo{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false, // 是否喜欢，在gateway处理
			Title:         video.Title,
		})
	}

	resp := &pb.GetFeedListResponse{
		NextTime:  nextTime,
		VideoList: VideoInfoList,
	}

	log.Infof("GetFeedList|resp:%+v", resp)
	return resp, nil
}

// GetFavoriteVideoList 获取用户喜欢的视频列表
func (v VideoService) GetFavoriteVideoList(ctx context.Context, req *pb.GetFavoriteVideoListReq) (*pb.GetFavoriteVideoListRsp, error) {
	// 获取用户喜欢的视频id列表（这个得调用favorite服务处理）
	resp, err := utils.GetFavoriteSvrClient().GetFavoriteVideoIdList(ctx, &pb.GetFavoriteVideoIdListReq{
		UserId: req.UserId,
	})
	if err != nil {
		log.Errorf("GetFavoriteVideoList | GetFavoriteVideoIdList err: %v", err)
		return nil, err
	}

	videoInfoListRsp, err := v.GetVideoInfoList(ctx, &pb.GetVideoInfoListReq{
		VideoId: resp.VideoIdList,
	})
	if err != nil {
		log.Errorf("GetFavoriteVideoList | GetVideoInfoList err: %v", err)
		return nil, err
	}

	if videoInfoListRsp == nil {
		return nil, nil
	}

	// 返回
	return &pb.GetFavoriteVideoListRsp{
		VideoList: videoInfoListRsp.VideoInfoList,
	}, nil
}

// GetVideoInfoList 通过video_id列表 获取 对应的视频信息
func (v VideoService) GetVideoInfoList(ctx context.Context, req *pb.GetVideoInfoListReq) (*pb.GetVideoInfoListRsp, error) {
	videoList, err := repository.GetVideoListByVideoIdList(req.VideoId)
	if err != nil {
		return nil, err
	}

	videoInfoList := make([]*pb.VideoInfo, len(videoList))
	for i, video := range videoList {
		videoInfoList[i] = &pb.VideoInfo{
			Id:            video.Id,
			AuthorId:      video.AuthorId,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    false,
			Title:         video.Title,
		}
	}
	return &pb.GetVideoInfoListRsp{
		VideoInfoList: videoInfoList,
	}, nil
}

func VideoInfo(videoList []repository.Video, userId int64) []*pb.VideoInfo {
	// var err error
	// FollowList := make(map[int64]struct{})
	// favList := make(map[int64]struct{})
	// if userId != int64(0) {
	// 	FollowList, err = tokenFollowList(userId)
	// 	if err != nil {
	// 		return nil
	// 	}
	// 	favList, err = tokenFavList(userId)
	// 	if err != nil {
	// 		return nil
	// 	}
	// }
	//
	// lists := make([]*pb.VideoInfo, len(videoList))
	// for i, video := range videoList {
	// 	v := &pb.VideoInfo{
	// 		Id:            video.Id,
	// 		PlayUrl:       video.PlayUrl,
	// 		CoverUrl:      video.CoverUrl,
	// 		FavoriteCount: video.FavoriteCount,
	// 		CommentCount:  video.CommentCount,
	// 		IsFavorite:    false,
	// 		// Author:        messageUserInfo(video.Author),
	// 		Title: video.Title,
	// 	}
	//
	// 	if _, ok := FollowList[video.AuthorId]; ok {
	// 		v.Author.IsFollow = true
	// 	}
	// 	if _, ok := favList[video.Id]; ok {
	// 		v.IsFavorite = true
	// 	}
	//
	// 	lists[i] = v
	// }
	// return lists
	return nil
}

func tokenFollowList(userId int64) (map[int64]struct{}, error) {
	m := make(map[int64]struct{})
	reply, err := utils.GetRelationSvrClient().GetRelationFollowList(context.Background(), &pb.GetRelationFollowListReq{
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

// func tokenFavList(tokenUserId int64) (map[int64]struct{}, error) {
// 	m := make(map[int64]struct{})
//
// 	reply, err := utils.NewFavoriteSvrClient(config.GetGlobalConfig().SvrConfig.FavoriteSvrName).GetFavoriteVideoList(context.Background(), &pb.GetFavoriteVideoListReq{
// 		UserId: tokenUserId,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	list := reply.VideoInfoList
// 	for _, v := range list {
// 		m[v.Id] = struct{}{}
// 	}
// 	return m, nil
// }

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
