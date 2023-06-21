package main

import (
	"fmt"
	"github.com/Percygu/camp_tiktok/pkg/pb"
	"relationsvr/config"
	"relationsvr/log"
	"relationsvr/middleware/consul"
	"relationsvr/service"
	// "github.com/Percygu/litetiktok_proto/pb"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	_ "google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	log.Info("log init success...")
}

func Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GetGlobalConfig().Host, config.GetGlobalConfig().Port))
	if err != nil {
		log.Fatalf("listen: error %v", err)
		return fmt.Errorf("listen: error %v", err)
	}
	// 端口监听启动成功，启动grpc server
	server := grpc.NewServer()
	// 注册grpc server
	pb.RegisterCommentServiceServer(server, &service.CommentService{}) // 注册服务
	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 注册服务到consul中
	consulClient := consul.NewRegistryClient(config.GetGlobalConfig().ConsulConfig.Host, config.GetGlobalConfig().ConsulConfig.Port)
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	if err := consulClient.Register(config.GetGlobalConfig().Host, config.GetGlobalConfig().Port,
		config.GetGlobalConfig().Name, config.GetGlobalConfig().ConsulConfig.Tags, serviceID); err != nil {
		log.Fatal("consul.Register error: ", zap.Error(err))
		return fmt.Errorf("consul.Register error: ", zap.Error(err))
	}
	log.Info("Init Consul Register success")

	// 启动
	log.Infof("TikTokLite.comment_svr listening on %s:%d", config.GetGlobalConfig().Host, config.GetGlobalConfig().Port)
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 服务终止，注销 consul 服务
	if err = consulClient.DeRegister(serviceID); err != nil {
		log.Info("注销失败")
		return fmt.Errorf("注销失败")
	} else {
		log.Info("注销成功")
	}
	return nil
}

func main() {
	Init()
	defer log.Sync()
	if err := Run(); err != nil {
		log.Errorf("commentsvr run err:%v", err)
	}
}
