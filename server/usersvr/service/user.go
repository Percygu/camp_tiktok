package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (u UserService) CacheGetAuthor(ctx context.Context, req *pb.CacheGetAuthorReq) (*pb.CacheGetAuthorRsp, error) {
	key := strconv.FormatInt(req.VideoId, 10)
	data, err := repository.CacheHGet("video", key)
	if err != nil {
		log.Errorf("CacheGetAuthor err", req.VideoId)
		return nil, err
	}

	uid := int64(0)
	err = json.Unmarshal(data, &uid)
	if err != nil {
		return nil, err
	}

	return &pb.CacheGetAuthorRsp{UserId: uid}, nil
}

func (u UserService) GetUserInfoList(ctx context.Context, request *pb.GetUserInfoListRequest) (response *pb.GetUserInfoListResponse, err error) {
	for _, user := range request.IdList {
		info, err := repository.GetUserInfo(user)
		if err != nil {
			return nil, err
		}
		response.UserInfoList = append(response.UserInfoList, UserToUserInfo(info))
	}

	return response, nil
}

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
