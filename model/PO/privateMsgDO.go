package PO

type PrivateMsgDO struct {
	ID        string `json:"id" db:"id"`
	Context   string `json:"context" db:"context"`
	Time      string `json:"time" db:"time"`
	SendID    string `json:"send_id" db:"send_id"`
	ReceiveID string `json:"receive_id" db:"receive_id"`
	Type      string `json:"type" db:"type"`
}
