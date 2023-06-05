package param

type QueryFriendListParam struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

type QueryFriendInfoParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type AddFriendParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
	Reason   string `json:"reason" form:"reason" binding:"required"`
}

type DeleteFriendParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetPrivateChatBlackParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type UnBlockPrivateChatParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetFriendCircleBlackParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type UnBlockFriendCircleParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type QueryFriendApplyParam struct {
	UserID string `json:"user_id" form:"user_id" binding:"required"`
}

type AgreeFriendApplyParam struct {
	ApplyID  string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type DisagreeFriendApplyParam struct {
	ApplyID  string `json:"apply_id" form:"apply_id" binding:"required"`
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetFriendRemark struct {
	UserID   string  `json:"user_id" form:"user_id" binding:"required"`
	FriendID string  `json:"friend_id" form:"friend_id" binding:"required"`
	Remark   *string `json:"remark" form:"remark"`
}

type SetReadTime struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetPrivateChatGrayParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}

type UnGrayPrivateChatParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
}
