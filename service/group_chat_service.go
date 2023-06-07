package service

import (
	"dou_yin/dao/mysql/group_chat_dao"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/dao/redis"
	"dou_yin/model/VO"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"fmt"
	"strings"
)

type GroupListType int

const (
	GROUP_LIST_UKNOW GroupListType = 5  // 不清楚
	GROUP_WHITE_LIST GroupListType = 6  // 群聊处于白名单
	GROUP_GRAY_LIST  GroupListType = 7  // 群聊处于灰名单
	GROUP_BLACK_LIST GroupListType = 8  // 群聊处于黑名单
	GROUP_APPLY_TYPE GroupListType = 10 // 申请/邀请通知
)

// 群聊名单
func GroupMSGTypeRedis(UserID, GroupID string) (GroupListType, error) {
	// TODO: 加缓存
	var Type GroupListType
	key := UserID + GroupID
	value, err := redis.GetUserGroup(key)
	if err == nil && len(value) > 0 { // 查到缓存
		node, err := redis.TurnNodeFromNode(value[0])
		if err != nil {
			return GROUP_LIST_UKNOW, err
		}
		return GroupListType(node.Type), nil
	} else { // 没有缓存
		Type, err = GroupMSGType(UserID, GroupID)
		if err != nil {
			fmt.Println("[GroupMSGTypeRedis], GroupMSGType err is ", err.Error())
			return GROUP_LIST_UKNOW, err
		}
		IsSlience, err := IsGroupSlienceList(UserID, GroupID)
		if err == nil {
			AddRedisUserGroup(UserID, GroupID, int(Type), IsSlience)
		}

		return Type, nil
	}
}

// 群聊名单
func GroupMSGType(UserID, GroupID string) (GroupListType, error) {
	// TODO: 加缓存
	var Type GroupListType
	userInfoPO, err := user_dao.QueryUserInfo(utils.ShiftToNum64(UserID))
	if err != nil {
		return GROUP_LIST_UKNOW, err
	}
	whiteList, grayList, blackList, err := turnUserGroupList(userInfoPO)
	if err != nil {
		Type = GROUP_LIST_UKNOW
	}
	for _, id := range *whiteList {
		if id == utils.ShiftToNum64(GroupID) {
			Type = GROUP_WHITE_LIST
		}
	}
	for _, id := range *grayList {
		if id == utils.ShiftToNum64(GroupID) {
			Type = GROUP_GRAY_LIST
		}
	}
	for _, id := range *blackList {
		if id == utils.ShiftToNum64(GroupID) {
			Type = GROUP_BLACK_LIST
		}
	}

	return Type, nil

}

func AddRedisUserGroup(UserID, GroupID string, Type int, IsSlience bool) error {
	// 加入缓存
	key := UserID + GroupID
	value, _ := redis.TurnStringFromNode(&redis.UserGroup{
		UserId:    utils.ShiftToNum64(UserID),
		GroupID:   utils.ShiftToNum64(GroupID),
		Type:      int(Type),
		IsSlience: IsSlience,
	})
	err := redis.AddUserGroup(value, key)
	if err != nil {
		fmt.Println("[GroupMSGType], redis.AddUserGroup err is ", err.Error())
		return err
	}
	return nil
}

// 在禁言名单
func IsGroupSlienceListRedis(UserID, GroupID string) (bool, error) {
	// TODO: 加缓存
	key := UserID + GroupID
	value, err := redis.GetUserGroup(key)
	if err == nil && len(value) > 0 { // 查到缓存
		node, err := redis.TurnNodeFromNode(value[0])
		if err == nil {
			return node.IsSlience, nil
		}
	} else { // 无缓存
		IsSlience, err := IsGroupSlienceList(UserID, GroupID)
		Type, _ := GroupMSGType(UserID, GroupID)
		if err != nil {
			fmt.Println("[IsGroupSlienceListRedis], IsGroupSlienceList err is ", err.Error())
		}
		if err == nil {
			err := AddRedisUserGroup(UserID, GroupID, int(Type), IsSlience)
			fmt.Println("[IsGroupSlienceListRedis], AddRedisUserGroup err is ", err.Error())
		}
		return IsSlience, nil
	}

	return false, nil

}

func IsGroupSlienceList(UserID, GroupID string) (bool, error) {
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

// 获取历史消息by用户ID、群ID，pageNum， Num
func QueryGroupOldMsgList(UserID, GroupID string, pageNum, num int) (*[]VO.MessageVO, error) {
	groupPOList, err := group_dao.MGetGroupListByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupDOList, err := DO.MGetGroupListfromPOList(groupPOList)
	if err != nil {
		return nil, err
	}
	Names := make(map[int64]string)
	for _, group := range *groupDOList {
		Names[group.UserID] = group.Extra.MyName
	}

	var result []VO.MessageVO
	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(userID, groupID)
	if err != nil {
		return &result, err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return &result, err
	}

	list, err := group_chat_dao.MGetGroupMsgOldList(groupID, groupDO.Extra.ReadTime, pageNum, num)
	if err != nil {
		return &result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
			MSG.SenderName = Names[msgDO.SenderID]
			result = append(result, *MSG)
		}
	}

	return &result, nil
}

