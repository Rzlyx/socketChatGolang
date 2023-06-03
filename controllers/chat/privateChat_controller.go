package chat

import (
	"dou_yin/logger"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"github.com/gin-gonic/gin"
)

func QueryPrivateChatMsg(c *gin.Context) {
	param := new(param.QueryPrivateChatMsgParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	messages, err := service.QueryPrivateChatMsg(*param)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeInternError)
		return
	}

	resp := new(response.QueryPrivateChatMsgResp)
	resp.MessageList = messages
	response.ResponseSuccess(c, resp)
}

func DeletePrivateChatMsg(c *gin.Context) {
	param := new(param.DeletePrivateChatMsgParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.DeletePrivateChatMsg(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}
