package controller

import (
	"TikTokLite/response"
	"TikTokLite/service"
	"TikTokLite/web/global"
	"TikTokLite/web/proto"
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

	resp, err := global.UserSrvClient.CheckPassWord(ctx, &proto.CheckPassWordRequest{
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
	registerResponse, err := service.UserRegister(userName, password)
	if err != nil {
		zap.L().Error("register error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", registerResponse)

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

	userinfo, err := service.UserInfo(uid)
	if err != nil {
		zap.L().Error("get userinfo error", zap.Error(err))
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", userinfo)

}