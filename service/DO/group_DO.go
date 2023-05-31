package DO

import (
	"dou_yin/dao/mysql/group_dao"
	"encoding/json"
	"fmt"
)

type GroupInfoDO struct {
	GroupID     int64
	OwnerID     int64
	AdminIds    *[]int64
	SilenceList *[]int64
	UserIds     *[]int64
	GroupName   string
	Description *string
	CreateTime  string
	IsDeleted   bool
	Extra       *GroupInfoExtra
}

type GroupInfoExtra struct {
	Avatar string `json:"avatar"`
}

type GroupDO struct {
	GroupID    int64      `json:"group_id"`
	GroupName  string     `json:"group_name"`
	UserID     int64      `json:"user_id"`
	Type       int        `json:"type"`
	CreateTime string     `json:"create_time"`
	Extra      *GroupExtra `json:"extra"`
}

type GroupExtra struct {
	// Avatar string `json:"avatar"`
}

func MGetGroupDOfromPO(group group_dao.GroupDO) (*GroupDO, error) {
	var extra GroupExtra
	if group.Extra != nil {
		err := json.Unmarshal([]byte(*group.Extra), &extra)
		if err != nil {
			fmt.Println("[MGetGroupDOfromPO] Unmarshal err is ", err.Error())
			return nil, err
		}
	}

	return &GroupDO{
		GroupID:    group.GroupID,
		GroupName:  group.GroupName,
		UserID:     group.UserID,
		Type:       group.Type,
		CreateTime: group.CreateTime,
		Extra:      &extra,
	}, nil
}

func MGetGroupListfromPOList(list []group_dao.GroupDO) (*[]GroupDO, error) {
	var result []GroupDO
	for _, group := range list {
		record, err := MGetGroupDOfromPO(group)
		if err != nil {
			if err != nil {
				fmt.Println("[MGetGroupListfromPOList] MGetGroupDOfromPO err is ", err.Error())
				return nil, err
			}
		}
		result = append(result, *record)
	}
	return &result, nil
}
