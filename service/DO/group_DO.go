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
	// Avatar string `json:"avatar"`
}

type GroupDO struct {
	GroupID    int64
	GroupName  string
	UserID     int64
	Type       int
	CreateTime string
	Extra      *GroupExtra
}

type GroupExtra struct {
	// Avatar string `json:"avatar"`
}

func MGetGroupInfofromPO(info group_dao.GroupInfoPO) (*GroupInfoDO, error) {
	result := GroupInfoDO{
		GroupID:     info.GroupID,
		OwnerID:     info.OwnerID,
		GroupName:   info.GroupName,
		CreateTime:  info.CreateTime,
		IsDeleted:   info.IsDeleted,
		Description: info.Description,
	}

	var adminIds []int64
	if info.AdminIds != nil {
		err := json.Unmarshal([]byte(*info.AdminIds), &adminIds)
		if err != nil {
			fmt.Println("[MGetGroupInfofromPO] Unmarshal adminIds err is ", err.Error())
			return nil, err
		}
	}
	result.AdminIds = &adminIds

	var silenceList []int64
	if info.SilenceList != nil {
		err := json.Unmarshal([]byte(*info.SilenceList), &silenceList)
		if err != nil {
			fmt.Println("[MGetGroupInfofromPO] Unmarshal silenceList err is ", err.Error())
			return nil, err
		}
	}
	result.SilenceList = &silenceList

	var userIds []int64
	if info.UserIds != nil {
		err := json.Unmarshal([]byte(*info.UserIds), &userIds)
		if err != nil {
			fmt.Println("[MGetGroupInfofromPO] Unmarshal userIds err is ", err.Error())
			return nil, err
		}
	}
	result.UserIds = &userIds

	var extra GroupInfoExtra
	if info.Extra != nil {
		err := json.Unmarshal([]byte(*info.Extra), &extra)
		if err != nil {
			fmt.Println("[MGetGroupInfofromPO] Unmarshal extra err is ", err.Error())
			return nil, err
		}
	}
	result.Extra = &extra

	return &result, nil
}

func TurnGroupInfoPOfromDO(info GroupInfoDO) (*group_dao.GroupInfoPO, error) {
	result := group_dao.GroupInfoPO{
		GroupID:     info.GroupID,
		OwnerID:     info.OwnerID,
		GroupName:   info.GroupName,
		CreateTime:  info.CreateTime,
		IsDeleted:   info.IsDeleted,
		Description: info.Description,
	}

	var admin string
	if len(*info.AdminIds) != 0 {
		data, err := json.Marshal(*info.AdminIds)
		if err != nil {
			fmt.Println("[TurnGroupInfoPOfromDO] Marshal err is ", err.Error())
			return nil, err
		}
		admin = string(data)
		result.AdminIds = &admin
	}else{
		result.AdminIds = nil
	}
	

	var silence string
	if len(*info.SilenceList) != 0 {
		data, err := json.Marshal(*info.SilenceList)
		if err != nil {
			fmt.Println("[TurnGroupInfoPOfromDO] Marshal err is ", err.Error())
			return nil, err
		}
		silence = string(data)
		result.SilenceList = &silence
	}else{
		result.SilenceList = nil
	}
	

	var users string
	if len(*info.UserIds) != 0 {
		data, err := json.Marshal(*info.UserIds)
		if err != nil {
			fmt.Println("[TurnGroupInfoPOfromDO] Marshal err is ", err.Error())
			return nil, err
		}
		users = string(data)
		result.UserIds = &users
	}else{
		result.UserIds = nil
	}
	

	var extra string
	if info.Extra != nil {
		data, err := json.Marshal(*info.Extra)
		if err != nil {
			fmt.Println("[TurnGroupInfoPOfromDO] Marshal err is ", err.Error())
			return nil, err
		}
		extra = string(data)
		result.Extra = &extra
	}else{
		result.Extra = nil
	}

	return &result, nil
}

func MGetGroupDOfromPO(group group_dao.GroupPO) (*GroupDO, error) {
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

func TurnGroupPOfromDO(group GroupDO) (*group_dao.GroupPO, error) {
	var extra string
	if group.Extra != nil {
		data, err := json.Marshal(*group.Extra)
		if err != nil {
			fmt.Println("[TurnGroupPOfromDO], Marshal is err ", err.Error())
			return nil, err
		}
		extra = string(data)
	}

	return &group_dao.GroupPO{
		GroupID:    group.GroupID,
		GroupName:  group.GroupName,
		UserID:     group.UserID,
		Type:       group.Type,
		CreateTime: group.CreateTime,
		Extra:      &extra,
	}, nil
}

func MGetGroupListfromPOList(list []group_dao.GroupPO) (*[]GroupDO, error) {
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
