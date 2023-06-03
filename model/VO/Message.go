package VO

type MessageVO struct {
	MsgID      string `json:"id" form:"id" binding:"required"`
	MsgType    int    `json:"msg_type" form:"msg_type" binding:"required"`
	Message    string `json:"context" form:"context" binding:"required"`
	CreateTime string `json:"time" form:"time" binding:"required"`
	SenderID   string `json:"send_id" form:"send_id" binding:"required"`
	ReceiverID string `json:"receive_id" form:"receive_id" binding:"required"`
	DataType   int    `json:"type" form:"type" binding:"required"`
	IsAnonymous bool `json:"is_anonymous" form:"is_anonymous" binding:"required"`
}
