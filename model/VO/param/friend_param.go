package param

type QueryFriendListParam struct {
	UserID int64
}

type QueryFriendInfoParam struct {
	FriendshipID int64
	FriendID     int64
}

type AddFriendParam struct {
}

type RetractAddFriendParam struct {
}

type DeleteFriendParam struct {
}

type SetPrivateChatBlackParam struct {
}

type UnBlockPrivateChatParam struct {
}

type QueryFriendApplyParam struct {
}

type AgreeFriendApplyParam struct {
}

type DisagreeFriendApplyParam struct {
}
