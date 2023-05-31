package friend

import (
	param2 "dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"github.com/gin-gonic/gin"
)

func QueryFriendList(c *gin.Context) {
	param := new(param2.QueryFriendListParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	friendlist, err := service.QueryFriendList(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	queryFriendListResp := new(response.QueryFriendListResp)
	queryFriendListResp.Friendlist = friendlist
	response.ResponseSuccess(c, queryFriendListResp)
}

func QueryFriendInfo(c *gin.Context) {
	param := new(param2.QueryFriendInfoParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	friend, err := service.QueryFriendInfo(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	queryFriendInfoResp := new(response.QueryFriendInfoResp)
	queryFriendInfoResp.FriendInfo = friend
	response.ResponseSuccess(c, queryFriendInfoResp)
}

func AddFriend(c *gin.Context) {

}

func RetractAddFriend(c *gin.Context) {

}

func DeleteFriend(c *gin.Context) {

}

func SetPrivateChatBlack(c *gin.Context) {

}

func UnBlockPrivateChat(c *gin.Context) {

}

// todo: friend circle

func QueryFriendApply(c *gin.Context) {

}

func AgreeFriendApply(c *gin.Context) {

}

func DisagreeFriendApply(c *gin.Context) {

}