func QueryGroupOldMsgLogin(UserID, GroupID string) (*[]VO.MessageVO, error) {
	groupPOList, err := group_dao.MGetGroupListByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupDOList, err := DO.MGetGroupListfromPOList(groupPOList)
	if err != nil {
		return nil, err
	}
	Names := make(map[int64]string)
	for _, group := range *groupDOList {
		Names[group.UserID] = group.Extra.MyName
	}

	var result []VO.MessageVO

	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(userID, groupID)
	if err != nil {
		return &result, err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return &result, err
	}

	list, err := group_chat_dao.MGetGroupOldMsgListLogin(groupID, groupDO.Extra.ReadTime)
	if err != nil {
		return &result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
			MSG.SenderName = Names[msgDO.SenderID]
			result = append(result, *MSG)
		}
	}

	return &result, nil
}

func QueryGroupOldMsgUp(UserID, GroupID, TimeTag string) (*[]VO.MessageVO, error) {
	groupPOList, err := group_dao.MGetGroupListByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupDOList, err := DO.MGetGroupListfromPOList(groupPOList)
	if err != nil {
		return nil, err
	}
	Names := make(map[int64]string)
	for _, group := range *groupDOList {
		Names[group.UserID] = group.Extra.MyName
	}

	var result []VO.MessageVO

	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	list, err := group_chat_dao.MGetGroupOldMsgListLogin(groupID, TimeTag)
	if err != nil {
		return &result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
			MSG.SenderName = Names[msgDO.SenderID]
			result = append(result, *MSG)
		}
	}

	return &result, nil
}

func QueryGroupOldMsgDay(UserID, GroupID, StartTime, EndTime string) (*[]VO.MessageVO, error) {
	groupPOList, err := group_dao.MGetGroupListByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupDOList, err := DO.MGetGroupListfromPOList(groupPOList)
	if err != nil {
		return nil, err
	}
	Names := make(map[int64]string)
	for _, group := range *groupDOList {
		Names[group.UserID] = group.Extra.MyName
	}

	var result []VO.MessageVO

	groupID := utils.ShiftToNum64(GroupID)
	userID := utils.ShiftToNum64(UserID)

	Type, err := GroupMSGType(UserID, GroupID)
	if err != nil {
		return nil, err
	}

	list, err := group_chat_dao.MGetGroupOldMsgListDay(groupID, StartTime, EndTime)
	if err != nil {
		return &result, err
	}

	for _, msg := range *list {
		msgDO, err := DO.MGetGroupMsgDOfromPO(&msg)
		if err != nil {
			return nil, err
		}
		if !IsDeletedMsg(userID, *msgDO.DeletedList) {
			MSG := DO.MGetMsgVOfromDO(msgDO, int(Type))
			MSG.SenderName = Names[msgDO.SenderID]
			result = append(result, *MSG)
		}
	}

	return &result, nil
}

func SendGroupNewMsg(UserId string) error {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(UserId))
	if err != nil {
		return err
	}
	whiteList, grayList, _, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	if whiteList != nil {
		for _, groupID := range *whiteList {
			err := StartSendGroupNewMsg(UserId, utils.ShiftToStringFromInt64(groupID), int(GROUP_WHITE_LIST))
			if err != nil {
				return err
			}
		}
	}
	if grayList != nil {
		for _, groupID := range *grayList {
			err := StartSendGroupNewMsg(UserId, utils.ShiftToStringFromInt64(groupID), int(GROUP_GRAY_LIST))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取未读信息by用户ID、群ID
func QueryGroupNewMsgList(UserID, GroupID string) ([]VO.MessageVO, error) {
	groupPOList, err := group_dao.MGetGroupListByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupDOList, err := DO.MGetGroupListfromPOList(groupPOList)
	if err != nil {
		return nil, err
	}
	Names := make(map[int64]string)
	for _, group := range *groupDOList {
		Names[group.UserID] = group.Extra.MyName
	}

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

	list, err := group_chat_dao.MGetGroupNewList(groupID, groupDO.Extra.ReadTime)
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
				MSG.SenderName = Names[msgDO.SenderID]
				result = append(result, *MSG)
			}
		}
	}

	return result, nil
}

func HandleGroupChatMsg(msg *VO.MessageVO) {
	IsSlience, err := IsGroupSlienceListRedis(msg.SenderID, msg.ReceiverID)
	if err != nil {
		msg.ErrString = "系统内部错误，请稍后再试"
	}
	if IsSlience {
		msg.ErrString = "已被禁言"
	} else {
		err = CreatGroupMsg(*msg)
		if err != nil {
			msg.ErrString = "系统内部错误，请稍后再试"
		}
		fmt.Println("转发群聊消息 ", msg.SenderID,"-->", msg.ReceiverID," msg:", msg)
		MsgChan <- *msg
		if strings.HasPrefix(msg.Message, "@GPT") {
			result, err := GetGPTMessage(msg)
			if err != nil {
				msg.ErrString = "系统内部错误，请稍后再试"
			}
			fmt.Println("转发GPT群聊消息 ", msg.SenderID,"-->", msg.ReceiverID," msg:", msg)
			MsgChan <- *result
		}
	}
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
