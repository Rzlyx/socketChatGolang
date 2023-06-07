package privateChat

import (
	"dou_yin/logger"
	"dou_yin/model/VO"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/utils"
	"dou_yin/service"
	"encoding/json"
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

func QueryPrivateChatMsgByDate(c *gin.Context) {
	param := new(param.QueryPrivateChatMsgByDateParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	messages, err := service.QueryPrivateChatMsgByDate(*param)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeInternError)
		return
	}

	resp := new(response.QueryPrivateChatMsgResp)
	resp.MessageList = messages
	response.ResponseSuccess(c, resp)
}

func QueryPrivateChatMsgByReadTime(c *gin.Context) {
	param := new(param.QueryPrivateChatMsgByReadTimeParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	messages, err := service.QueryPrivateChatMsgByReadTime(*param)
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
	file, err := c.FormFile("img")
	if err != nil {
		fmt.Println("[UploadPrivateChatPhoto], FormFile err is ", err.Error())
		return
	}
	message := c.PostForm("id")
	fmt.Println("[UploadPrivateChatPhoto],message is ", message)

	pwd := utils.GetCurrentPath()
	dst := fmt.Sprintf("%v/img/%v", pwd, message)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		fmt.Println("[UploadPrivateChatPhoto], SaveUploadedFile err is ", err.Error())
		return
	}
	// p := new(param.UploadPrivateChatPhotoParam)
	// err := c.ShouldBind(p)
	// if err != nil {
	//   response.ResponseError(c, response.CodeInvalidParams)
	//   return
	// }

	// file, err := c.FormFile("img")
	// if err != nil {
	//   response.ResponseError(c, response.CodeServerBusy)
	//   return
	// }

	// pwd := utils.GetCurrentPath()
	// dst := fmt.Sprintf("%v/img/%v", pwd, p.Message.MsgID)
	// err = c.SaveUploadedFile(file, dst)
	// if err != nil {
	//   response.ResponseError(c, response.CodeServerBusy)
	//   return
	// }
	// service.HandlePrivateChatMsg(p.Message)

	// response.ResponseSuccess(c, struct{}{})
}

func UploadPrivateChatFile(c *gin.Context) {
	fmt.Println("********['")
	p := new(param.UploadPrivateChatFileParam)
	err := c.ShouldBind(p)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	var message VO.MessageVO
	err = json.Unmarshal([]byte(p.Message), &message)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	pwd := utils.GetCurrentPath()
	dst := fmt.Sprintf("%v/file/%v", pwd, message.MsgID)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		logger.Log.Error(err.Error())
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	service.HandlePrivateChatMsg(message)

	response.ResponseSuccess(c, struct{}{})
}

func GetFileByID(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	pwd := utils.GetCurrentPath()
	filePath := fmt.Sprintf("%v/file/%v", pwd, ID)
	c.File(filePath)
}
