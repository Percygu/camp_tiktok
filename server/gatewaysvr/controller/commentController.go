package controller

import (
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"go.uber.org/zap"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 发布评论
func CommentAction(ctx *gin.Context) {
	var err error

	tokenUids, _ := ctx.Get("UserId")

	tokenUid := tokenUids.(int64)

	video_id := ctx.Query("video_id")
	comment_text := ctx.Query("comment_text")
	actionTypeStr := ctx.Query("action_type")
	comment_id := ctx.Query("comment_id")
	commentId := int64(0)
	if actionTypeStr == "2" {
		commentId, err = strconv.ParseInt(comment_id, 10, 64)
		if err != nil {
			zap.L().Error("commentId error", zap.Error(err))
			response.Fail(ctx, err.Error(), nil)
			return
		}
	}
	videoId, err := strconv.ParseInt(video_id, 10, 64)
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		zap.L().Error("videoId error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.GetCommentSvrClient().CommentAction(ctx, &pb.CommentActionReq{
		UserId:      tokenUid,
		VideoId:     videoId,
		CommentId:   commentId,
		CommentText: comment_text,
		ActionType:  actionType,
	})

	if err != nil {
		zap.L().Error("comment error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}

// 获取评论列表
func GetCommentList(ctx *gin.Context) {
	var err error
	video_id := ctx.Query("video_id")
	/* token := ctx.Query("token")
	_, err = util.VerifyToken(token)
	if err != nil {
		log.Errorf("token error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	videoId, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		zap.L().Error("videoId error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.GetCommentSvrClient().GetCommentList(ctx, &pb.GetCommentListReq{
		VideoId: videoId,
	})

	if err != nil {
		zap.L().Error("commentList error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}
