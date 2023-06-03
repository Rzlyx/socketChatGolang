package service

import (
	"dou_yin/dao/mysql/group_chat_dao"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/VO"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
)

type GroupListType int

const (
	GROUP_LIST_UKNOW GroupListType = 0 // 不清楚
	GROUP_WHITE_LIST GroupListType = 1 // 群聊处于白名单
	GROUP_GRAY_LIST  GroupListType = 2 // 群聊处于灰名单
	GROUP_BLACK_LIST GroupListType = 3 // 群聊处于黑名单
)

// 群聊是否在白名单
func GroupMSGType(UserID, GroupID string) (GroupListType, error) {
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

func IsDeletedMsg(UserID int64, List []int64) bool {
	for _, id := range List {
		if id == UserID {
			return true
		}
	}
	return false
}

// 获取历史消息by用户ID、群ID
func QueryGroupOldMsgList(UserID, GroupID string, pageNum, num int) ([]VO.MessageVO, error) {
	var result []VO.MessageVO
	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(userID, groupID)
	if err != nil {
		return result, err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return result, err
	}

	list, err := group_chat_dao.MGetGroupMsgOldList(groupID, groupDO.Extra.ReadTime, pageNum, num)
	if err != nil {
		return result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
			result = append(result, *MSG)
		}
	}

	return result, nil
}

// 获取未读信息by用户ID、群ID
func QueryGroupNewMsgList(UserID, GroupID string, pageNum, num int) ([]VO.MessageVO, error) {
	var result []VO.MessageVO
	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(userID, groupID)
	if err != nil {
		return result, err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return result, err
	}

	list, err := group_chat_dao.MGetGroupNewList(groupID, groupDO.Extra.ReadTime, pageNum, num)
	if err != nil {
		return result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			if Type == GROUP_WHITE_LIST || Type == GROUP_GRAY_LIST {
				MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
				result = append(result, *MSG)
			}
		}
	}

	return result, nil
}

// 接收新消息并保存
func CreatGroupMsg(msg VO.MessageVO) error {
	msgDO, err := DO.MGetMsgDOfromVO(&msg)
	if err != nil {
		return err
	}
	msgPO, err := DO.MGetGroupMsgPOfromDO(msgDO)
	if err != nil {
		return err
	}
	err = group_chat_dao.WriteGroupMsg(*msgPO)
	if err != nil {
		return err
	}

	return nil
}
