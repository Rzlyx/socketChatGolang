package DO

type ContactList struct {
	ContactorList []ContactInfo `json:"contactor_list"`
}

type ContactInfo struct {
	ID           int64  `json:"id,string"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	FriendshipID int64  `json:"friendship_id,string"`
}
