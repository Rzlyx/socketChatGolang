package response

import "dou_yin/service/DO"

// 查看群聊详细信息
type GroupInfo struct {
	GroupID     int64             `json:"group_id,string"`
	OwnerID     int64             `json:"owner_id,string"`
	AdminIds    []int64           `json:"admin_ids"`
	SilenceList []int64           `json:"silence_list"`
	UserIds     []int64           `json:"user_ids"`
	GroupName   string            `json:"group_name"`
	Description string            `json:"description"`
	CreateTime  string            `json:"create_time"`
	IsDeleted   bool              `json:"is_deleted"`
	Extra       DO.GroupInfoExtra `json:"extra"`
}

// 加入群聊信息,仅展示 群名和头像
type GroupJoin struct {
	GroupID   int64  `json:"group_id,string"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
}

// 加群申请
type GroupJoinApply struct {
	ApplyID   int64  `json:"apply_id,string"`
	Applicant int64  `json:"applicant,string"`
	TargetID  int64  `json:"target_id,string"`
	Reason    string `json:"reason"`
}

// 邀请加群申请
type InviteGroupInfo struct {
	ApplyID   int64 `json:"apply_id,string"`
	Applicant int64 `json:"applicant,string"`  // 发起邀请人
	TargetID  int64 `json:"target_id,string"`  // 群ID
	InvitedID int64 `json:"invited_id,string"` // 被邀请人
}
