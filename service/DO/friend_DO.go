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
	ApplyID     int64
	ApplicantID int64
	FriendID    int64
	Type        int
	Reason      string
	CreateTime  string
}

type FriendApplicationList struct {
	Applications []FriendApplication
}

type FriendApplication struct {
	ApplyID       int64
	UserID        int64
	ApplicantID   int64
	ApplicantName string
	Reason        string
	Status        int
}

type Friendship struct {
	FriendshipID      int64
	FirstID           int64
	SecondID          int64
	FirstRemarkSecond string
	SecondRemarkFirst string
}
