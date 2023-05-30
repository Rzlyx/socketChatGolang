package model

type User struct {
	UserID   int64  `json:"user_id" db:"user_id"`
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
	EMail    string `json:"e_mail" db:"e_mail"`
}
type Contactor struct {
	UserID   int64  `json:"id" db:"user_id"`
	UserName string `json:"user_name" db:"username"`
	Remark   string `json:"remark" db:"remark"`
}
type ContactorList struct {
	ContactorList []*Contactor `json:"contactor_list"`
}
