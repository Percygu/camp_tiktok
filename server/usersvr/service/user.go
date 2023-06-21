package service

import (
	"context"
	"errors"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"usersvr/repository"
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
	pb.UnimplementedCommentServiceServer
}

func (u UserService) GetUserInfoList(ctx context.Context, request *pb.GetUserInfoListRequest) (response *pb.GetUserInfoListResponse, err error) {
	for _, user := range request.IdList {
		info, err := repository.GetUserInfo(user)
		if err != nil {
			return nil, err
		}
		// TODO: 没写
		response.UserInfoList = append(response.UserInfoList, info)
	}

	return response, nil
}

func (u UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
	user, err := repository.GetUserInfo(req.Id)
	if err != nil {
		return nil, err
	}

	return UserInfoToResponse(user), nil
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
	err := repository.UserNameIsExist(req.Username)
	if err != nil {
		return nil, err
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
		UserID: info.Id,
		Token:  token,
	}

	return registerResponse, nil
}

func UserInfoToResponse(info repository.User) *pb.GetUserInfoResponse {
	return &pb.GetUserInfoResponse{
		Id:              info.Id,
		Name:            info.Name,
		FollowCount:     info.Follow,
		FollowerCount:   info.Follower,
		IsFollow:        false,
		Avatar:          info.Avatar,
		BackgroundImage: info.BackgroundImage,
		Signature:       info.Signature,
		TotalFavorite:   info.TotalFav,
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

// 解析token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// 验证token
func VerifyToken(tokenString string) (int64, error) {
	zap.L().Debug("tokenString", zap.String("tokenString", tokenString))

	if tokenString == "" {
		return int64(0), nil
	}
	claims, err := ParseToken(tokenString)
	if err != nil {
		return int64(0), err
	}

	return claims.UserId, nil
}
