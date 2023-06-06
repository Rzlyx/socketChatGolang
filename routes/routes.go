package routes

import (
	"dou_yin/controllers/friend"
	"dou_yin/controllers/friendCircle"
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

	r.GET("/getPhotoByID/:id", user.GetPhotoByID)
	r.GET("/getFileByID/:id", privateChat.GetFileByID)

	// user
	r.POST("/register", user.Register)
	r.POST("/login", user.Login)
	r.GET("/test", user.TestRedis)
	g := r.Group("/api", middleware.JWTAuthMiddleware())
	g.GET("/get_contactor_list/:id", user.GetContactorList)
	r.POST("/queryContactorList", user.QueryContactorList)
	r.POST("/uploadPhoto", user.UploadPhoto)
	r.POST("/setContactorList", user.SetContactorList)
	r.POST("/updateUserInfo", user.UpdateUserInfo)
	r.POST("/queryUserInfo", user.QueryUserInfo)
	r.POST("/searchFriendOrGroup", user.SearchFriendOrGroup)

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
	r.POST("/setPrivateChatGray", friend.SetPrivateChatGray)
	r.POST("/unGrayPrivateChat", friend.UnGrayPrivateChat)
	r.POST("/addFriendTag", friend.AddFriendTag)
	r.POST("/removeFriendTag", friend.RemoveFriendTag)

	// privateChat
	r.POST("/queryPrivateChatMsg", privateChat.QueryPrivateChatMsg)
	r.POST("/queryPrivateChatMsgByDate", privateChat.QueryPrivateChatMsgByDate)
	r.POST("/queryPrivateChatMsgByReadTime", privateChat.QueryPrivateChatMsgByReadTime)
	r.POST("/deletePrivateChatMsg", privateChat.DeletePrivateChatMsg)
	r.POST("/uploadPrivateChatPhoto", privateChat.UploadPrivateChatPhoto)
	r.POST("/uploadPrivateChatFile", privateChat.UploadPrivateChatFile)

	// group
	r.POST("/CreateGroupInfo", group.CreateGroupInfo)           // 创建群聊
	r.POST("/UploadGroupPhoto", group.UploadGroupPhoto)         // 上传头像
	r.POST("/QueryGroupInfo", group.QueryGroupInfo)             // 查看群信息
	r.POST("/QueryGroupList", group.QueryGroupList)             // 查看已添加群
	r.POST("/GetGroupAllUser", group.GetGroupAllUser)           // 获取所有群成员
	r.POST("/UpdateGroupInfo", group.UpdateGroupInfo)           // 更新群信息
	r.POST("/DissolveGroupInfo", group.DissolveGroupInfo)       // 解散群
	r.POST("/ApplyJoinGroup", group.ApplyJoinGroup)             // 申请加群
	r.POST("/QuitGroup", group.QuitGroup)                       // 退群
	r.POST("/KickUserFromGroup", group.KickUserFromGroup)       // 踢人
	r.POST("/QueryGroupApplyList", group.QueryGroupApplyList)   // 查看加群申请（管理员/群主）
	r.POST("/AgreeGroupApply", group.AgreeGroupApply)           // 同意加群申请（管理员/群主）
	r.POST("/DisAgreeGroupApply", group.DisAgreeGroupApply)     // 拒绝加群申请（管理员/群主）
	r.POST("/Silence", group.Silence)                           // 禁言（管理员/群主）
	r.POST("/UnSilence", group.UnSilence)                       // 解除禁言（管理员/群主）
	r.POST("/TransferGroup", group.TransferGroup)               // 转让群（群主）
	r.POST("/SetBlackList", group.SetBlackList)                 // 屏蔽消息(将该群聊加入黑名单)
	r.POST("/SetGrayList", group.SetGrayList)                   // 免打扰消息(将该群聊设置为免打扰)
	r.POST("/SetWhiteList", group.SetWhiteList)                 // 通知消息(将该群聊设为可通知消息)
	r.POST("/SetGroupAdmin", group.SetGroupAdmin)               // 将普通成员设为管理员
	r.POST("/SetGroupUser", group.SetGroupUser)                 // 将管理员设为普通成员
	r.POST("/InviteJoinGroup", group.InviteJoinGroup)           // 邀请加入群聊
	r.POST("/QueryInviteGroup", group.QueryInviteGroup)         // 查看邀请
	r.POST("/AgreeInviteGroup", group.AgreeInviteGroup)         // 同意邀请
	r.POST("/DisAgreeInviteGroup", group.DisAgreeInviteGroup)   // 拒绝邀请
	r.POST("/SetGroupName", group.SetGroupName)                 // 设置群备注
	r.POST("/SetMyName", group.SetMyName)                       // 设置在本群的昵称
	r.POST("/SetGroupReadTime", group.SetGroupReadTime)         // 设置群聊已读时间
	r.POST("/GetPageOldMsg", group.GetPageOldMsg)               // 按页查询历史消息
	r.POST("/GetGroupOldMsgLogin", group.GetGroupOldMsgLogin)   // 登录获取历史消息--获取15条消息
	r.POST("/GetGroupOldMsgUp", group.GetGroupOldMsgUp)         // 加载向上的信息
	r.POST("/GetGroupOldMsgDay", group.GetGroupOldMsgDay)       // 按天数获取信息
	r.POST("/UploadGroupChatPhoto", group.UploadGroupChatPhoto) // 上传群相册
	r.POST("/UploadGroupChatFile", group.UploadGroupChatFile)   // 上传群文件
	// r.POST("/", group.)

	// friend_circle
	r.POST("/SendCircle", friendCircle.SendCircle)               // 发朋友圈
	r.POST("/UploadCirclePhoto", friendCircle.UploadCirclePhoto) // 上传图片
	r.POST("/queryAllFriendCircle", friendCircle.QueryAllFriendCircle)
	r.POST("/queryFriendCircle", friendCircle.QueryFriendCircle)

	return r
}
