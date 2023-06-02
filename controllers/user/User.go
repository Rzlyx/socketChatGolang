package user

import (
	"dou_yin/dao/redis"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func Register(c *gin.Context) {
	p := new(param.ParamRegister)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}
	err, p1 := service.Register(p)
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
	err, user, token := service.Login(p)
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
	param := new(param.QueryContactorList)
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
	img := c.Param("id")
	if img == "" {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	pwd := GetCurrentPath()
	imgPath := fmt.Sprintf("%v/img/%v", pwd, img)
	c.File(imgPath)
}

func UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("img")
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	pwd := GetCurrentPath()
	dst := fmt.Sprintf("%v/img/%v", pwd, file.Filename)
	c.SaveUploadedFile(file, dst)

	response.ResponseSuccess(c, struct{}{})
}

func GetCurrentPath() string {
	path, _ := os.Getwd()
	return path
}
