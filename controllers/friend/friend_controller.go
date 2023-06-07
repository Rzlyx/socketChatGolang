package friend

import (
	"dou_yin/logger"
	param "dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"github.com/gin-gonic/gin"
)

func QueryFriendList(c *gin.Context) {
	param := new(param.QueryFriendListParam)
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
	queryFriendListResp.FriendList = friendlist
	response.ResponseSuccess(c, queryFriendListResp)
}

func QueryFriendInfo(c *gin.Context) {
	param := new(param.QueryFriendInfoParam)
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
	param := new(param.AddFriendParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	application, err := service.AddFriend(*param)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeInternError)
		return
	}

	addFriendResp := new(response.AddFriendResp)
	addFriendResp.Application = application
	response.ResponseSuccess(c, addFriendResp)
}

func DeleteFriend(c *gin.Context) {
	param := new(param.DeleteFriendParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.DeleteFriend(*param)
	if err != nil {
		if err.Error() == "not friend" {
			response.ResponseError(c, response.CodeNotFriend)
		} else {
			logger.Log.Error(err.Error())
			response.ResponseError(c, response.CodeInternError)
		}
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func SetPrivateChatBlack(c *gin.Context) {
	param := new(param.SetPrivateChatBlackParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.SetPrivateChatBlack(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func UnBlockPrivateChat(c *gin.Context) {
	param := new(param.UnBlockPrivateChatParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.UnBlockPrivateChat(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func SetFriendCircleBlack(c *gin.Context) {
	param := new(param.SetFriendCircleBlackParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.SetFriendCircleBlack(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func UnBlockFriendCircle(c *gin.Context) {
	param := new(param.UnBlockFriendCircleParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.UnBlockFriendCircle(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func QueryFriendApply(c *gin.Context) {
	param := new(param.QueryFriendApplyParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	applications, err := service.QueryFriendApply(*param)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeInternError)
		return
	}

	resp := new(response.QueryFriendApplyResp)
	resp.ApplicationList = applications
	response.ResponseSuccess(c, resp)
}

func AgreeFriendApply(c *gin.Context) {
	param := new(param.AgreeFriendApplyParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.AgreeFriendApply(*param)
	if err != nil {
		if err.Error() == "there is no application" {
			response.ResponseError(c, response.CodeNotApplied)
		} else {
			logger.Log.Error(err.Error())
			response.ResponseError(c, response.CodeInternError)
		}
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func DisagreeFriendApply(c *gin.Context) {
	param := new(param.DisagreeFriendApplyParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.DisagreeFriendApply(*param)
	if err != nil {
		if err.Error() == "there is no application" {
			response.ResponseError(c, response.CodeNotApplied)
		} else {
			logger.Log.Error(err.Error())
			response.ResponseError(c, response.CodeInternError)
		}
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func SetFriendRemark(c *gin.Context) {
	param := new(param.SetFriendRemark)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	logger.Log.Info((*param).FriendID + " " + (*param).UserID + " " + *param.Remark)
	err = service.SetFriendRemark(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func SetReadTime(c *gin.Context) {
	p := new(param.SetReadTime)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	if len(p.FriendID) > 15 {
		err = service.SetGroupReadTimebyParam(&param.SetGroupReadTimeParam{
			UserID: p.UserID,
			GroupID: p.FriendID,
		})
		if err != nil {
			response.ResponseError(c, response.CodeInternError)
			return
		}	
	}

	err = service.SetReadTime(*p)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func SetPrivateChatGray(c *gin.Context) {
	param := new(param.SetPrivateChatGrayParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.SetPrivateChatGray(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func UnGrayPrivateChat(c *gin.Context) {
	param := new(param.UnGrayPrivateChatParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.UnGrayPrivateChat(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func AddFriendTag(c *gin.Context) {
	param := new(param.AddFriendTagParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.AddFriendTag(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func RemoveFriendTag(c *gin.Context) {
	param := new(param.RemoveFriendTagParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.RemoveFriendTag(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}
