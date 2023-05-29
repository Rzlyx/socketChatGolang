package model

type User struct {
	ID       int64  `json:"id" db:"id"`
	UserName string `json:"user_name" db:"username"`
	Password string `json:"password" db:"password"`
	EMail    string `json:"e_mail" db:"email"`
}
type Contactor struct {
	UserID   int64  `json:"id" db:"user_id"`
	UserName string `json:"user_name" db:"username"`
	Remark   string `json:"remark" db:"remark"`
}
type ContactorList struct {
	ContactorList []*Contactor `json:"contactor_list"`
}
