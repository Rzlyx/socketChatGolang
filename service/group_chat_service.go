package service

import (
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
)

type GroupListType int

const (
	GROUP_LIST_UKNOW GroupListType = 0	// 不清楚
	GROUP_WHITE_LIST GroupListType = 1 	// 群聊处于白名单
	GROUP_GRAY_LIST  GroupListType = 2 	// 群聊处于灰名单
	GROUP_BLACK_LIST GroupListType = 3	// 群聊处于黑名单
)

// 群聊是否在白名单
func IsGroupWhite(UserID, GroupID string) (GroupListType, error) {
	// TODO: 加缓存
	userInfoPO, err := user_dao.QueryUserInfo(utils.ShiftToNum64(UserID))
	if err != nil {
		return GROUP_LIST_UKNOW, err
	}
	whiteList, grayList, blackList, err := turnUserGroupList(userInfoPO)
	if err != nil {
		return GROUP_LIST_UKNOW, err
	}
	for _, id := range *whiteList {
		if id == utils.ShiftToNum64(GroupID) {
			return GROUP_WHITE_LIST, nil
		}
	}
	for _, id := range *grayList {
		if id == utils.ShiftToNum64(GroupID) {
			return GROUP_GRAY_LIST, nil
		}
	}
	for _, id := range *blackList {
		if id == utils.ShiftToNum64(GroupID) {
			return GROUP_BLACK_LIST, nil
		}
	}
	return GROUP_LIST_UKNOW, nil
}

// 在禁言名单
func IsGroupSlienceList(UserID, GroupID string) (bool, error) {
	// TODO: 加缓存
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return false, err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return false, err
	}
	for _, id := range *groupInfoDO.SilenceList {
		if id == utils.ShiftToNum64(UserID) {
			return true, nil
		}
	}
	return false, nil
}

// 获取历史消息by用户ID、群ID
func QueryGroupOldMsgList() {

}

// 获取未读信息by用户ID、群ID
func QueryGroupNewMsgList() {

}

// 接收新消息并保存
func CreatGroupMsg() error {

	return nil
}
