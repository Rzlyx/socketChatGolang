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
		fmt.Println("[QueryGroupList] MGetGroupListByParam err is ", err)
		return
	}
	response.ResponseSuccess(c, p1)
}

// 申请加群
func ApplyJoinGroup(c *gin.Context) {

}

// 撤回申请加群
func RetractApply(c *gin.Context) {

}

// 退群
func QuitGroup(c *gin.Context) {

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

}

// 解除禁言（管理员/群主）
func UnSilence(c *gin.Context) {

}

// 转让群（群主）
func TransferGroup(c *gin.Context) {

}

// 屏蔽消息(将该群聊加入黑名单)
func SetBlackList(c *gin.Context) {

}

// 免打扰消息(将该群聊设置为免打扰)
func SetGrayList(c *gin.Context) {

}

// 通知消息(将该群聊设为可通知消息)
func SetWhiteList(c *gin.Accounts) {

}