package service

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"usersvr/log"
	"usersvr/middleware/lock"
	"usersvr/repository"

	"github.com/Percygu/camp_tiktok/pkg/pb"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var (
	Secret = []byte("TikTok")
	// TokenExpireDuration = time.Hour * 2 过期时间
)

type JWTClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

type UserService struct {
	pb.UnimplementedUserServiceServer
}

// UpdateUserFavoritedCount 更新用户 获赞数，ActionType 1：表示+1 2：-1
func (u UserService) UpdateUserFavoritedCount(ctx context.Context, req *pb.UpdateUserFavoritedCountReq) (*pb.UpdateUserFavoritedCountRsp, error) {
	err := repository.UpdateUserFavoritedNum(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateUserFavoritedCount err", req.UserId)
		return nil, err
	}
	return &pb.UpdateUserFavoritedCountRsp{}, nil
}

// UpdateUserFollowCount 更新用户 喜爱的视频数，ActionType 1：表示+1 2：-1
func (u UserService) UpdateUserFavoriteCount(ctx context.Context, req *pb.UpdateUserFavoriteCountReq) (*pb.UpdateUserFavoriteCountRsp, error) {
	err := repository.UpdateUserFavoriteNum(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateUserFavoriteCount err", req.UserId)
		return nil, err
	}
	return &pb.UpdateUserFavoriteCountRsp{}, nil
}

// UpdateUserFollowCount 更新用户 关注数，ActionType 1：表示+1 2：-1
func (u UserService) UpdateUserFollowCount(ctx context.Context, req *pb.UpdateUserFollowCountReq) (*pb.UpdateUserFollowCountRsp, error) {
	err := repository.UpdateUserFollowNum(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateUserFollowCount err", req.UserId)
		return nil, err
	}
	return &pb.UpdateUserFollowCountRsp{}, nil
}

// UpdateUserFollowerCount 更新用户 粉丝数，ActionType 1：表示+1 2：-1
func (u UserService) UpdateUserFollowerCount(ctx context.Context, req *pb.UpdateUserFollowerCountReq) (*pb.UpdateUserFollowerCountRsp, error) {
	err := repository.UpdateUserFollowerNum(req.UserId, req.ActionType)
	if err != nil {
		log.Errorf("UpdateUserFollowerCount err", req.UserId)
		return nil, err
	}
	return &pb.UpdateUserFollowerCountRsp{}, nil
}

func (u UserService) GetUserInfoDict(ctx context.Context, req *pb.GetUserInfoDictRequest) (*pb.GetUserInfoDictResponse, error) {
	userList, err := repository.GetUserList(req.UserIdList)
	if err != nil {
		log.Errorf("GetUserInfoDict err", req.UserIdList)
		return nil, err
	}
	resp := &pb.GetUserInfoDictResponse{UserInfoDict: make(map[int64]*pb.UserInfo)}

	for _, user := range userList {
		resp.UserInfoDict[user.Id] = &pb.UserInfo{
			Id:              user.Id,
			Name:            user.Name,
			Avatar:          user.Avatar,
			FollowCount:     user.Follow,
			FollowerCount:   user.Follower,
			BackgroundImage: user.BackgroundImage,
			Signature:       user.Signature,
			TotalFavorited:  user.TotalFav,
			FavoriteCount:   user.FavCount,
		}
	}

	return resp, nil
}

func (u UserService) CacheChangeUserCount(ctx context.Context, req *pb.CacheChangeUserCountReq) (*pb.CacheChangeUserCountRsp, error) {
	uid := strconv.FormatInt(req.UserId, 10)
	mutex := lock.GetLock("user_" + uid)
	defer lock.UnLock(mutex)
	user, err := repository.CacheGetUser(req.UserId)
	if err != nil {
		log.Infof("CacheChangeUserCount err", req.UserId)
		return nil, err
	}

	switch req.CountType {
	case "follow":
		user.Follow += req.Op
	case "follower":
		user.Follower += req.Op
	case "like":
		user.FavCount += req.Op
	case "liked":
		user.TotalFav += req.Op
	}
	repository.CacheSetUser(user)

	return &pb.CacheChangeUserCountRsp{}, nil
}

// func (u UserService) CacheGetAuthor(ctx context.Context, req *pb.CacheGetAuthorReq) (*pb.CacheGetAuthorRsp, error) {
// 	key := strconv.FormatInt(req.VideoId, 10)
// 	data, err := repository.CacheHGet("video", key)
// 	if err != nil {
// 		log.Errorf("CacheGetAuthor err", req.VideoId)
// 		return nil, err
// 	}
//
// 	uid := int64(0)
// 	err = json.Unmarshal(data, &uid)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &pb.CacheGetAuthorRsp{UserId: uid}, nil
// }

// GetUserInfoList 获取用户信息列表
func (u UserService) GetUserInfoList(ctx context.Context, req *pb.GetUserInfoListRequest) (*pb.GetUserInfoListResponse, error) {
	response := new(pb.GetUserInfoListResponse)

	userList, err := repository.GetUserList(req.IdList)
	if err != nil {
		zap.L().Error("GetUserList error", zap.Error(err))
		return nil, err
	}

	for _, user := range userList {
		response.UserInfoList = append(response.UserInfoList, UserToUserInfo(*user))
	}

	return response, nil
}

// GetUserInfo 获取用户信息
func (u UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	user, err := repository.GetUserInfo(req.Id)
	if err != nil {
		return nil, err
	}
	response := &pb.GetUserInfoResponse{
		UserInfo: UserToUserInfo(user),
	}

	return response, nil
}

// CheckPassWord 登录验证
func (u UserService) CheckPassWord(ctx context.Context, req *pb.CheckPassWordRequest) (*pb.CheckPassWordResponse, error) {
	info, err := repository.GetUserInfo(req.Username)
	if err != nil {
		return nil, err
	}
	// 验证密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(info.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("password error")
	}
	token, err := GenToken(info.Id, req.Username)
	if err != nil {
		return nil, err
	}
	response := &pb.CheckPassWordResponse{
		UserId: info.Id,
		Token:  token,
	}
	return response, nil
}

// Register 注册
func (u UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	sign, err := repository.UserNameIsExist(req.Username)
	if err != nil {
		log.Error("UserNameIsExist err ", err)
		return nil, err
	}
	if sign {
		return nil, fmt.Errorf("user %s exists", req.Username)
	}
	info, err := repository.InsertUser(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := GenToken(info.Id, req.Username)
	if err != nil {
		return nil, err
	}
	registerResponse := &pb.RegisterResponse{
		UserId: info.Id,
		Token:  token,
	}

	return registerResponse, nil
}

func UserToUserInfo(info repository.User) *pb.UserInfo {
	return &pb.UserInfo{
		Id:              info.Id,
		Name:            info.Name,
		FollowCount:     info.Follow,
		FollowerCount:   info.Follower,
		IsFollow:        false,
		Avatar:          info.Avatar,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorited:  info.TotalFav,
		FavoriteCount:   info.FavCount,
	}
}

// 生成token
func GenToken(userid int64, userName string) (string, error) {
	claims := JWTClaims{
		UserId:   userid,
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "server",
			// ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),可用于设定token过期时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("TikTok"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
