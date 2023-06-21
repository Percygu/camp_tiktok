package controller

import (
	"fmt"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/goccy/go-json"
	"github.com/mgechev/revive/config"
	"gorm.io/gorm/logger"
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
	videoPath := config.GetConfig().Path.Videofile
	saveFile := filepath.Join(videoPath, finalName)

	logger.Info("saveFile:", saveFile)
	// 将data序列化成json 字符串
	dataJson, err := json.Marshal(data)

	if err := ctx.SaveUploadedFile(data, saveFile); err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	publish, err := service.PublishVideo(userId.(int64), saveFile, title)
	// publish, err := service.PublishVideo(userId, saveFile)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	logger.Infof("publish:%v", publish)
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
	list, err := service.PublishList(tokenUserId.(int64), userId)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}
	response.Success(ctx, "success", list)
}
