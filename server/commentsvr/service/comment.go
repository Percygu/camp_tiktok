package service

import (
	"commentsvr/constant"
	"commentsvr/log"
	"commentsvr/repository"
	"commentsvr/utils"
	"context"

	"github.com/Percygu/camp_tiktok/pkg/pb"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
}

func (c CommentService) CommentAction(ctx context.Context, req *pb.CommentActionReq) (*pb.CommentActionRsp, error) {
	// 增加评论
	if req.ActionType == 1 {
		comment, err := repository.CommentAdd(req.UserId, req.VideoId, req.CommentText)
		if err != nil {
			log.Errorf("CommentAction|CommentAdd err:%v", err)
			return nil, err
		}
		getUserInfoRsp, err := utils.GetUserSvrClient().GetUserInfo(ctx, &pb.GetUserInfoRequest{
			Id: req.UserId,
		})
		if err != nil {
			log.Errorf("CommentAction|GetUserInfo err %v", err)
			return nil, err
		}
		result := &pb.CommentActionRsp{Comment: &pb.Comment{
			Id:         comment.Id,
			User:       getUserInfoRsp.UserInfo,
			Content:    comment.CommentText,
			CreateDate: comment.CreateTime.Format(constant.DefaultTime),
		}}
		return result, nil

	} else {
		// 删除评论
		err := repository.CommentDelete(req.VideoId, req.CommentId)
		if err != nil {
			return nil, err
		}
		return &pb.CommentActionRsp{
			Comment: nil,
		}, nil
	}
}

func (c CommentService) GetCommentList(ctx context.Context, req *pb.GetCommentListReq) (*pb.GetCommentListRsp, error) {
	comments, err := repository.CommentList(req.VideoId)
	if err != nil {
		log.Errorf("GetCommentList|CommentList err:%v", err)
		return nil, err
	}

	userIDList := make([]int64, len(comments))
	for i, comment := range comments {
		userIDList[i] = comment.UserId
	}

	userInfoListRsp, err := utils.GetUserSvrClient().GetUserInfoList(context.Background(), &pb.GetUserInfoListRequest{
		IdList: userIDList,
	})
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
			User:       userInfo,
			Content:    comment.CommentText,
			CreateDate: comment.CreateTime.Format(constant.DefaultTime),
		}
		list.CommentList[i] = v
		log.Infof("commentsvr|comment====%+v", v)
	}
	return list, nil
}
