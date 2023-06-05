package param

import "dou_yin/model/VO"

type QueryPrivateChatMsgParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
	Num      int    `json:"num" form:"num" binding:"required"`
	PageNum  int    `json:"page_num" form:"page_num" binding:"required"`
}

type DeletePrivateChatMsgParam struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
	MsgID  string `json:"msg_id" form:"msg_id" binding:"required"`
}

type UploadPrivateChatPhotoParam struct {
	Message VO.MessageVO `json:"message" binding:"required"`
}

type UploadPrivateChatFileParam struct {
	Message VO.MessageVO `json:"message" binding:"required"`
}
