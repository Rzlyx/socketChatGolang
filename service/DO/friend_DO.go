package DO

type FriendList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	FriendshipID int64  `json:"friendship_id,string"`
	FriendID     int64  `json:"friend_id,string"`
	Name         string `json:"name"`
	Status       int    `json:"status"`
}

type FriendInfo struct {
	UserID              int64  `json:"user_id,string"`
	UserName            string `json:"user_name"`
	Sex                 int    `json:"sex"`
	PhoneNumber         string `json:"phone_number"`
	Email               string `json:"e_mail"`
	Signature           string `json:"signature"`
	Birthday            string `json:"birthday"`
	Status              int    `json:"status"`
	Remark              string `json:"remark"`
	IsRemark            bool   `json:"is_remark"`
	IsPrivateChatBlack  bool   `json:"is_private_chat_black"`
	IsFriendCircleBlack bool   `json:"is_friend_circle_black"`
	IsPrivateChatGray   bool   `json:"is_private_chat_gray"`
}

type AddFriendApplication struct {
	// 是否通过对方申请之间成为好友
	IsBeFriend  bool   `json:"is_be_friend"`
	ApplyID     int64  `json:"apply_id,string"`
	ApplicantID int64  `json:"applicant_id,string"`
	FriendID    int64  `json:"friend_id,string"`
	Type        int    `json:"type"`
	Reason      string `json:"reason"`
	CreateTime  string `json:"create_time"`
}

type FriendApplicationList struct {
	Applications []FriendApplication `json:"applications"`
}

type FriendApplication struct {
	ApplyID       int64  `json:"apply_id,string"`
	UserID        int64  `json:"user_id,string"`
	ApplicantID   int64  `json:"applicant_id,string"`
	ApplicantName string `json:"applicant_name"`
	Reason        string `json:"reason"`
	Status        int    `json:"status"`
}

type Friendship struct {
	FriendshipID      int64  `json:"friendship_id,string"`
	FirstID           int64  `json:"first_id,string"`
	SecondID          int64  `json:"second_id,string"`
	FirstRemarkSecond string `json:"first_remark_second"`
	SecondRemarkFirst string `json:"second_remark_first"`
}

type CheckFriendApply struct {
	FriendshipID int64 `json:"friendship_id,string"`
	FirstID      int64 `json:"first_id,string"`
	SecondID     int64 `json:"second_id,string"`
}
