package controller

import (
	"gatewaysvr/log"
	"gatewaysvr/response"
	"gatewaysvr/utils"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"github.com/gin-gonic/gin"
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

// DouyinFavoriteListResponse  GetFavoriteList返回的结构体
type DouyinFavoriteListResponse struct {
	StatusCode int32    `protobuf:"varint,1,opt,name=status_code,json=statusCode,proto3" json:"status_code"`
	StatusMsg  string   `protobuf:"bytes,2,opt,name=status_msg,json=statusMsg,proto3" json:"status_msg,omitempty"`
	VideoList  []*Video `protobuf:"bytes,3,rep,name=video_list,json=videoList,proto3" json:"video_list,omitempty"`
}
type Video struct {
	Id            int64        `json:"id"`
	Author        *pb.UserInfo `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

// *******************************************

// 点赞或取消点赞 一个视频
func FavoriteAction(ctx *gin.Context) {
	var favInfo FavActionParams
	err := ctx.ShouldBindQuery(&favInfo)
	if err != nil {
		log.Errorf("ShouldBindQuery failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}
	tokenUidStr, _ := ctx.Get("UserId")
	tokenUid := tokenUidStr.(int64)

	// 这里只是插入了一条favorite记录，没有更新video表的favorite_count，还有对应视频作者的favorite_count
	_, err = utils.GetFavoriteSvrClient().FavoriteAction(ctx, &pb.FavoriteActionReq{
		UserId:     tokenUid,
		VideoId:    favInfo.VideoId,
		ActionType: int64(favInfo.ActionType),
	})

	if err != nil {
		log.Errorf("FavoriteAction failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 更新video表的favorite_count（更新视频获赞数）
	_, err = utils.GetVideoSvrClient().UpdateFavoriteCount(ctx, &pb.UpdateFavoriteCountReq{
		VideoId:    favInfo.VideoId,
		ActionType: int64(favInfo.ActionType),
	})
	if err != nil {
		log.Errorf("UpdateFavoriteCount failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 查询video表的author_id
	videoInfoResp, err := utils.GetVideoSvrClient().GetVideoInfoList(ctx, &pb.GetVideoInfoListReq{
		VideoId: []int64{favInfo.VideoId},
	})
	if err != nil {
		log.Errorf("GetVideoInfoList failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	var authorId = videoInfoResp.VideoInfoList[0].AuthorId
	// 更新user表的 total_favorited_count（更新视频作者获赞数）
	_, err = utils.GetUserSvrClient().UpdateUserFavoritedCount(ctx, &pb.UpdateUserFavoritedCountReq{
		UserId:     authorId,
		ActionType: int64(favInfo.ActionType),
	})
	if err != nil {
		log.Errorf("UpdateUserFavoritedCount failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// 更新user表的 favorite_count更新我喜欢的视频数）
	_, err = utils.GetUserSvrClient().UpdateUserFavoriteCount(ctx, &pb.UpdateUserFavoriteCountReq{
		UserId:     tokenUid,
		ActionType: int64(favInfo.ActionType),
	})
	if err != nil {
		log.Errorf("UpdateUserFavoriteCount failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", nil)
}

// 获取点赞列表
func GetFavoriteList(ctx *gin.Context) {

	tokenUidStr, _ := ctx.Get("UserId")
	tokenUid := tokenUidStr.(int64)

	// 拿videoID List
	videoListResp, err := utils.GetVideoSvrClient().GetFavoriteVideoList(ctx, &pb.GetFavoriteVideoListReq{
		UserId: tokenUid,
	})

	var userIdList []int64
	for _, v := range videoListResp.VideoList {
		userIdList = append(userIdList, v.AuthorId)
	}

	// 拿userInfo List
	userInfoResp, err := utils.GetUserSvrClient().GetUserInfoList(ctx, &pb.GetUserInfoListRequest{
		IdList: userIdList,
	})

	if err != nil {
		log.Errorf("GetUserInfoList failed, err:%v", err)
		response.Fail(ctx, err.Error(), nil)
		return
	}

	// log.Info(userInfoResp)

	userMap := make(map[int64]*pb.UserInfo)
	for _, v := range userInfoResp.UserInfoList {
		userMap[v.Id] = v
	}

	videoList := make([]*Video, 0)

	for _, v := range videoListResp.VideoList {
		videoList = append(videoList, &Video{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
			Author: &pb.UserInfo{
				Id:              v.AuthorId,
				Name:            userMap[v.AuthorId].Name,
				Avatar:          userMap[v.AuthorId].Avatar,
				FollowCount:     userMap[v.AuthorId].FollowCount,
				FollowerCount:   userMap[v.AuthorId].FollowerCount,
				IsFollow:        userMap[v.AuthorId].IsFollow,
				BackgroundImage: userMap[v.AuthorId].BackgroundImage,
				Signature:       userMap[v.AuthorId].Signature,
				TotalFavorited:  userMap[v.AuthorId].TotalFavorited,
				FavoriteCount:   userMap[v.AuthorId].FavoriteCount,
			},
		})
	}

	if err != nil {
		response.Fail(ctx, err.Error(), nil)
		return
	}

	response.Success(ctx, "success", &DouyinFavoriteListResponse{
		VideoList: videoList,
	})
}
