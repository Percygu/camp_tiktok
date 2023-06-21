package routes

import (
	"gatewaysvr/controller"
	"gatewaysvr/utils/common"
	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	r := gin.Default()
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
