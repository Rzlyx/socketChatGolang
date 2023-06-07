package param

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
	Message string `json:"message" binding:"required"`
}

type UploadPrivateChatFileParam struct {
	Message string `json:"message" binding:"required"`
}

type QueryPrivateChatMsgByDateParam struct {
	UserID    string `json:"user_id" form:"user_id" binding:"required"`
	FriendID  string `json:"friend_id" form:"friend_id" binding:"required"`
	StartTime string `json:"start_time" form:"start_time" binding:"required"`
	EndTime   string `json:"end_time" form:"end_time" binding:"required"`
}

type QueryPrivateChatMsgByReadTimeParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
	ReadTime string `json:"read_time" form:"read_time" binding:"required"`
	Num      int    `json:"num" form:"num" binding:"required"`
}
