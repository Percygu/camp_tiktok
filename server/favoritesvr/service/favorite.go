package service

import (
	"context"
	"favoritesvr/constant"
	"favoritesvr/log"
	"favoritesvr/repository"
	"github.com/Percygu/camp_tiktok/pkg/pb"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
}

func (f *FavoriteService) FavoriteAction(ctx context.Context, req *pb.FavoriteActionReq) (*pb.FavoriteActionRsp, error) {
	if req.ActionType == 1 {
		log.Infof("like action uid:%v,vid:%v", req.UserId, req.VideoId)
		err := repository.LikeAction(req.UserId, req.VideoId)
		if err != nil {
			return nil, err
		}
	} else {
		log.Infof("unlike action uid:%v,vid:%v", req.UserId, req.VideoId)
		err := repository.UnLikeAction(req.UserId, req.VideoId)
		if err != nil {
			return nil, err
		}
	}
	return &pb.FavoriteActionRsp{
		CommonRsp: &pb.CommonResponse{
			Code: constant.SuccessCode,
			Msg:  constant.SuccessMsg,
		},
	}, nil
}

func (f *FavoriteService) GetFavoriteVideoIdList(ctx context.Context, req *pb.GetFavoriteVideoIdListReq) (*pb.GetFavoriteVideoIdListRsp, error) {
	videoIdList, err := repository.GetFavoriteIdList(req.UserId)
	if err != nil {
		return nil, err
	}
	log.Infof("get favorite video id list success", videoIdList)

	return &pb.GetFavoriteVideoIdListRsp{
		VideoIdList: videoIdList,
	}, nil
}

// func tokenFavList(tokenUserId int64) (map[int64]struct{}, error) {
// 	m := make(map[int64]struct{})
// 	list, err := repository.GetFavoriteList(tokenUserId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, v := range list {
// 		m[v.Id] = struct{}{}
// 	}
// 	return m, nil
// }
