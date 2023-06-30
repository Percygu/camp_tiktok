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

type DouyinFeedResponse struct {
	StatusCode int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string   `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	VideoList  []*Video `protobuf:"bytes,3,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
	NextTime   int64    `protobuf:"varint,4,opt,name=next_time,json=nextTime,proto3" json:"next_time,omitempty"`
}

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
	// 这里还需要知道用户是否关注这个视频 作者 以及是否点赞
	feedListResponse, err := utils.GetVideoSvrClient().GetFeedList(ctx, &pb.GetFeedListRequest{
		CurrentTime: currentTime,
		TokenUserId: tokenId,
	})

	var resp = &DouyinFeedResponse{VideoList: make([]*Video, 0), NextTime: feedListResponse.NextTime}

	// 调用远程方法获取视频作者信息（一次性）
	var userIdList = make([]int64, 0)
	var followUintList = make([]*pb.FollowUint, 0)
	var favoriteUnitList = make([]*pb.FavoriteUnit, 0)
	for _, video := range feedListResponse.VideoList {
		userIdList = append(userIdList, video.AuthorId)
		followUintList = append(followUintList, &pb.FollowUint{
			SelfUserId: tokenId,
			// TODO: 命名不规范，应该是 toUserId
			UserIdList: video.AuthorId,
		})
		favoriteUnitList = append(favoriteUnitList, &pb.FavoriteUnit{
			UserId:  tokenId,
			VideoId: video.Id,
		})

	}
	getUserInfoRsp, err := utils.GetUserSvrClient().GetUserInfoDict(ctx, &pb.GetUserInfoDictRequest{
		UserIdList: userIdList,
	})
	if err != nil {
		log.Errorf("GetUserSvrClient GetUserInfoDict err %v", err.Error())
		response.Fail(ctx, fmt.Sprintf("GetUserSvrClient GetUserInfoDict err %v", err.Error()), nil)
		return
	}

	isFollowedRsp, err := utils.GetRelationSvrClient().IsFollowDict(ctx, &pb.IsFollowDictReq{
		FollowUintList: followUintList,
	})
	if err != nil {
		log.Errorf("GetRelationSvrClient IsFollowDict err %v", err.Error())
		response.Fail(ctx, fmt.Sprintf("GetRelationSvrClient IsFollowDict err %v", err.Error()), nil)
		return
	}

	isFavoriteVideoRsp, err := utils.GetFavoriteSvrClient().IsFavoriteVideoDict(ctx, &pb.IsFavoriteVideoDictReq{
		FavoriteUnitList: favoriteUnitList,
	})
	log.Infof("IsFavoriteVideoDict isFavoriteVideoRsp %+v", isFavoriteVideoRsp)

	if err != nil {
		log.Errorf("GetFavoriteSvrClient IsFavoriteVideoDict err %v", err.Error())
		response.Fail(ctx, fmt.Sprintf("GetFavoriteSvrClient IsFavoriteVideoDict err %v", err.Error()), nil)
		return
	}

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
		// 获取视频作者信息
		videoRsp.Author = getUserInfoRsp.UserInfoDict[video.AuthorId]
		var followUint = strconv.FormatInt(tokenId, 10) + "_" + strconv.FormatInt(videoRsp.Author.Id, 10)
		// 我是否关注了这个作者
		// if tokenId == videoRsp.Author.Id { // 自己的视频，不需要判断是否关注
		// 	videoRsp.Author.IsFollow = true
		// } else {
		videoRsp.Author.IsFollow = isFollowedRsp.IsFollowDict[followUint]

		var favoriteUint = strconv.FormatInt(tokenId, 10) + "_" + strconv.FormatInt(videoRsp.Id, 10)

		videoRsp.IsFavorite = isFavoriteVideoRsp.IsFavoriteDict[favoriteUint]
		resp.VideoList = append(resp.VideoList, videoRsp)
	}

	// bytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("json.Marshal err %v", err.Error())
		response.Fail(ctx, fmt.Sprintf("json.Marshal err %v", err.Error()), nil)
		return
	}

	// log.Info("Feed resp VideoList %v", resp.VideoList)
	// log.Infof("Feed resp json %v", string(bytes))

	for i := range resp.VideoList {
		log.Infof("Feed resp VideoList %v", resp.VideoList[i].Author.IsFollow, tokenId, resp.VideoList[i].Author.Id)
	}
	response.Success(ctx, "success", resp)
}
