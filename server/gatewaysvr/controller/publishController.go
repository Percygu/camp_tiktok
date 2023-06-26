package controller

import (
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 视频发布
func PublishAction(ctx *gin.Context) {
	// publishResponse := &message.DouyinPublishActionResponse{}
	userId, _ := ctx.Get("UserId")
	// token := ctx.PostForm("token")
	// userId, err := common.VerifyToken(token)
	title := ctx.PostForm("title")
	data, err := ctx.FormFile("data")
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	filename := filepath.Base(data.Filename)

	finalName := fmt.Sprintf("%s_%s", utils.RandomString(), filename)
	videoPath := config.GetGlobalConfig().VideoPath
	saveFile := filepath.Join(videoPath, finalName)

	log.Infof("videoPath:%v", videoPath)

	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	publish, err := utils.NewVideoSvrClient(config.GetGlobalConfig().SvrConfig.VideoSvrName).PublishVideo(ctx, &pb.PublishVideoRequest{
		UserId:   userId.(int64),
		Title:    title,
		SaveFile: saveFile,
	})

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	log.Infof("publish:%v", err)
	response.Success(ctx, "success", publish)

}

// 获取视频列表
func GetPublishList(ctx *gin.Context) {
	tokenUserId, _ := ctx.Get("UserId")
	id := ctx.Query("user_id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
	}

	resp, err := utils.NewVideoSvrClient(config.GetGlobalConfig().SvrConfig.VideoSvrName).GetPublishVideoList(ctx, &pb.GetPublishVideoListRequest{
		TokenUserId: tokenUserId.(int64),
		UserID:      userId,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", resp)
}
