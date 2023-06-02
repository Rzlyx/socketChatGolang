package routes

import (
	"dou_yin/controllers/chat"
	"dou_yin/controllers/friend"
	"dou_yin/controllers/user"
	"dou_yin/logger"
	"dou_yin/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ws/:token", chat.Connect)
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	r.GET("/test", user.TestRedis)
	g := r.Group("/api", middleware.JWTAuthMiddleware())
	g.GET("/get_contactor_list/:id", user.GetContactorList)

	// friend
	r.POST("/queryFriendList", friend.QueryFriendList)
	r.POST("/queryFriendInfo", friend.QueryFriendInfo)
	r.POST("/addFriend", friend.AddFriend)
	r.POST("/deleteFriend", friend.DeleteFriend)
	r.POST("/setPrivateChatBlack", friend.SetPrivateChatBlack)
	r.POST("/unBlockPrivateChat", friend.UnBlockPrivateChat)
	r.POST("/setFriendCircleBlack", friend.SetFriendCircleBlack)
	r.POST("/unBlockFriendCircle", friend.UnBlockFriendCircle)
	r.POST("/queryFriendApply", friend.QueryFriendApply)
	r.POST("/agreeFriendApply", friend.AgreeFriendApply)
	r.POST("/disagreeFriendApply", friend.DisagreeFriendApply)
	r.POST("/setFriendRemark", friend.SetFriendRemark)

	return r
}
