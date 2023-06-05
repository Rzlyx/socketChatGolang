package PO

type FriendCirclePO struct {
	NewsID     int64   `json:"news_id" db:"news_id"`
	SenderID   int64   `json:"sender_id" db:"sender_id"`
	News       *string `json:"news" db:"news"`
	Type       int     `json:"type" db:"type"`
	BlackList  *string `json:"black_list" db:"black_list"`
	WhiteList  *string `json:"white_list" db:"white_list"`
	CteateTime string  `json:"cteate_time" db:"cteate_time"`
	Likes      *string `json:"likes" db:"likes"`
	IsDeleted  bool    `json:"is_deleted" db:"is_deleted"`
	Extra      *string `json:"extra" db:"extra"`
}
