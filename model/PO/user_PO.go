package PO

type User struct {
	UserID   int64  `json:"user_id" db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
	EMail    string `json:"e_mail" db:"e_mail"`
}
type Contactor struct {
	UserID   int64  `json:"id" db:"user_id"`
	UserName string `json:"user_name" db:"username"`
	Remark   string `json:"remark" db:"remark"`
}
type ContactorList struct {
	ContactorList []*Contactor `json:"contactor_list"`
}

type UserPO struct {
	UserID               int64   `json:"user_id" db:"user_id"`
	UserName             string  `json:"user_name" db:"user_name"`
	Password             string  `json:"password" db:"password"`
	Sex                  int     `json:"sex" db:"sex"`
	PhoneNumber          string  `json:"phone_number" db:"phone_number"`
	Email                string  `json:"e_mail" db:"e_mail"`
	Signature            *string `json:"signature" db:"signature"`
	Birthday             string  `json:"birthday" db:"birthday"`
	Status               int     `json:"status" db:"status"`
	PrivateChatWhite     *string `json:"private_chat_white" db:"private_chat_white"`
	PrivateChatBlack     *string `json:"private_chat_black" db:"private_chat_black"`
	FriendCircleWhite    *string `json:"friend_circle_white" db:"friend_circle_white"`
	FriendCircleBlack    *string `json:"friend_circle_black" db:"friend_circle_black"`
	FriendCircleVisiable int     `json:"friend_circle_visiable" db:"friend_circle_visiable"`
	GroupChatWhite       *string `json:"group_chat_white" db:"group_chat_white"`
	GroupChatBlack       *string `json:"group_chat_black" db:"group_chat_black"`
	GroupChatGray        *string `json:"group_chat_grey" db:"group_chat_gray"`
	CreateTime           *string `json:"create_time" db:"create_time"`
	IsDeleted            bool    `json:"is_deleted" db:"is_deleted"`
	Extra                *string `json:"extra" db:"extra"`
}

type PrivateChatWhite struct {
	WhiteList []int64 `json:"white_list"`
}

type PrivateChatBlack struct {
	BlackList []int64 `json:"black_list"`
}

type FriendCircleWhite struct {
	WhiteList []int64 `json:"white_list"`
}

type FriendCircleBlack struct {
	BlackList []int64 `json:"black_list"`
}
