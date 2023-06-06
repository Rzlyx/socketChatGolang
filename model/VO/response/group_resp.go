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
	MyGroupName string            `json:"my_group_name"` // 对群名的备注
	MyName      string            `json:"my_name"`       // 在此群的昵称
	Type        int               `json:"type"`          // 群身份 0-普通成员，1-管理员，2-群主
	OnlineNum   int               `json:"online_num"`    // 在线人数
	MsgType     int               `json:"msg_type"`      // 消息类型 6-接收并通知，7-接收不通知，8-不接收
}

type GroupUserInfo struct {
	UserID     int64  `json:"user_id,string"`
	MyName     string `json:"name"`     // 在此群的昵称
	InsertTime string `json:"date"`     // 入群时间
	Status     int    `json:"online"`   // 0-离线，1-在线,2-潜水
	Type       int    `json:"identity"` // 群身份 0-普通成员，1-管理员，2-群主
	IsSlienced bool   `json:"status"`   // 是否被禁言
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
	ApplyID   int64  `json:"apply_id,string"`
	Applicant int64  `json:"applicant,string"`  // 被邀请人
	TargetID  int64  `json:"target_id,string"`  // 群ID
	InvitedID int64  `json:"invited_id,string"` // 发起邀请人
	Reason    string `json:"reason"`
}
