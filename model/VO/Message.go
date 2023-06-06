package VO

type MessageVO struct {
	MsgID      string `json:"id" form:"id" binding:"required"`
	MsgType    int    `json:"msg_type" form:"msg_type" binding:"required"` // 6-接收并通知， 7-接收不通知， 10-好友申请/邀请通知,11-群聊申请/邀请通知 12-朋友圈通知
	Message    string `json:"context" form:"context" binding:"required"`
	CreateTime string `json:"time" form:"time" binding:"required"`
	SenderID   string `json:"send_id" form:"send_id" binding:"required"`
	SenderName string `json:"send_name" form:"send_name" binding:"required"`
	ReceiverID string `json:"receive_id" form:"receive_id" binding:"required"`
	DataType   int    `json:"type" form:"type" binding:"required"`

	ErrString   string `json:"err" form:"err"`
	IsAnonymous bool   `json:"is_anonymous" form:"is_anonymous" binding:"required"`
}
