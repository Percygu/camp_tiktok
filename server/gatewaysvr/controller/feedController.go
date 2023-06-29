package controller

import (
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
	resp, err := utils.GetVideoSvrClient().GetFeedList(ctx, &pb.GetFeedListRequest{
		CurrentTime: currentTime,
		TokenUserId: tokenId,
	})

	// 填充是否关注
	for i, video := range resp.VideoList {
		vi
		utils.GetRelationSvrClient().GetRelationFollowList(ctx, &pb.GetRelationFollowListReq{
			UserId: tokenId,
		})
	}

	// m := make(map[int64]struct{})
	// list, err := repository.GetFollowList(userId, "follow")
	// if err != nil {
	// 	return nil, err
	// }
	// for _, u := range list {
	// 	m[u.Id] = struct{}{}
	// }

	// log.Info("resp:", resp.VideoList, err)

	if err != nil {
		if err != nil {
			response.Fail(ctx, err.Error(), nil)
			return
		}

		response.Success(ctx, "success", resp.VideoList)
	}
}

type DouyinFeedResponse struct {
	StatusCode int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string   `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	VideoList  []*Video `protobuf:"bytes,3,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
	NextTime   int64    `protobuf:"varint,4,opt,name=next_time,json=nextTime,proto3" json:"next_time,omitempty"`
}
