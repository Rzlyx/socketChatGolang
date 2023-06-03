package DO

type ContactList struct {
	ContactorList []ContactInfo `json:"contactor_list"`
}

type ContactInfo struct {
	ID      int64  `json:"id,string"`
	Name    string `json:"name"`
	Message string `json:"message"`
	Time    string `json:"time"`
}
