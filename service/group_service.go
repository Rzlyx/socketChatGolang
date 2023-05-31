package service

import (
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/service/DO"
)

func MGetGroupListByParam(info *param.QueryGroupListParam) (*[]response.GroupJoin, error) {
	list, err := group_dao.MGetGroupListByUserID(info.UserID)
	if err != nil {
		return nil, err
	}

	groups, err :=  DO.MGetGroupListfromPOList(list)
	if err != nil {
		return nil, err
	}
	
	var result []response.GroupJoin
	for _, group := range *groups {
		result = append(result, response.GroupJoin{
			GroupID: group.GroupID,
			GroupName: group.GroupName,
			// Avatar: group.Extra,
		})
	}
	return &result, nil
}