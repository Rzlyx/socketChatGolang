package DO

type ContactList struct {
	ContactorList []ContactInfo `json:"contactor_list"`
}

type ContactInfo struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Message string `json:"new_msg" binding:"required"`
	Time    string `json:"time" binding:"required"`
	Status  string `json:"status" binding:"status"`
}

type UserInfo struct {
	UserID      int64  `json:"user_id,string"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	Sex         int    `json:"sex"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"e_mail"`
	Signature   string `json:"signature"`
	Birthday    string `json:"birthday" `
	Status      int    `json:"status"`
	//PrivateChatWhite     *string `json:"private_chat_white"`
	//PrivateChatBlack     *string `json:"private_chat_black"`
	//FriendCircleWhite    *string `json:"friend_circle_white"`
	//FriendCircleBlack    *string `json:"friend_circle_black"`
	FriendCircleVisiable int `json:"friend_circle_visiable"`
	//GroupChatWhite       *string `json:"group_chat_white""`
	//GroupChatBlack       *string `json:"group_chat_black"`
	//GroupChatGray        *string `json:"group_chat_grey"`
	//CreateTime           *string `json:"create_time"`
	//IsDeleted            bool    `json:"is_deleted"`
	//Extra                *string `json:"extra"`
}

type SearchFriendOrGroupContexts struct {
	Result []SearchFriendOrGroupContext `json:"result"`
}

type SearchFriendOrGroupContext struct {
	Type int    `json:"type"`
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}
