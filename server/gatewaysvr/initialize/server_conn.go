package initialize

import (
	"errors"
	"fmt"
	"gatewaysvr/global"
	"gatewaysvr/utils/otgrpc"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() (err error) {
	consulInfo := global.Conf.ConsulConfig

	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port,
			global.Conf.UserServerConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)

	if err != nil {
		return errors.New("连接用户服务失败")
	}
	// global.UserSrvClient = proto.NewUserClient(userConn)

	// global.VideoSrvClient = proto.NewUserClient(userConn)
	// global.CommentSrvClient = proto.NewUserClient(userConn)
	// global.FollowSrvClient = proto.NewUserClient(userConn)
	// global.UserSrvClient = proto.NewUserClient(userConn)

	return nil
}
