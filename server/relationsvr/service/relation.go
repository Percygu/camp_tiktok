package service

import (
	"context"
	"fmt"
	"relationsvr/constant"
	"relationsvr/repository"

	"github.com/Percygu/camp_tiktok/pkg/pb"
	"relationsvr/log"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

func (c RelationService) CommentAction(ctx context.Context, req *pb.RelationActionReq) (*pb.RelationActionRsp, error) {
	if req.SelfUserId == req.ToUserId {
		return nil, fmt.Errorf("you can't follow yourself")
	}
	if req.ActionType == 1 {
		log.Infof("follow action id:%v,toid:%v", req.SelfUserId, req.ToUserId)
		err := repository.FollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, err
		}
	} else {
		log.Infof("unfollow action id:%v,toid:%v", req.SelfUserId, req.ToUserId)
		err := repository.UnFollowAction(req.SelfUserId, req.ToUserId)
		if err != nil {
			return nil, err
		}
	}
	return &pb.RelationActionRsp{
		CommonRsp: &pb.CommonResponse{
			Code: constant.SuccessCode,
			Msg:  constant.SuccessMsg,
		},
	}, nil
}

// GetRelationFollowList 获取被关注者列表
func GetRelationFollowList(ctx context.Context, userId int64) (*pb.GetRelationFollowListRsp, error) {
	userInfoList, err := RelationFollowList(userId, 1)
	if err != nil {
		return nil, err
	}
	return &pb.GetRelationFollowListRsp{
		CommonRsp: &pb.CommonResponse{
			Code: constant.SuccessCode,
			Msg:  constant.SuccessMsg,
		},
		UserInfo: userInfoList,
	}, nil
}

// GetRelationFollowerList 获取关注者列表
func GetRelationFollowerList(userId int64, tokenUserId int64) (*pb.GetRelationFollowerListRsp, error) {
	userInfoList, err := RelationFollowList(userId, 2)
	if err != nil {
		return nil, err
	}
	return &pb.GetRelationFollowerListRsp{
		CommonRsp: &pb.CommonResponse{
			Code: constant.SuccessCode,
			Msg:  constant.SuccessMsg,
		},
		UserInfo: userInfoList,
	}, nil
}

func RelationFollowList(userID, relationType int64) ([]*pb.UserInfo, error) {
	var (
		relationList []*repository.Relation
		err          error
	)
	if relationType == 1 {
		// 获取关注者
		relationList, err = repository.GetFollowList(userID)
	} else {
		// 获取被关注者
		relationList, err = repository.GetFollowerList(userID)
	}
	if err != nil {
		return nil, err
	}
	if len(relationList) == 0 {
		return []*pb.UserInfo{}, nil
	}
	log.Infof("user:%v, relationList:%+v", userID, relationList)
	userIdList := make([]int64, 0)
	for _, relation := range relationList {
		userIdList = append(userIdList, relation.Follow)
	}
	// todo 获取用户列表
	return []*pb.UserInfo{}, nil
}
