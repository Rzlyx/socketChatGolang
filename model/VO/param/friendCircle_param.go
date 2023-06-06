package param

type QueryAllFriendCircleParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	ReadTime string `json:"read_time" form:"read_time" binding:"required"`
	Num      int    `json:"num" form:"num" binding:"required"`
}

type QueryFriendCircleParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}
