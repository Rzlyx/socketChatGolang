package PO

type ApplyPO struct {
	ApplyID    int64  `json:"apply_id" db:"apply_id"`
	Applicant  int64  `json:"applicant" db:"applicant"`
	TargetID   int64  `json:"target_id" db:"target_id"`
	Type       int    `json:"type" db:"type"`
	Status     int    `json:"status" db:"status"`
	Reason     string `json:"reason" db:"reason"`
	CreateTime string `json:"create_time" db:"create_time"`
	Extra      string `json:"extra" db:"extra"`
}
