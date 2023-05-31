package PO

type FriendPO struct {
	FriendshipID        int64  `json:"friendship_id" db:"friendship_id"`
	FirstID             int64  `json:"first_id" db:"first_id"`
	SecondID            int64  `json:"second_id" db:"second_id"`
	FirstRemarkSecond   string `json:"f_remark_s" db:"f_remark_s"`
	SecondRemarkFirst   string `json:"s_remark_f" db:"s_remark_f"`
	IsFirstRemarkSecond bool   `json:"is_f_remark_s" db:"is_f_remark_s"`
	IsSecondRemarkFirst bool   `json:"is_s_remark_f" db:"is_s_remark_f"`
	Extra               string `json:"extra" db:"extra"`
}
