package utils

import (
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"videosvr/config"
	"videosvr/log"
	"videosvr/utils/otgrpc"
)

var (
	FavoriteSvrClient pb.FavoriteServiceClient
	RelationSvrClient pb.RelationServiceClient
)

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	return conn, nil
}

// func NewUserSvrClient(svrName string) pb.UserServiceClient {
// 	conn, err := NewSvrConn(svrName)
// 	if err != nil {
// 		return nil
// 	}
// 	return pb.NewUserServiceClient(conn)
// }

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

func GetRelationSvrClient() pb.RelationServiceClient {
	return RelationSvrClient
}

func GetFavoriteSvrClient() pb.FavoriteServiceClient {
	return FavoriteSvrClient
}

func InitSvrConn() {
	RelationSvrClient = NewRelationSvrClient(config.GetGlobalConfig().SvrConfig.RelationSvrName)
	FavoriteSvrClient = NewFavoriteSvrClient(config.GetGlobalConfig().SvrConfig.FavoriteSvrName)
}
