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
	GroupID   int64  `json:"group_id"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
}

// 加群申请
type GroupJoinApply struct {
	ApplyID   int64  `json:"apply_id"`
	Applicant int64  `json:"applicant"`
	TargetID  int64  `json:"target_id"`
	Reason    string `json:"reason"`
}
