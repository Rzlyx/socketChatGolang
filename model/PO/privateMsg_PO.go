package PO

type PrivateMsgPOf struct {
	ID        string `json:"id" db:"id"`
	Context   string `json:"context" db:"context"`
	Time      string `json:"time" db:"time"`
	SendID    string `json:"send_id" db:"send_id"`
	ReceiveID string `json:"receive_id" db:"receive_id"`
	Type      string `json:"type" db:"type"`
}

type PrivateMsgPO struct {
	MsgID        int64   `json:"message_id" db:"message_id"`
	FriendshipID int64   `json:"friendship_id" db:"friendship_id"`
	SenderID     int64   `json:"sender" db:"sender"`
	ReceiverID   int64   `json:"receiver" db:"receiver"`
	Message      string  `json:"message" db:"message"`
	Type         int     `json:"type" db:"type"`
	CreateTime   string  `json:"create_time" db:"create_time"`
	Deleted_list int     `json:"deleted_list" db:"deleted_list"`
	Extra        *string `json:"extra" db:"extra"`
}
