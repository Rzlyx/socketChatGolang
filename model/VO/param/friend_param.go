package param

type QueryFriendListParam struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
}

type QueryFriendInfoParam struct {
	FriendshipID int64 `json:"friendship_id" form:"friendship_id" binding:"required"`
	UserID       int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID     int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type AddFriendParam struct {
	UserID   int64  `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64  `json:"friend_id" form:"friend_id" binding:"required"`
	Reason   string `json:"reason" form:"reason" binding:"required"`
}

type DeleteFriendParam struct {
	UserID       int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID     int64 `json:"friend_id" form:"friend_id" binding:"required"`
	FriendshipID int64 `json:"friendship_id" form:"friendship_id" binding:"required"`
}

type SetPrivateChatBlackParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type UnBlockPrivateChatParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetFriendCircleBlackParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type UnBlockFriendCircleParam struct {
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type QueryFriendApplyParam struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
}

type AgreeFriendApplyParam struct {
	ApplyID  int64 `json:"apply_id" form:"apply_id" binding:"required"`
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type DisagreeFriendApplyParam struct {
	ApplyID  int64 `json:"apply_id" form:"apply_id" binding:"required"`
	UserID   int64 `json:"user_id" form:"user_id" binding:"required"`
	FriendID int64 `json:"friend_id" form:"friend_id" binding:"required"`
}

type SetFriendRemark struct {
	UserID       int64   `json:"user_id" form:"user_id" binding:"required"`
	FriendID     int64   `json:"friend_id" form:"friend_id" binding:"required"`
	FriendshipID int64   `json:"friendship_id" form:"friendship_id" binding:"required"`
	Remark       *string `json:"remark" form:"remark"`
}
