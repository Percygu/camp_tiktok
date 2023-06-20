package service

import (
	"commentsvr/constant"
	"commentsvr/repository"
	"context"
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
	//TODO implement me
	panic("implement me")
}
