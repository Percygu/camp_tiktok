package controller

import (
	"gatewaysvr/config"
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

	resp, err := utils.NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName).CheckPassWord(ctx, &pb.CheckPassWordRequest{
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

	resp, err := utils.NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName).Register(ctx, &pb.RegisterRequest{
		Username: userName,
		Password: password,
	})
	if err != nil {
		zap.L().Error("register error", zap.Error(err))
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

	userInfo, err := utils.NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName).GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Id: uid,
	})
	if err != nil {
		zap.L().Error("get userinfo error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", userInfo)

}
