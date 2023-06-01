package DO

type FriendList struct {
	Friends []Friend `json:"friends"`
}

type Friend struct {
	FriendshipID int64  `json:"friendship_id"`
	FriendID     int64  `json:"friend_id"`
	Name         string `json:"name"`
}

type FriendInfo struct {
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	Sex         int    `json:"sex"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"e_mail"`
	Signature   string `json:"signature"`
	Birthday    string `json:"birthday"`
	Status      int    `json:"status"`
	Remark      string `json:"remark"`
	IsRemark    bool   `json:"is_remark"`
}

type AddFriendApplication struct {
	ApplyID     int64  `json:"apply_id"`
	ApplicantID int64  `json:"applicant_id"`
	FriendID    int64  `json:"friend_id"`
	Type        int    `json:"type"`
	Reason      string `json:"reason"`
	CreateTime  string `json:"create_time"`
}

type FriendApplicationList struct {
	Applications []FriendApplication `json:"applications"`
}

type FriendApplication struct {
	ApplyID       int64  `json:"apply_id"`
	UserID        int64  `json:"user_id"`
	ApplicantID   int64  `json:"applicant_id"`
	ApplicantName string `json:"applicant_name"`
	Reason        string `json:"reason"`
	Status        int    `json:"status"`
}

type Friendship struct {
	FriendshipID      int64  `json:"friendship_id"`
	FirstID           int64  `json:"first_id"`
	SecondID          int64  `json:"second_id"`
	FirstRemarkSecond string `json:"first_remark_second"`
	SecondRemarkFirst string `json:"second_remark_first"`
}
