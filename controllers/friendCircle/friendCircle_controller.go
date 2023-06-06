package friendCircle

import (
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"github.com/gin-gonic/gin"
)

func QueryFriendCircle(c *gin.Context) {
	param := new(param.QueryAllFriendCircleParam)
	err := c.ShouldBind(param)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	context, err := service.QueryAllFriendCircle(*param)
	if err != nil {
		response.ResponseError(c, response.CodeInternError)
		return
	}

	response.ResponseSuccess(c, context)
}
