package DO

type ContactList struct {
	ContactorList []ContactInfo `json:"contactor_list"`
}

type ContactInfo struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Message string `json:"message" binding:"required"`
	Time    string `json:"time" binding:"required"`
}
