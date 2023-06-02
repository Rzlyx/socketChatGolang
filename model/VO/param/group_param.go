package param

// --------------需要返回详细数据

// 查看群信息
type QueryGroupInfoParam struct {
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 查看已添加群
type QueryGroupListParam struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
}

// 查看加群申请（管理员/群主）
type QueryGroupApplyListParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// ---------------仅需要返回成功与否

// 创建群聊
type CreateGroupInfoParam struct {
	OwnerID     int64    `json:"owner_id" form:"owner_id" binding:"required"`
	GroupName   string   `json:"group_name" form:"group_name" binding:"required"`
	Description *string  `json:"description" form:"description" binding:"required"`
	UserIDs     *[]int64 `json:"user_ids" form:"user_ids"`
}

// 解散群聊
type DissolveGroupInfoParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 申请加群
type ApplyJoinGroupParam struct {
	UserID  int64  `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64  `json:"group_id" form:"group_id" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

// 退群
type QuitGroupParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 同意加群申请（管理员/群主）
type AgreeGroupApplyParam struct {
	ApplyID   int64 `json:"apply_id" form:"apply_id" binding:"required"`
	UserID    int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID   int64 `json:"group_id" form:"group_id" binding:"required"`
	Applicant int64 `json:"applicant" form:"applicant" binding:"required"`
}

// 拒绝加群申请（管理员/群主）
type DisAgreeGroupApplyParam struct {
	ApplyID   int64 `json:"apply_id" form:"apply_id" binding:"required"`
	UserID    int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID   int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 禁言（管理员/群主）
type SilenceParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID  int64 `json:"group_id" form:"group_id" binding:"required"`
	TargetID int64 `json:"target_id" form:"target_id" binding:"required"`
}

// 解除禁言（管理员/群主）
type UnSilenceParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID  int64 `json:"group_id" form:"group_id" binding:"required"`
	TargetID int64 `json:"target_id" form:"target_id" binding:"required"`
}

// 转让群（群主）
type TransferGroupParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID  int64 `json:"group_id" form:"group_id" binding:"required"`
	TargetID int64 `json:"target_id" form:"target_id" binding:"required"`
}

// 屏蔽消息(将该群聊加入黑名单)
type SetBlackListParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 免打扰消息(将该群聊设置为免打扰)
type SetGrayListParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

// 通知消息(将该群聊设为可通知消息)
type SetWhiteListParam struct {
	UserID  int64 `json:"user_id" form:"user_id" binding:"required"`
	GroupID int64 `json:"group_id" form:"group_id" binding:"required"`
}

type SetGroupAdminParam struct {
	GroupID  int64 `json:"group_id" form:"group_id" binding:"required"`
	TargetID int64 `json:"target_id" form:"target_id" binding:"required"`
}

type SetGroupUserParam struct {
	GroupID  int64 `json:"group_id" form:"group_id" binding:"required"`
	TargetID int64 `json:"target_id" form:"target_id" binding:"required"`
}
