package service

import (
	"commentsvr/constant"
	"commentsvr/log"
	"commentsvr/repository"
	"commentsvr/utils"
	"context"
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
}

func (c CommentService) CommentAction(ctx context.Context, req *pb.CommentActionReq) (*pb.CommentActionRsp, error) {
	// 增加评论
	if req.ActionType == 1 {
		// commentInfo, err := repository.CommentAdd(userId, videoId, comment_text)
		_, err := repository.CommentAdd(req.UserId, req.VideoId, req.CommentText)
		if err != nil {
			return nil, err
		}
		return &pb.CommentActionRsp{
			CommonRsp: &pb.CommonResponse{
				Code: constant.SuccessCode,
				Msg:  constant.SuccessMsg,
			},
		}, nil
	} else { // 删除评论
		err := repository.CommentDelete(req.VideoId, req.CommentId)
		if err != nil {
			return nil, err
		}
		return &pb.CommentActionRsp{
			CommonRsp: &pb.CommonResponse{
				Code: constant.SuccessCode,
				Msg:  constant.SuccessMsg,
			},
		}, nil
	}
}

func (c CommentService) GetCommentList(ctx context.Context, req *pb.GetCommentListReq) (*pb.GetCommentListRsp, error) {
	comments, err := repository.CommentList(req.VideoId)
	if err != nil {
		return nil, err
	}
	log.Infof("comments:%v\n", comments)

	userIDList := make([]int64, len(comments))
	for _, comment := range comments {
		userIDList = append(userIDList, comment.UserId)
	}

	userSvrClient := utils.GetUserSvrClient()
	if userSvrClient == nil {
		return nil, fmt.Errorf("userSvrClient is nil")
	}
	userInfoListReq := &pb.GetUserInfoListRequest{
		IdList: userIDList,
	}
	userInfoListRsp, err := userSvrClient.GetUserInfoList(context.Background(), userInfoListReq)
	if err != nil {
		log.Errorf("GetCommentList|GetUserInfoList err:%v", err)
		return nil, err
	}
	uerMap := make(map[int64]*pb.UserInfo)
	for _, userInfo := range userInfoListRsp.UserInfoList {
		uerMap[userInfo.Id] = userInfo
	}

	list := &pb.GetCommentListRsp{
		CommentList: make([]*pb.Comment, len(comments)),
	}

	for i, comment := range comments {
		// 为了找到video_id所对应的user_id，在通过user_id找到user_name.传递给前端
		userInfo := uerMap[comment.UserId]
		v := &pb.Comment{
			Id:         comment.Id,
			UserInfo:   userInfo,
			Content:    comment.CommentText,
			CreateDate: fmt.Sprint(comment.CreateTime),
		}
		list.CommentList[i] = v
	}
	return list, nil
}
