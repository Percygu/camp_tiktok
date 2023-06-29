package controller

import (
	"fmt"
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频流
func Feed(ctx *gin.Context) {
	var tokenId int64
	currentTime, err := strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	if err != nil || currentTime == int64(0) {
		currentTime = utils.GetCurrentTime()
	}
	// token := ctx.Query("token")
	// userId, err = common.VerifyToken(token)
	userId, _ := ctx.Get("UserId")
	tokenId = userId.(int64)

	// log.Info("currentTime:", currentTime, "tokenId:", tokenId)
	// if err != nil {
	// 	response.Fail(ctx, err.Error(), nil)
	// 	return
	// }

	// 这里还需要知道用户是否关注这个视频 作者 以及是否点赞
	feedListResponse, err := utils.GetVideoSvrClient().GetFeedList(ctx, &pb.GetFeedListRequest{
		CurrentTime: currentTime,
		TokenUserId: tokenId,
	})

	var resp = &DouyinFeedResponse{VideoList: make([]*Video, 0), NextTime: feedListResponse.NextTime}

	for _, video := range feedListResponse.VideoList {
		videoRsp := &Video{
			Id:            video.Id,
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
		// 调用远程方法获取视频作者信息
		GetUserInfoRsp, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
			Id: video.AuthorId,
		})
		if err != nil {
			log.Errorf("GetUserSvrClient GetUserInfo err %v", err.Error())
			response.Fail(ctx, fmt.Sprintf("GetUserSvrClient GetUserInfo err %v", err.Error()), nil)
			return
		}
		// 调用远程方法获取 我是否关注了这个视频的作者
		isFollowedRsp, err := utils.GetRelationSvrClient().IsFollowed(ctx, &pb.IsFollowedReq{
			SelfUserId: tokenId,
			ToUserId:   video.AuthorId,
		})
		if err != nil {
			log.Errorf("GetRelationSvrClient IsFollowed err %v", err.Error())
			response.Fail(ctx, fmt.Sprintf("GetRelationSvrClient IsFollowed err %v", err.Error()), nil)
			return
		}
		// 调用远程方法获取 我是否点赞了这个视频
		IsFavoriteVideoRsp, err := utils.GetFavoriteSvrClient().IsFavoriteVideo(ctx, &pb.IsFavoriteVideoReq{
			UserId:  tokenId,
			VideoId: video.Id,
		})
		if err != nil {
			log.Errorf("GetFavoriteSvrClient IsFavoriteVideo err %v", err.Error())
			response.Fail(ctx, fmt.Sprintf("GetFavoriteSvrClient IsFavoriteVideo err %v", err.Error()), nil)
			return
		}

		videoRsp.Author = GetUserInfoRsp.UserInfo
		videoRsp.Author.IsFollow = isFollowedRsp.IsFollowed
		videoRsp.IsFavorite = IsFavoriteVideoRsp.IsFavorite

	}
	response.Success(ctx, "success", resp)

}

type DouyinFeedResponse struct {
	StatusCode int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string   `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	VideoList  []*Video `protobuf:"bytes,3,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
	NextTime   int64    `protobuf:"varint,4,opt,name=next_time,json=nextTime,proto3" json:"next_time,omitempty"`
}
