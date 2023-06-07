package response

type CreateCircle struct {
	NewsID int64 `json:"news_id,string"`
}
type FriendCircleContext struct {
	NewsID     int64      `json:"news_id,string"`
	SenderID   int64      `json:"sender_id,string"`
	SenderName string     `json:"sender_name"`
	News       *string    `json:"news"`
	Type       int        `json:"type"`
	CreateTime string     `json:"create_time"`
	Likes      *[]string  `json:"likes"`
	PhotoPaths *[]string  `json:"photo_paths"`
	Comments   *[]Comment `json:"comments"`
}

type Comment struct {
	SenderID   int64  `json:"sender_id"`
	Message    string `json:"message"`
	Caller     int64  `json:"caller"`
	CreateTime string `json:"create_time"`
}
