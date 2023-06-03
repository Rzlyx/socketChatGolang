package DO

type GroupMsgDO struct {
	MsgID       int64          `json:"message_id" db:"message_id"`
	GroupID     int64          `json:"group_id" db:"group_id"`
	SenderID    int64          `json:"sender_id" db:"sender_id"`
	Message     string         `json:"message" db:"message"`
	Type        int            `json:"type" db:"type"`
	IsAnonymous bool           `json:"is_anonymous" db:"is_anonymous"`
	CreateTime  string         `json:"create_time" db:"create_time"`
	DeletedList *[]int64       `json:"deleted_list" db:"deleted_list"`
	Extra       *GroupMsgExtra `json:"extra" db:"extra"`
}

type GroupMsgExtra struct {
}