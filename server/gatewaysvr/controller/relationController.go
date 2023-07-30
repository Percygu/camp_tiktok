package controller

import (
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DouyinRelationListResponse struct {
	StatusCode int32          `json:"status_code"`
	StatusMsg  string         `json:"status_msg,omitempty"`
	UserList   []*pb.UserInfo `json:"user_list,omitempty"`
}

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
	log.Infof("RelationAction tokenUserId:%d, toUid:%d, actionType:%d", tokenUserId, toUid, actionType)
	// 1.关注 2.取消关注
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

	// 获取关注列表
	getRelationFollowListRsp, err := utils.GetRelationSvrClient().GetRelationFollowList(ctx, &pb.GetRelationFollowListReq{
		UserId: uid,
	})

	if err != nil {
		log.Errorf("GetFollowList error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 去用户服务获取用户信息
	resp, err := utils.GetUserSvrClient().GetUserInfoList(ctx, &pb.GetUserInfoListRequest{
		IdList: getRelationFollowListRsp.FollowList,
	})
	if err != nil {
		log.Errorf("GetFollowList error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", &DouyinRelationListResponse{
		UserList: resp.UserInfoList,
	})
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
		log.Errorf("GetFollowerList ParseInt error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 获取关注者列表
	getRelationFollowerListRsp, err := utils.GetRelationSvrClient().GetRelationFollowerList(ctx, &pb.GetRelationFollowerListReq{
		UserId: uid,
	})
	if err != nil {
		log.Errorf("GetFollowerList error : %s", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 去用户服务获取用户信息
	resp, err := utils.GetUserSvrClient().GetUserInfoList(ctx, &pb.GetUserInfoListRequest{
		IdList: getRelationFollowerListRsp.FollowerList,
	})

	response.Success(ctx, "success", &DouyinRelationListResponse{
		UserList: resp.UserInfoList,
	})
}
