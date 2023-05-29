package routes

import (
	"dou_yin/controllers"
	"dou_yin/logger"
	"dou_yin/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine{
	r:=gin.New()
	r.Use(middleware.Cors())
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ws/:token",controllers.Connect)
	r.POST("/register",controllers.Register)
	r.POST("/login",controllers.Login)
	r.GET("/test",controllers.TestRedis)
	g:=r.Group("/api",middleware.JWTAuthMiddleware())
	g.GET("/get_contactor_list/:id",controllers.GetContactorList)
	return r
}
