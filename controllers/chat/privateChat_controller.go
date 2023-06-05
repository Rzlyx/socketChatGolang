package chat

import (
	"dou_yin/controllers/user"
	"dou_yin/logger"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"fmt"
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

func UploadPrivateChatPhoto(c *gin.Context) {
	p := new(param.UploadPrivateChatPhotoParam)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	file, err := c.FormFile("img")
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	pwd := user.GetCurrentPath()
	dst := fmt.Sprintf("%v/img/%v", pwd, p.Message.MsgID)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	service.HandlePrivateChatMsg(p.Message)

	response.ResponseSuccess(c, struct{}{})
}
