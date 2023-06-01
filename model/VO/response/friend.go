package response

import "dou_yin/service/DO"

type QueryFriendListResp struct {
	Friendlist DO.FriendList `json:"friendlist"`
}

type QueryFriendInfoResp struct {
	FriendInfo DO.FriendInfo `json:"friend_info"`
}

type AddFriendResp struct {
	Application DO.AddFriendApplication `json:"application"`
}

type QueryFriendApplyResp struct {
	Applications DO.FriendApplicationList `json:"applications"`
}
