package response

import "dou_yin/service/DO"

type QueryFriendListResp struct {
	Friendlist DO.FriendList `json:"friendlist"`
}

type QueryFriendInfoResp struct {
	FriendInfo DO.FriendInfo `json:"friend_info"`
}
