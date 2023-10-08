package router

import (
	"net/http"
	"server/application/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(uc controller.UserController, ct controller.ChatController) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/signup", uc.SignUp)
	r.POST("/login", uc.Login)

	ws := r.Group("/ws")
	ws.Use(Auth)
	ws.POST("/rooms", ct.CreateRoom)
	ws.GET("/rooms", ct.GetRooms)
	ws.GET("/rooms/:id/clients", ct.GetClients)
	r.GET("/ws/rooms/:id/join", ct.JoinRoom)

	return r
}

func Auth(gc *gin.Context) {
	tokenCookie, err := gc.Cookie("token")
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
		return
	}
	userIdCookie, err := gc.Cookie("userId")
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
		return
	}
	usernameCookie, err := gc.Cookie("username")
	if err != nil {
		gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
		return
	}

	if tokenCookie == "" || userIdCookie == "" || usernameCookie == "" {
		gc.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "you must be logged in"})
		return
	}

	gc.Next()
}
