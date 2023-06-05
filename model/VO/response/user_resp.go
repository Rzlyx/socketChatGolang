package response

import "dou_yin/service/DO"

type QueryContactorList struct {
	ContactorList DO.ContactList `json:"contactor_list"`
}

type QueryUserInfo struct {
	UserInfo DO.UserInfo `json:"user_info"`
}

type SearchFriendOrGroup struct {
	Context DO.SearchFriendOrGroupContexts `json:"context"`
}
