package utils

import (
	"context"
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"math/rand"
	"os"
	"time"

	// 必须要导入这个包，否则grpc会报错
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
)

var (
	VideoSvrClient    pb.VideoServiceClient
	UserSvrClient     pb.UserServiceClient
	CommentSvrClient  pb.CommentServiceClient
	RelationSvrClient pb.RelationServiceClient
	FavoriteSvrClient pb.FavoriteServiceClient
)

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		// grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	log.Info("NewSvrConn success")
	return conn, nil
}

func GetVideoSvrClient() pb.VideoServiceClient {
	return VideoSvrClient
}

func GetUserSvrClient() pb.UserServiceClient {
	return UserSvrClient
}

func GetCommentSvrClient() pb.CommentServiceClient {
	return CommentSvrClient
}

func GetRelationSvrClient() pb.RelationServiceClient {
	return RelationSvrClient
}

func GetFavoriteSvrClient() pb.FavoriteServiceClient {
	return FavoriteSvrClient
}

func NewVideoSvrClient(svrName string) pb.VideoServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewVideoServiceClient(conn)
}

func NewUserSvrClient(svrName string) pb.UserServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewUserServiceClient(conn)
}

func NewCommentSvrClient(svrName string) pb.CommentServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewCommentServiceClient(conn)
}

func NewRelationSvrClient(svrName string) pb.RelationServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewRelationServiceClient(conn)
}

func NewFavoriteSvrClient(svrName string) pb.FavoriteServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewFavoriteServiceClient(conn)
}

func InitSvrConn() {
	VideoSvrClient = NewVideoSvrClient(config.GetGlobalConfig().SvrConfig.VideoSvrName)
	UserSvrClient = NewUserSvrClient(config.GetGlobalConfig().SvrConfig.UserSvrName)
	CommentSvrClient = NewCommentSvrClient(config.GetGlobalConfig().SvrConfig.CommentSvrName)
	RelationSvrClient = NewRelationSvrClient(config.GetGlobalConfig().SvrConfig.RelationSvrName)
	FavoriteSvrClient = NewFavoriteSvrClient(config.GetGlobalConfig().SvrConfig.FavoriteSvrName)
}

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}

// 随机生成字符
func RandomString() string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, 16)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
