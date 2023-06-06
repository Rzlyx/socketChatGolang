package response

type FriendCircleContext struct {
	NewsID     int64     `json:"news_id,string"`
	SenderID   int64     `json:"sender_id"`
	News       *string   `json:"news"`
	Type       int       `json:"type"`
	CreateTime string    `json:"create_time"`
	Likes      *[]string `json:"likes"`
}
