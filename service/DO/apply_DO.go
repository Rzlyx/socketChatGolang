package DO

import (
	"dou_yin/model/PO"
	"encoding/json"
	"fmt"
)

type ApplyDO struct {
	ApplyID    int64       `json:"apply_id" db:"apply_id"`
	Applicant  int64       `json:"applicant" db:"applicant"`
	TargetID   int64       `json:"target_id" db:"target_id"`
	Type       int         `json:"type" db:"type"`
	Status     int         `json:"status" db:"status"`
	Reason     string      `json:"reason" db:"reason"`
	CreateTime string      `json:"create_time" db:"create_time"`
	Extra      *ApplyExtra `json:"extra" db:"extra"`
}

type ApplyExtra struct {
	InvitedID int64 `json:"invited_id"`
}


func MGetApplyDOFromPO(apply *PO.ApplyPO) (*ApplyDO, error) {
	var extra ApplyExtra
	if apply.Extra != nil {
		err := json.Unmarshal([]byte(*apply.Extra), &extra)
		if err != nil {
			fmt.Println("[MGetApplyDOFromPO], Unmarshal err is ", err.Error())
			return nil, err
		}
	}
	return &ApplyDO{
		ApplyID: apply.ApplyID,
		Applicant: apply.Applicant,
		TargetID: apply.TargetID,
		Type: apply.Type,
		Status: apply.Status,
		Reason: apply.Reason,
		CreateTime: apply.CreateTime,
		Extra: &extra,
	}, nil
}

func MGetApplyPOFromDO(apply *ApplyDO) (*PO.ApplyPO, error) {
	result := PO.ApplyPO{
		ApplyID: apply.ApplyID,
		Applicant: apply.Applicant,
		TargetID: apply.TargetID,
		Type: apply.Type,
		Status: apply.Status,
		Reason: apply.Reason,
		CreateTime: apply.CreateTime,
	}

	var extra string
	if apply.Extra != nil {
		data, err := json.Marshal(*apply.Extra)
		if err != nil {
			fmt.Println("[MGetApplyDOFromPO], Marshal err is ", err.Error())
			return nil, err
		}
		extra = string(data)
		result.Extra = &extra
	}else{
		result.Extra = nil
	}
	return &result, nil
}