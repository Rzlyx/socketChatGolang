package group_dao

type GroupInfoPO struct {
	GroupID     int64   `json:"group_id" db:"group_id"`
	GroupName   string  `json:"group_name" db:"group_name"`
	Description *string `json:"description" db:"description"`
	UserIds     *string `json:"user_ids" db:"user_ids"`
	OwnerID     string  `json:"owner_id" db:"owner_id"`
	AdminIds    *string `json:"admin_ids" db:"admin_ids"`
	SilenceList *string `json:"silence_list" db:"silence_list"`
	CreateTime  string  `json:"create_time" db:"create_time"`
	IsDeleted   bool    `json:"is_deleted" db:"is_deleted"`
	Extra       *string `json:"extra" db:"extra"`
}

type GroupDO struct {
	GroupID    int64   `json:"group_id" db:"group_id"`
	GroupName  string  `json:"group_name" db:"group_name"`
	UserID     int64   `json:"user_id" db:"user_id"`
	Type       int     `json:"type" db:"type"`
	CreateTime string  `json:"create_time" db:"create_time"`
	Extra      *string `json:"extra" db:"extra"`
}
