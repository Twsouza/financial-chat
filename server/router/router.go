package router

import (
	"server/application/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter(uc controller.UserController) *gin.Engine {
	r := gin.Default()

	r.POST("/signup", uc.SignUp)
	r.POST("/login", uc.Login)

	return r
}
