package service

import (
	"favoritesvr/log"
	"favoritesvr/repository"

	"github.com/Percygu/camp_tiktok/pkg/pb"
)

func FavoriteAction(uid, vid int64, actionType int64) error {
	if actionType == 1 {
		log.Infof("like action uid:%v,vid:%v", uid, vid)
		err := repository.LikeAction(uid, vid)
		if err != nil {
			return err
		}
	} else {
		log.Infof("unlike action uid:%v,vid:%v", uid, vid)
		err := repository.UnLikeAction(uid, vid)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFavoriteVideoList(tokenUid, uid int64) (*pb.GetFavoriteListRsp, error) {
	favList, err := repository.GetFavoriteList(uid)
	if err != nil {
		return nil, err
	}
	// log.Infof("user:%v, followList:%+v", uid, favList)

	favListResponse := message.DouyinFavoriteListResponse{
		VideoList: VideoList(favList, tokenUid),
	}

	return &favListResponse, nil
}

func tokenFavList(tokenUserId int64) (map[int64]struct{}, error) {
	m := make(map[int64]struct{})
	list, err := repository.GetFavoriteList(tokenUserId)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		m[v.Id] = struct{}{}
	}
	return m, nil
}
