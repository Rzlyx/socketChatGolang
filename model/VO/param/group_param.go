package param

import "dou_yin/model/VO"

// --------------需要返回详细数据

// 查看群信息
type QueryGroupInfoParam struct {
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
}

type GetGroupAllUserParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 查看已添加群
type QueryGroupListParam struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

// 查看加群申请（管理员/群主）
type QueryGroupApplyListParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// ---------------仅需要返回成功与否

// 创建群聊
type CreateGroupInfoParam struct {
	OwnerID     string    `json:"owner_id" form:"owner_id" binding:"required"`
	GroupName   string    `json:"group_name" form:"group_name" binding:"required"`
	Description *string   `json:"description" form:"description" binding:"required"`
	UserIDs     *[]string `json:"user_ids" form:"user_ids"`
}

type UpdateGroupInfoParam struct {
	GroupID     string  `json:"group_id" form:"group_id" binding:"required"`
	GroupName   string  `json:"group_name" form:"group_name" binding:"required"`
	Description *string `json:"description" form:"description" binding:"required"`
}

type UploadGroupPhotoParam struct {
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 解散群聊
type DissolveGroupInfoParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 申请加群
type ApplyJoinGroupParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

// 退群
type QuitGroupParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// type KickUserFromGroupParam struct {
// 	UserID   string `json:"user_id" form:"user_id" binding:"required"`
// 	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
// 	TargetID string `json:"target_id" form:"target_id" binding:"required"`
// }

// 同意加群申请（管理员/群主）
type AgreeGroupApplyParam struct {
	ApplyID   string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	GroupID   string `json:"group_id" form:"group_id" binding:"required"`
	Applicant string `json:"applicant" form:"applicant" binding:"required"`
}

// 拒绝加群申请（管理员/群主）
type DisAgreeGroupApplyParam struct {
	ApplyID string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 禁言（管理员/群主）
type SilenceParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
}

// 解除禁言（管理员/群主）
type UnSilenceParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
}

// 转让群（群主）
type TransferGroupParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
}

// 屏蔽消息(将该群聊加入黑名单)
type SetBlackListParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 免打扰消息(将该群聊设置为免打扰)
type SetGrayListParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

// 通知消息(将该群聊设为可通知消息)
type SetWhiteListParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

type SetGroupAdminParam struct {
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
}

type SetGroupUserParam struct {
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
}

type InviteJoinGroupParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	GroupID  string `json:"group_id" form:"group_id" binding:"required"`
	TargetID string `json:"target_id" form:"target_id" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

type QueryInviteGroupParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
}

type AgreeInviteGroupParam struct {
	ApplyID   string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	GroupID   string `json:"group_id" form:"group_id" binding:"required"`
	Applicant string `json:"applicant" form:"applicant" binding:"required"`
}

type DisAgreeInviteGroupParam struct {
	ApplyID string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}
type SetGroupNameParam struct {
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	GroupID   string `json:"group_id" form:"group_id" binding:"required"`
	GroupName string `json:"group_name" form:"group_name" binding:"required"`
}

type SetMyNameParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
	MyName  string `json:"my_name" form:"my_name" binding:"required"`
}

type SetGroupReadTimeParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

type SetAIGPTParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

type GetPageOldMsgParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
	PageNum int    `json:"page_num" form:"page_num" binding:"required"`
	Num     int    `json:"num" form:"num" binding:"required"`
}

type GetGroupOldMsgLoginParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
}

type GetGroupOldMsgUpParam struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	GroupID string `json:"group_id" form:"group_id" binding:"required"`
	TimeTag string `json:"time_tag" form:"time_tag" binding:"required"`
}

type GetGroupOldMsgDayParam struct {
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	GroupID   string `json:"group_id" form:"group_id" binding:"required"`
	StartTime string `json:"start_time" form:"start_time" binding:"required"`
	EndTime   string `json:"end_time" form:"end_time" binding:"required"`
}

type UploadGroupChatPhotoParam struct {
	Message VO.MessageVO `json:"message" binding:"required"`
}

type UploadGroupChatFileParam struct {
	Message VO.MessageVO `json:"message" binding:"required"`
}
