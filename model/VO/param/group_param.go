package param

// --------------需要返回详细数据

// 查看群信息
type QueryGroupInfoParam struct {
	GroupID int64 `json:"group_id"`
}

// 查看已添加群
type QueryGroupListParam struct {
	UserID int64 `json:"user_id"`
}

// 查看加群申请（管理员/群主）
type QueryGroupApplyListParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

// ---------------仅需要返回成功与否

// 申请加群
type ApplyJoinGroupParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

// 撤回申请加群
type RetractApplyParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	ApplyID int64 `json:"apply_id"`
}

// 退群
type QuitGroupParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

// 同意加群申请（管理员/群主）
type AgreeGroupApplyParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	Applicant int64 `json:"applicant"`
}

// 拒绝加群申请（管理员/群主）
type DisAgreeGroupApplyParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	Applicant int64 `json:"applicant"`
}

// 禁言（管理员/群主）
type SilenceParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	TargetID int64 `json:"target_id"`
}

// 解除禁言（管理员/群主）
type UnSilenceParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	TargetID int64 `json:"target_id"`
}

// 转让群（群主）
type TransferGroupParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
	TargetID int64 `json:"target_id"`
}

// 屏蔽消息(将该群聊加入黑名单)
type SetBlackListParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

// 免打扰消息(将该群聊设置为免打扰)
type SetGrayListParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}

// 通知消息(将该群聊设为可通知消息)
type SetWhiteListParam struct {
	UserID  int64 `json:"user_id"`
	GroupID int64 `json:"group_id"`
}