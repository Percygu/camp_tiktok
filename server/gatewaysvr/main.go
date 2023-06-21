package main

import (
	"fmt"
	"gatewaysvr/global"
	"gatewaysvr/initialize"
	"gatewaysvr/routes"
	"go.uber.org/zap"
	"log"
)

func main() {
	// 1.加载配置
	if err := initialize.InitConfig(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	// 2.初始化日志
	if err := initialize.InitLogger(); err != nil {
		log.Fatalf("init logger failed, err:%v\n", err)
	}
	defer zap.L().Sync() // 把缓存区的日志追加到文件中
	zap.L().Info("init config success...")
	zap.L().Info("init logger success...")

	// 4. 初始化server的连接
	if err := initialize.InitSrvConn(); err != nil {
		zap.L().Sugar().Fatalf("SrvConn initialization error: %v", err)
	}

	// 5.注册路由
	r := routes.SetRoute()
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", global.Conf.Port)); err != nil {
			zap.L().Panic("Router.Run error: ", zap.Error(err))
		}
	}()
	zap.L().Sugar().Infof("listen on %s:%d", global.Conf.Host, global.Conf.Port)

}
