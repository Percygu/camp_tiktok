package controller

import (
	"gatewaysvr/config"
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

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.NewVideoSvrClient(config.GetGlobalConfig().VideoServerConfig.Name).GetFeedList(ctx, &pb.GetFeedListRequest{
		CurrentTime: currentTime,
		TokenUserId: tokenId,
	})
	if err != nil {
		if err != nil {
			response.Fail(ctx, err.Error(), nil)
			return
		}

		response.Success(ctx, "success", resp.VideoList)
	}
}
