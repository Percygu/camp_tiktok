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

// DouyinPublishActionResponse PublishAction返回的数据结构
type DouyinPublishActionResponse struct {
	StatusCode int32  `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
}

// DouyinPublishListResponse GetPublishList返回的数据结构
type DouyinPublishListResponse struct {
	StatusCode int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string   `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	VideoList  []*Video `protobuf:"bytes,3,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
}

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

	_, err = utils.GetVideoSvrClient().PublishVideo(ctx, &pb.PublishVideoRequest{
		UserId:   userId.(int64),
		Title:    title,
		SaveFile: saveFile,
	})

	if err != nil {
		log.Errorf("utils.GetVideoSvrClient().PublishVideo err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", &DouyinPublishActionResponse{})

}

// 获取自己发布视频列表
func GetPublishList(ctx *gin.Context) {
	tokenUserId, _ := ctx.Get("UserId")
	id := ctx.Query("user_id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
	}

	getPublishVideoList, err := utils.GetVideoSvrClient().GetPublishVideoList(ctx, &pb.GetPublishVideoListRequest{
		TokenUserId: tokenUserId.(int64),
		UserID:      userId,
	})
	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 获取我自己的用户信息
	getUserInfo, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
		Id: tokenUserId.(int64),
	})
	if err != nil {
		log.Errorf("utils.GetUserSvrClient().GetUserInfo err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	video := make([]*Video, 0, len(getPublishVideoList.VideoList))
	for _, v := range getPublishVideoList.VideoList {
		video = append(video, &Video{
			Id:            v.Id,
			Author:        getUserInfo.UserInfo,
			Title:         v.Title,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			// TODO:
			IsFavorite: false, // 自己是否喜欢
		})
	}

	response.Success(ctx, "success", &DouyinPublishListResponse{
		VideoList: video,
	})
}
