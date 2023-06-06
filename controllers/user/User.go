package user

import (
	"dou_yin/dao/redis"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/utils"
	"dou_yin/service"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	p := new(param.ParamRegister)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	p1, err := service.Register(p)
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccess(c, p1)
}

func Login(c *gin.Context) {
	p := new(param.ParamLogin)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	user, token, err := service.Login(p)
	if err != nil {
		response.ResponseErrorWithMsg(c, response.CodeInvalidPassword, "")
		return
	}

	response.ResponseSuccess(c, gin.H{
		"id":    user.UserID,
		"token": token,
	})
}

func GetContactorList(c *gin.Context) {
	Id := c.Param("id")
	p, err := service.GetContactorList(Id)
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	response.ResponseSuccess(c, p)
}

func TestRedis(c *gin.Context) {
	err := redis.AddMsg("", "123456")
	fmt.Println(err)
}

func QueryContactorList(c *gin.Context) {
	param := new(param.QueryContactorListParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	contactors, err := service.QueryContactorList(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	queryContactorListResp := new(response.QueryContactorList)
	queryContactorListResp.ContactorList = contactors
	response.ResponseSuccess(c, queryContactorListResp)
}

func GetPhotoByID(c *gin.Context) {
	ID := c.Param("id")
	if ID == "" {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	pwd := utils.GetCurrentPath()
	imgPath := fmt.Sprintf("%v/img/%v", pwd, ID)
	c.File(imgPath)
}

func UploadPhoto(c *gin.Context) {
	param := new(param.UploadPhotoParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	file, err := c.FormFile("img")
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	pwd := utils.GetCurrentPath()
	dst := fmt.Sprintf("%v/img/%v", pwd, param.UserID)
	c.SaveUploadedFile(file, dst)

	response.ResponseSuccess(c, struct{}{})
}

func SetContactorList(c *gin.Context) {
	param := new(param.SetContactorListParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	fmt.Println(*param)
	err = service.SetContactorList(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func UpdateUserInfo(c *gin.Context) {
	param := new(param.UpdateUserInfoParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	err = service.UpdateUserInfo(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, struct{}{})
}

func QueryUserInfo(c *gin.Context) {
	param := new(param.QueryUserInfoParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	userInfo, err := service.QueryUserInfo(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	resp := new(response.QueryUserInfo)
	resp.UserInfo = userInfo
	response.ResponseSuccess(c, resp)
}

func SearchFriendOrGroup(c *gin.Context) {
	param := new(param.SearchFriendOrGroupParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	context, err := service.SearchFriendOrGroup(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}
	resp := new(response.SearchFriendOrGroup)
	resp.Context = context
	response.ResponseSuccess(c, resp)
}

func StartSendWebSocket(c *gin.Context) {
	p := new(param.StartSendWebSocketParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[StartSendWebSocket] ShouldBind err is ", err.Error())
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	// 群聊
	go func() {
		err := service.SendGroupNewMsg(p.UserID)
		if err != nil {
			fmt.Println("[StartSendWebSocket], SendGroupNewMsg err is ", err.Error())
		}
		wg.Done()
	}()
	// 私聊
	go func() {
		err := service.QueryAllUnreadPrivateChatMsg(utils.ShiftToNum64(p.UserID))
		if err != nil {
			fmt.Println("[StartSendWebSocket], QueryAllUnreadPrivateChatMsg err is ", err.Error())
		}
		wg.Done()
	}()

	wg.Wait()

	response.ResponseSuccess(c, struct{}{})
}
