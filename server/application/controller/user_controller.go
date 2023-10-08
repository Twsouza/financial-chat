package controller

import (
	"net/http"
	"server/application/dto"
	"server/application/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(us *services.UserService) UserController {
	return UserController{
		UserService: us,
	}
}

func (uc *UserController) SignUp(gc *gin.Context) {
	req := dto.CreateUserReq{}
	if err := gc.Bind(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.UserService.CreateUser(gc.Request.Context(), &req)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gc.JSON(http.StatusCreated, res)
}

func (uc UserController) Login(gc *gin.Context) {
	req := dto.LoginUserReq{}
	if err := gc.Bind(&req); err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := uc.UserService.Login(gc.Request.Context(), &req)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gc.SetCookie("token", res.AccessToken, 3600, "/", "localhost", false, true)
	gc.SetCookie("userId", res.ID, 3600, "/", "localhost", false, true)
	gc.SetCookie("username", res.Username, 3600, "/", "localhost", false, true)

	gc.JSON(http.StatusOK, res)
}
