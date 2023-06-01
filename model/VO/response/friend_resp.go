package response

import "dou_yin/service/DO"

type QueryFriendListResp struct {
	FriendList DO.FriendList `json:"friend_list"`
}

type QueryFriendInfoResp struct {
	FriendInfo DO.FriendInfo `json:"friend_info"`
}

type AddFriendResp struct {
	Application DO.AddFriendApplication `json:"application"`
}

type QueryFriendApplyResp struct {
	ApplicationList DO.FriendApplicationList `json:"applications"`
}
