package controllers

import (
	"dou_yin/dao/redis"
	"dou_yin/model"
	"dou_yin/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	p := new(model.ParamRegister)
	err := c.ShouldBind(p)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	err, p1 := service.Register(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, p1)
}
func Login(c *gin.Context) {
	p := new(model.ParamLogin)
	err := c.ShouldBind(p)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}
	err, user, token := service.Login(p)
	if err != nil {
		ResponseErrorWithMsg(c, CodeInvalidPassword, "")
		return
	}

	ResponseSuccess(c, gin.H{
		"id":    user.UserID,
		"token": token,
	})
}
func GetContactorList(c *gin.Context) {
	Id := c.Param("id")
	p, err := service.GetContactorList(Id)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, p)
}
func TestRedis(c *gin.Context) {
	err := redis.AddMsg("", "123456")
	fmt.Println(err)
}
