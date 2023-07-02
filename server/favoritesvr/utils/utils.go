package utils

import (
	"favoritesvr/config"
	"favoritesvr/log"
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"google.golang.org/grpc"
	// 必须要导入这个包，否则grpc会报错
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc/credentials/insecure"
)

var (
	VideoSvrClient pb.VideoServiceClient
)

func NewSvrConn(svrName string) (*grpc.ClientConn, error) {
	consulInfo := config.GetGlobalConfig().ConsulConfig
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, svrName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		log.Errorf("NewSvrConn with svrname %s err:%v", svrName, err)
		return nil, err
	}
	return conn, nil
}

func NewVideoSvrClient(svrName string) pb.VideoServiceClient {
	conn, err := NewSvrConn(svrName)
	if err != nil {
		return nil
	}
	return pb.NewVideoServiceClient(conn)
}

func GetVideoSvrClient() pb.VideoServiceClient {
	return VideoSvrClient
}

func InitSvrConn() {
	VideoSvrClient = NewVideoSvrClient(config.GetGlobalConfig().SvrConfig.VideoSvrName)
}
