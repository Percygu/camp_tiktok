package initialize

import (
	"errors"
	"fmt"
	"gatewaysvr/global"
	"gatewaysvr/utils/otgrpc"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSrvConn() (err error) {
	consulConfig := global.Conf.ConsulConfig

	// 1.初始化用户服务连接
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port,
			global.Conf.UserServerConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return errors.New("连接User服务失败")
	}
	global.UserSrvClient = pb.NewUserServiceClient(userConn)

	// 2.初始化视频服务连接
	videoConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port,
			global.Conf.UserServerConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return errors.New("连接Video服务失败")
	}
	global.VideoSrvClient = pb.NewVideoServiceClient(videoConn)

	// 3. 初始化评论服务连接
	commentConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port,
			global.Conf.UserServerConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return errors.New("连接Comment服务失败")
	}
	global.CommentSrvClient = pb.NewCommentServiceClient(commentConn)

	relationConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulConfig.Host, consulConfig.Port,
			global.Conf.UserServerConfig.Name),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return errors.New("连接Relation服务失败")
	}
	global.RelationSrvClient = pb.NewRelationServiceClient(relationConn)

	// TODO： 初始化其他服务连接

	return nil
}
