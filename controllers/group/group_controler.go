package group

import (
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/utils"
	"dou_yin/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// 查看群信息
func QueryGroupInfo(c *gin.Context) {
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

func GetGroupAllUser(c *gin.Context) {
	p := new(param.GetGroupAllUserParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[GetGroupAllUser] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.GetGroupAllUserbyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[GetGroupAllUser] MGetGroupListByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}

// 创建群聊
func CreateGroupInfo(c *gin.Context) {
	p := new(param.CreateGroupInfoParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[CreateGroupInfo] ShouldBind err is ", err.Error())
		return
	}
	err = service.CreateGroupInfoByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[CreateGroupInfo] CreateGroupInfoByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 更新群信息
func UpdateGroupInfo(c *gin.Context) {
	p := new(param.UpdateGroupInfoParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[CreateGroupInfo] ShouldBind err is ", err.Error())
		return
	}
	err = service.UpdateGroupInfoByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[CreateGroupInfo] CreateGroupInfoByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 上传头像
func UploadGroupPhoto(c *gin.Context) {
	param := new(param.UploadGroupPhotoParam)
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
	dst := fmt.Sprintf("%v/img/%v", pwd, param.GroupID)
	c.SaveUploadedFile(file, dst)

	response.ResponseSuccess(c, struct{}{})
}

// 解散群
func DissolveGroupInfo(c *gin.Context) {
	p := new(param.DissolveGroupInfoParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[DissolveGroupInfo] ShouldBind err is ", err.Error())
		return
	}
	err = service.DissolveGroupInfoByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[DissolveGroupInfo] MGetGroupListByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 申请加群
func ApplyJoinGroup(c *gin.Context) {
	p := new(param.ApplyJoinGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[ApplyJoinGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.ApplyJoinGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[ApplyJoinGroup] MGetGroupListByParam err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
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

// 踢人
func KickUserFromGroup(c *gin.Context) {
	p := new(param.QuitGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[KickUserFromGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.QuitGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[KickUserFromGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 查看加群申请（管理员/群主）
func QueryGroupApplyList(c *gin.Context) {
	p := new(param.QueryGroupApplyListParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[QueryGroupApplyList] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.QueryGroupApplyListByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[QueryGroupApplyList]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, *p1)
}

// 同意加群申请（管理员/群主）
func AgreeGroupApply(c *gin.Context) {
	p := new(param.AgreeGroupApplyParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[AgreeGroupApply] ShouldBind err is ", err.Error())
		return
	}
	err = service.AgreeGroupApplyByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[AgreeGroupApply]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 拒绝加群申请（管理员/群主）
func DisAgreeGroupApply(c *gin.Context) {
	p := new(param.DisAgreeGroupApplyParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[DisAgreeGroupApply] ShouldBind err is ", err.Error())
		return
	}
	err = service.DisAgreeGroupApplyByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[DisAgreeGroupApply]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
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

// 将普通成员设为管理员
func SetGroupAdmin(c *gin.Context) {
	p := new(param.SetGroupAdminParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetGroupAdmin] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetGroupAdminByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetGroupAdmin]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 将管理员设为普通成员
func SetGroupUser(c *gin.Context) {
	p := new(param.SetGroupUserParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetGroupUser] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetGroupUserByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetGroupUser]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 邀请加入群聊
func InviteJoinGroup(c *gin.Context) {
	p := new(param.InviteJoinGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[InviteJoinGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.InviteJoinGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[InviteJoinGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 查看邀请
func QueryInviteGroup(c *gin.Context) {
	p := new(param.QueryInviteGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[QueryInviteGroup] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.QueryInviteGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[QueryInviteGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}

// 同意邀请
func AgreeInviteGroup(c *gin.Context) {
	p := new(param.AgreeInviteGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[AgreeInviteGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.AgreeInviteGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[AgreeInviteGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 拒绝邀请
func DisAgreeInviteGroup(c *gin.Context) {
	p := new(param.DisAgreeInviteGroupParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[DisAgreeInviteGroup] ShouldBind err is ", err.Error())
		return
	}
	err = service.DisAgreeInviteGroupByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[DisAgreeInviteGroup]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 设置群备注
func SetGroupName(c *gin.Context) {
	p := new(param.SetGroupNameParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetGroupName] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetGroupNameByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetGroupName]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

func SetMyName(c *gin.Context) {
	p := new(param.SetMyNameParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetMyName] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetMyNamebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetMyName]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 设置群聊已读时间
func SetGroupReadTime(c *gin.Context) {
	p := new(param.SetGroupReadTimeParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[SetGroupReadTime] ShouldBind err is ", err.Error())
		return
	}
	err = service.SetGroupReadTimebyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[SetGroupReadTime]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, struct{}{})
}

// 按页查询历史消息
func GetPageOldMsg(c *gin.Context) {
	p := new(param.GetPageOldMsgParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[GetPageOldMsg] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.GetPageOldMsgByParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[GetPageOldMsg]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}


// 登录获取历史消息--获取15条消息
func GetGroupOldMsgLogin(c *gin.Context) {
	p := new(param.GetGroupOldMsgLoginParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[GetGroupOldMsgLogin] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.GetGroupOldMsgLoginbyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[GetGroupOldMsgLogin]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)

	go func(){
		err := service.SendGroupNewMsg(p.UserID)
		if err != nil {
			fmt.Println("[GetGroupOldMsgLogin], SenndGroupNewMsg err is ", err.Error())
		}
	}()
}

// 加载向上的信息
func GetGroupOldMsgUp(c *gin.Context) {
	p := new(param.GetGroupOldMsgUpParam)
	err := c.ShouldBind(p)
	if err != nil {
		// 无效参数
		response.ResponseError(c, response.CodeInvalidParams)
		fmt.Println("[GetGroupOldMsgUp] ShouldBind err is ", err.Error())
		return
	}
	p1, err := service.GetGroupOldMsgUpbyParam(p)
	if err != nil {
		// 内部错误
		response.ResponseError(c, response.CodeInternError)
		fmt.Println("[GetGroupOldMsgUp]  err is ", err.Error())
		return
	}
	response.ResponseSuccess(c, p1)
}

// 按天数获取信息
func GetGroupOldMsgDay() {

}

// 上传群相册
func UploadGroupChatPhoto(c *gin.Context) {
	p := new(param.UploadGroupChatPhotoParam)
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

	pwd := utils.GetCurrentPath()
	dst := fmt.Sprintf("%v/img/%v", pwd, p.Message.MsgID)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	service.HandleGroupChatMsg(&p.Message)

	response.ResponseSuccess(c, struct{}{})
}

// 上传群文件
func UploadGroupChatFile(c *gin.Context) {
	p := new(param.UploadGroupChatFileParam)
	err := c.ShouldBind(p)
	if err != nil {
		response.ResponseError(c, response.CodeInvalidParams)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}

	pwd := utils.GetCurrentPath()
	dst := fmt.Sprintf("%v/file/%v", pwd, p.Message.MsgID)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.ResponseError(c, response.CodeServerBusy)
		return
	}
	service.HandleGroupChatMsg(&p.Message)

	response.ResponseSuccess(c, struct{}{})
}