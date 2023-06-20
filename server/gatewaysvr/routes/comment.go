package routes

import (
	"gatewaysvr/controller"
	"gatewaysvr/utils/common"
	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.RouterGroup) {
	comment := r.Group("comment")
	{
		comment.POST("/action/", common.AuthMiddleware(), controller.CommentAction)
		comment.GET("/list/", common.AuthWithOutMiddleware(), controller.GetCommentList)
	}

}
