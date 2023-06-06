package param

type QueryAllFriendCircleParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	ReadTime string `json:"read_time" form:"read_time" binding:"required"`
	Num      int    `json:"num" form:"num" binding:"required"`
}

type QueryFriendCircleParam struct {
	UserID   string `json:"user_id" form:"user_id" binding:"required"`
	FriendID string `json:"friend_id" form:"friend_id" binding:"required"`
	ReadTime string `json:"read_time" form:"read_time" binding:"required"`
	Num      int    `json:"num" form:"num" binding:"required"`
}

type SendCircleParam struct {
	Sender     string   `json:"sender" form:"sender" binding:"required"`
	News       string   `json:"news" form:"news" binding:"required"`
	Type       int      `json:"type" form:"type" binding:"required"`
	BlackList  []string `json:"black_list" form:"black_list" binding:"required"`
	CircleType string   `json:"circle_type" form:"circle_type" binding:"required"` // 0-私密 1-公开
	// WhiteList []string `json:"" form:"" binding:"required"`
}

type UploadCirclePhotoParam struct {
	NewsID  string       `json:"news_id" form:"news_id" binding:"required"`
	// Message VO.MessageVO `json:"message" binding:"required"`
}
