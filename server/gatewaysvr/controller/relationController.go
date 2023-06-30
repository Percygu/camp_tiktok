package controller

import (
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 关注操作
func RelationAction(ctx *gin.Context) {
	tokens, _ := ctx.Get("UserId")
	tokenUserId := tokens.(int64)

	toUserId := ctx.Query("to_user_id")
	toUid, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	actionStr := ctx.Query("action_type")

	actionType, err := strconv.ParseInt(actionStr, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	_, err = utils.GetRelationSvrClient().RelationAction(ctx, &pb.RelationActionReq{
		ToUserId:   toUid,
		SelfUserId: tokenUserId,
		ActionType: actionType,
	})
	if err != nil {
		log.Errorf("RelationAction error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 我关注的人++
	_, err = utils.GetUserSvrClient().UpdateUserFollowCount(ctx, &pb.UpdateUserFollowCountReq{
		UserId:     tokenUserId,
		ActionType: actionType,
	})
	if err != nil {
		log.Errorf("RelationAction error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 被我关注的人的粉丝++
	_, err = utils.GetUserSvrClient().UpdateUserFollowerCount(ctx, &pb.UpdateUserFollowerCountReq{
		UserId:     toUid,
		ActionType: actionType,
	})
	if err != nil {
		log.Errorf("RelationAction error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", nil)
}

// 获取关注列表
func GetFollowList(ctx *gin.Context) {
	// token := ctx.Query("token")
	// tokenUserId, err := common.VerifyToken(token)
	/* if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	// tokens, _ := ctx.Get("UserId")
	// tokenUserId := tokens.(int64)

	UserId := ctx.Query("user_id")
	uid, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.GetRelationSvrClient().GetRelationFollowList(ctx, &pb.GetRelationFollowListReq{
		UserId: uid,
	})

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}

// 获取关注者列表
func GetFollowerList(ctx *gin.Context) {
	/* token := ctx.Query("token")
	tokenUserId, err := common.VerifyToken(token)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	} */
	// tokens, _ := ctx.Get("UserId")
	// tokenUserId := tokens.(int64)

	UserId := ctx.Query("user_id")
	uid, err := strconv.ParseInt(UserId, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	resp, err := utils.NewRelationSvrClient(config.GetGlobalConfig().SvrConfig.RelationSvrName).GetRelationFollowerList(ctx, &pb.GetRelationFollowerListReq{
		UserId: uid,
	})

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}
