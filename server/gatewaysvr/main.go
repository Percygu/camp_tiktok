package main

import (
	"fmt"
	"gatewaysvr/config"
	"gatewaysvr/log"
	"gatewaysvr/routes"
	"go.uber.org/zap"
)

func main() {
	// 1.加载配置
	config.Init()
	// 2.初始化日志
	log.InitLog()

	defer zap.L().Sync() // 把缓存区的日志追加到文件中
	zap.L().Info("init config success...")
	zap.L().Info("init logger success...")

	// 3.初始化路由
	r := routes.SetRoute()
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", config.GetGlobalConfig().Port)); err != nil {
			zap.L().Panic("Router.Run error: ", zap.Error(err))
		}
	}()
	zap.L().Sugar().Infof("listen on %s:%d", config.GetGlobalConfig().Host, config.GetGlobalConfig().Port)

}
