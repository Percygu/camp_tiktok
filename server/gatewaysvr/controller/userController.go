package controller

import (
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// 用户登录
func UserLogin(ctx *gin.Context) {
	var err error
	userName := ctx.Query("username")
	password := ctx.Query("password")
	if len(userName) > 32 || len(password) > 32 { // 最长32位字符
		response.Fail(ctx, "username or password invalid", nil)
		return
	}
	resp, err := utils.GetUserSvrClient().CheckPassWord(ctx, &pb.CheckPassWordRequest{
		Username: userName,
		Password: password,
	})
	if err != nil {
		zap.L().Error("login error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}

func UserRegister(ctx *gin.Context) {
	var err error
	userName := ctx.Query("username")
	password := ctx.Query("password")
	if len(userName) > 32 || len(password) > 32 { // 最长32位字符
		response.Fail(ctx, "username or password invalid", nil)
		return
	}
	log.Info(userName, password)
	resp, err := utils.GetUserSvrClient().Register(ctx, &pb.RegisterRequest{
		Username: userName,
		Password: password,
	})

	if err != nil {
		log.Error("register error", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)

}

// 获取用户信息
func GetUserInfo(ctx *gin.Context) {

	var err error
	userId := ctx.Query("user_id")
	uids, _ := ctx.Get("UserId")

	uid := uids.(int64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	if strconv.FormatInt(uid, 10) != userId {
		response.Fail(ctx, "token error", nil)
		return
	}

	getUserInfoRsp, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Id: uid,
	})
	if err != nil {
		log.Errorf("GetUserInfo err %v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	log.Infof("GetUserInfo %+v", getUserInfoRsp)
	response.Success(ctx, "success", &DouyinUserResponse{
		User: getUserInfoRsp.UserInfo,
	})
}

type DouyinUserResponse struct {
	StatusCode int32        `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string       `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	User       *pb.UserInfo `protobuf:"bytes,3,opt,name=user,proto3" json:"user"`
}
