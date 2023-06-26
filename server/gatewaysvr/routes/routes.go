package routes

import (
	"gatewaysvr/config"
	"gatewaysvr/controller"
	"gatewaysvr/utils/common"
	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	if config.GetGlobalConfig().SvrConfig.Mode == gin.ReleaseMode {
		// gin设置成发布模式：gin不在终端输出日志
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	douyin := r.Group("/douyin")
	{
		UserRoutes(douyin)
		PublishRoutes(douyin)
		CommentRoutes(douyin)
		FavoriteRoutes(douyin)
		RelationRoutes(douyin)
		douyin.GET("/feed/", common.AuthWithOutMiddleware(), controller.Feed)
	}

	return r
}
