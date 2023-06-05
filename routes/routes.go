package routes

import (
	"dou_yin/controllers/friend"
	"dou_yin/controllers/group"
	"dou_yin/controllers/privateChat"
	"dou_yin/controllers/user"
	"dou_yin/logger"
	"dou_yin/middleware"
	"dou_yin/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.Cors())
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/ws/:token", service.Connect)

	// user
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	r.GET("/test", user.TestRedis)
	g := r.Group("/api", middleware.JWTAuthMiddleware())
	g.GET("/get_contactor_list/:id", user.GetContactorList)
	r.POST("/queryContactorList", user.QueryContactorList)
	r.GET("/getPhotoByID/:id", user.GetPhotoByID)
	r.POST("/uploadPhoto", user.UploadPhoto)
	r.POST("/setContactorList", user.SetContactorList)
	r.POST("/updateUserInfo", user.UpdateUserInfo)

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
	r.POST("/setReadTime", friend.SetReadTime)

	// private privateChat
	r.POST("/queryPrivateChatMsg", privateChat.QueryPrivateChatMsg)
	r.POST("/deletePrivateChatMsg", privateChat.DeletePrivateChatMsg)
	r.POST("/uploadPrivateChatPhoto", privateChat.UploadPrivateChatPhoto)
	r.POST("/uploadPrivateChatFile", privateChat.UploadPrivateChatFile)
	r.GET("/getFileByID/:id", privateChat.GetFileByID)

	// group
	r.POST("/CreateGroupInfo", group.CreateGroupInfo)
	r.POST("/QueryGroupInfo", group.QueryGroupInfo)
	r.POST("/QueryGroupList", group.QueryGroupList)
	r.POST("/DissolveGroupInfo", group.DissolveGroupInfo)
	r.POST("/ApplyJoinGroup", group.ApplyJoinGroup)
	r.POST("/QuitGroup", group.QuitGroup)
	r.POST("/QueryGroupApplyList", group.QueryGroupApplyList)
	r.POST("/AgreeGroupApply", group.AgreeGroupApply)
	r.POST("/DisAgreeGroupApply", group.DisAgreeGroupApply)
	r.POST("/Silence", group.Silence)
	r.POST("/UnSilence", group.UnSilence)
	r.POST("/TransferGroup", group.TransferGroup)
	r.POST("/SetBlackList", group.SetBlackList)
	r.POST("/SetGrayList", group.SetGrayList)
	r.POST("/SetWhiteList", group.SetWhiteList)
	r.POST("/SetGroupAdmin", group.SetGroupAdmin)
	r.POST("/SetGroupUser", group.SetGroupUser)
	r.POST("/KickUserFromGroup", group.KickUserFromGroup)

	return r
}
