package user

import (
	"dou_yin/dao/redis"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"fmt"
	"github.com/gin-gonic/gin"
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
