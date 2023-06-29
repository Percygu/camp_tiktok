package controller

import (
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FavActionParams struct {
	// 暂时没 user_id ，因为客户端出于安全考虑没给出
	Token      string `form:"token" binding:"required"`
	VideoId    int64  `form:"video_id" binding:"required"`
	ActionType int8   `form:"action_type" binding:"required,oneof=1 2"`
}

type FavListParams struct {
	Token  string `form:"token" binding:"required"`
	UserId int64  `form:"user_id" binding:"required"`
}

// 点赞视频
func FavoriteAction(ctx *gin.Context) {
	var favInfo FavActionParams
	err := ctx.ShouldBindQuery(&favInfo)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	tokenUidStr, _ := ctx.Get("UserId")
	tokenUid := tokenUidStr.(int64)

	if err != nil {
		zap.L().Error("token error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.GetFavoriteSvrClient().FavoriteAction(ctx, &pb.FavoriteActionReq{
		UserId:     tokenUid,
		VideoId:    favInfo.VideoId,
		ActionType: int64(favInfo.ActionType),
	})

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}

// 获取点赞列表
func GetFavoriteList(ctx *gin.Context) {

	tokenUidStr, _ := ctx.Get("UserId")
	tokenUid := tokenUidStr.(int64)

	// 拿videoID List
	resp, err := utils.GetVideoSvrClient().GetFavoriteVideoList(ctx, &pb.GetFavoriteVideoListReq{
		UserId: tokenUid,
	})

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}
