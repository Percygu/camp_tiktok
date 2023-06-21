package routes

import (
	"gatewaysvr/controller"
	"gatewaysvr/utils/common"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	user := r.Group("user")
	{
		user.POST("/login/", controller.UserLogin)
		user.GET("/", common.AuthMiddleware(), controller.GetUserInfo)

		user.POST("/register/", controller.UserRegister)
	}

}
