package PO

type FriendCirclePO struct {
	NewsID     int64   `json:"news_id" db:"news_id"`
	SenderID   int64   `json:"sender_id" db:"sender"`
	News       *string `json:"news" db:"news"`
	Type       int     `json:"type" db:"type"`
	BlackList  *string `json:"black_list" db:"blacklist"`
	WhiteList  *string `json:"white_list" db:"whitelist"`
	CreateTime string  `json:"create_time" db:"create_time"`
	Likes      *string `json:"likes" db:"likes"`
	IsDeleted  bool    `json:"is_deleted" db:"is_deleted"`
	Extra      *string `json:"extra" db:"extra"`
}
