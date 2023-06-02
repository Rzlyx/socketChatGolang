package controllers

import (
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 查看群信息
func QueryGroupInfo(c *gin.Context){
	p := new(param.QueryGroupInfoParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[QueryGroupInfo] ShouldBind err is ", err.Error())
		return
	}

	p1, err := service.MGetGroupInfoByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[QueryGroupInfo] MGetGroupInfoByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}

// 查看已添加群
func QueryGroupList(c *gin.Context) {
	p := new(param.QueryGroupListParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[QueryGroupList] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.MGetGroupListByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[QueryGroupList] MGetGroupListByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}

// 创建群聊
func CreateGroupInfo(c *gin.Context) {

}

// 解散群
func DissolveGroupInfo(c *gin.Context) {

}

// 申请加群
func ApplyJoinGroup(c *gin.Context) {

}

// 撤回申请加群
func RetractApply(c *gin.Context) {

}

// 退群
func QuitGroup(c *gin.Context) {
	p := new(param.QuitGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[QuitGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.QuitGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[QuitGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 查看加群申请（管理员/群主）
func QueryGroupApplyList(c *gin.Context) {

}

// 同意加群申请（管理员/群主）
func AgreeGroupApply(c *gin.Context) {

}

// 拒绝加群申请（管理员/群主）
func DisAgreeGroupApply(c *gin.Context) {

}

// 禁言（管理员/群主）
func Silence(c *gin.Context) {
	p := new(param.SilenceParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[Silence] ShouldBind err is ", err.Error())
		return
	}
	err = service.SilenceByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[Silence]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 解除禁言（管理员/群主）
func UnSilence(c *gin.Context) {
	p := new(param.UnSilenceParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[UnSilence] ShouldBind err is ", err.Error())
		return
	}
	err = service.UnSilenceByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[UnSilence]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 转让群（群主）
func TransferGroup(c *gin.Context) {
	p := new(param.TransferGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[TransferGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.TransferGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[TransferGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 屏蔽消息(将该群聊加入黑名单)
func SetBlackList(c *gin.Context) {
	p := new(param.SetBlackListParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetBlackList] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetBlackListByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetBlackList]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 免打扰消息(将该群聊设置为免打扰)
func SetGrayList(c *gin.Context) {
	p := new(param.SetGrayListParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetGrayList] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetGrayListByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetGrayList]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 通知消息(将该群聊设为可通知消息)
func SetWhiteList(c *gin.Context) {
	p := new(param.SetWhiteListParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetWhiteList] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetWhiteListByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetWhiteList]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}