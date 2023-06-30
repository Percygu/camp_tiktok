package service

import (
	"context"
	"fmt"
	"relationsvr/constant"
	"relationsvr/repository"
	"strconv"

	"github.com/Percygu/camp_tiktok/pkg/pb"
	"relationsvr/log"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
}

func (c RelationService) IsFollowDict(ctx context.Context, req *pb.IsFollowDictReq) (*pb.IsFollowDictRsp, error) {
	var isFollowDict = make(map[string]bool)

	for _, unit := range req.FollowUintList {
		// TODO: UserIdList 可能是我关注的人
		isFollow, err := repository.IsFollow(unit.SelfUserId, unit.UserIdList)
		if err != nil {
			log.Errorf("IsFollowDict err", unit.SelfUserId, unit.UserIdList)
			return nil, err
		}
		isFollowKey := strconv.FormatInt(unit.SelfUserId, 10) + "_" + strconv.FormatInt(unit.UserIdList, 10)
		isFollowDict[isFollowKey] = isFollow
	}

	return &pb.IsFollowDictRsp{IsFollowDict: isFollowDict}, nil
}

func (c RelationService) RelationAction(ctx context.Context, req *pb.RelationActionReq) (*pb.RelationActionRsp, error) {
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
func (c RelationService) GetRelationFollowList(ctx context.Context, req *pb.GetRelationFollowListReq) (*pb.GetRelationFollowListRsp, error) {
	userInfoList, err := RelationFollowList(req.UserId, 1)
	if err != nil {
		return nil, err
	}
	return &pb.GetRelationFollowListRsp{
		UserInfoList: userInfoList,
	}, nil
}

// GetRelationFollowerList 获取关注者列表
func (c RelationService) GetRelationFollowerList(ctx context.Context, req *pb.GetRelationFollowerListReq) (*pb.GetRelationFollowerListRsp, error) {
	userInfoList, err := RelationFollowList(req.UserId, 2)
	if err != nil {
		return nil, err
	}
	return &pb.GetRelationFollowerListRsp{
		UserInfoList: userInfoList,
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
