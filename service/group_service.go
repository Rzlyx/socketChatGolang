package service

import (
	"database/sql"
	"dou_yin/dao/mysql"
	"dou_yin/dao/mysql/apply_dao"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"fmt"
)

func MGetGroupInfoByParam(info *param.QueryGroupInfoParam) (*response.GroupInfo, error) {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return nil, err
	}
	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return nil, err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return nil, err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return nil, err
	}
	userIds := *groupInfoDO.AdminIds
	userIds = append(userIds, *groupInfoDO.UserIds...)
	userIds = append(userIds, groupInfoDO.OwnerID)
	OnlineNum := 0
	for _, id := range userIds {
		if _, ok := UserChan[id]; ok {
			OnlineNum++
		}
	}

	result := response.GroupInfo{
		GroupID:     groupInfoDO.GroupID,
		OwnerID:     groupInfoDO.OwnerID,
		AdminIds:    *groupInfoDO.AdminIds,
		SilenceList: *groupInfoDO.SilenceList,
		UserIds:     *groupInfoDO.UserIds,
		GroupName:   groupInfoDO.GroupName,
		Description: *groupInfoDO.Description,
		CreateTime:  groupInfoDO.CreateTime,
		IsDeleted:   groupInfoDO.IsDeleted,
		Extra:       *groupInfoDO.Extra,
	}
	if groupDO.Extra.IsRemark {
		result.MyGroupName = groupDO.GroupName
	} else {
		result.MyGroupName = "无群备注"
	}
	result.Type = groupDO.Type
	result.MyName = groupDO.Extra.MyName
	msgType, err := GroupMSGType(info.UserID, info.GroupID)
	if err != nil {
		return nil, err
	}
	result.MsgType = int(msgType)
	result.OnlineNum = OnlineNum
	return &result, nil
}

func MGetGroupListByParam(info *param.QueryGroupListParam) (*[]response.GroupJoin, error) {
	list, err := group_dao.MGetGroupListByUserID(utils.ShiftToNum64(info.UserID))
	if err != nil {
		return nil, err
	}

	groups, err := DO.MGetGroupListfromPOList(list)
	if err != nil {
		return nil, err
	}

	var result []response.GroupJoin
	for _, group := range *groups {
		result = append(result, response.GroupJoin{
			GroupID:   group.GroupID,
			GroupName: group.GroupName,
			// Avatar: group.Extra,
		})
	}
	return &result, nil
}

func GetGroupAllUserbyParam(info *param.GetGroupAllUserParam) (*[]response.GroupUserInfo, error) {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return nil, err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return nil, err
	}
	userIds := *groupInfoDO.AdminIds
	userIds = append(userIds, *groupInfoDO.UserIds...)
	userIds = append(userIds, groupInfoDO.OwnerID)

	var List []response.GroupUserInfo

	for _, id := range userIds {
		groupPO, err := group_dao.MGetGroupByUserIDandGroupID(id, groupInfoDO.GroupID)
		if err != nil {
			return nil, err
		}
		groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
		if err != nil {
			return nil, err
		}
		status := 0
		if _, ok := UserChan[id]; ok {
			status = 1
		}
		isSlienced := false
		for _, id := range *groupInfoDO.SilenceList {
			if id == utils.ShiftToNum64(info.UserID) {
				isSlienced = true
			}
		}
		List = append(List, response.GroupUserInfo{
			UserID:     id,
			MyName:     groupDO.Extra.MyName,
			InsertTime: groupDO.CreateTime,
			Status:     status,
			Type:       groupDO.Type,
			IsSlienced: isSlienced,
		})
	}

	return &List, nil
}

// 创建群聊
func CreateGroupInfoByParam(info *param.CreateGroupInfoParam) error {
	var groupInfo DO.GroupInfoDO
	groupInfo.GroupID = snowflake.GenID()
	groupInfo.OwnerID = utils.ShiftToNum64(info.OwnerID)
	groupInfo.AdminIds = new([]int64)
	groupInfo.Description = info.Description
	groupInfo.CreateTime = utils.GetNowTime()
	groupInfo.Extra = &DO.GroupInfoExtra{
		AIGPT: true,
	}
	groupInfo.GroupName = info.GroupName
	groupInfo.IsDeleted = false
	groupInfo.SilenceList = new([]int64)
	// 将好友拉入群聊
	var userIds []int64
	if info.UserIDs != nil && len(*info.UserIDs) > 0 {
		for _, id := range *info.UserIDs {
			userIds = append(userIds, utils.ShiftToNum64(id))
		}
	}
	userIds = append(userIds, 999999) // AI机器人
	groupInfo.UserIds = &userIds

	groupInfoPO, err := DO.TurnGroupInfoPOfromDO(groupInfo)
	if err != nil {
		return err
	}

	var users []PO.UserPO
	var UserIDs []int64
	if info.UserIDs != nil && len(*info.UserIDs) > 0 {
		for _, id := range *info.UserIDs {
			UserIDs = append(UserIDs, utils.ShiftToNum64(id))
		}
	}
	UserIDs = append(UserIDs, utils.ShiftToNum64(info.OwnerID))
	UserIDs = append(UserIDs, 999999) // AI机器人
	var groups []group_dao.GroupPO
	for _, id := range UserIDs {
		userInfo, err := user_dao.QueryUserInfo(id)
		if err != nil {
			return err
		}

		var groupWhiteList []int64
		if userInfo.GroupChatWhite != nil {
			err := json.Unmarshal([]byte(*userInfo.GroupChatWhite), &groupWhiteList)
			if err != nil {
				fmt.Println("[CreateGroupInfoByParam], Unmarshal err is ", err.Error())
				return nil
			}
		}
		groupWhiteList = append(groupWhiteList, groupInfo.GroupID)

		var whiteList string
		data, err := json.Marshal(groupWhiteList)
		if err != nil {
			return err
		}
		whiteList = string(data)

		userInfo.GroupChatWhite = &whiteList

		users = append(users, userInfo)

		var groupDO DO.GroupDO
		groupDO.CreateTime = utils.GetNowTime()
		groupDO.Extra = new(DO.GroupExtra)
		groupDO.GroupID = groupInfo.GroupID
		groupDO.GroupName = groupInfoPO.GroupName
		groupDO.Type = 0
		groupDO.UserID = id
		var extra DO.GroupExtra
		extra.ReadTime = utils.GetNowTime()
		extra.IsRemark = false
		extra.MyName = userInfo.UserName
		groupDO.Extra = &extra

		groupPO, err := DO.TurnGroupPOfromDO(groupDO)
		if err != nil {
			return err
		}

		groups = append(groups, *groupPO)
	}

	// TODO: Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		// 创建群
		err = group_dao.CreateGroupInfo(tx, groupInfoPO)
		if err != nil {
			return err
		}
		// 修改用户白名单
		for _, user := range users {
			err = user_dao.UpdateUserInfoByPOTx(tx, &user)
			if err != nil {
				return err
			}
		}
		// 添加群关系信息
		for _, group := range groups {
			_, err = group_dao.CreateGroupByGroupPO(tx, group)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 修改群信息
func UpdateGroupInfoByParam(info *param.UpdateGroupInfoParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	groupInfoDO.GroupName = info.GroupName
	groupInfoDO.Description = info.Description

	var userIDs []int64
	userIDs = append(userIDs, groupInfoDO.OwnerID)
	userIDs = append(userIDs, *groupInfoDO.AdminIds...)
	userIDs = append(userIDs, *groupInfoDO.UserIds...)

	var groupPOs []group_dao.GroupPO

	for _, id := range userIDs {
		groupPO, err := group_dao.MGetGroupByUserIDandGroupID(id, groupInfoDO.GroupID)
		if err != nil {
			return err
		}
		groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
		if err != nil {
			return err
		}
		if !groupDO.Extra.IsRemark { // 未设有备注
			groupDO.GroupName = info.GroupName
			GroupPO, err := DO.TurnGroupPOfromDO(*groupDO)
			if err != nil {
				return err
			}
			groupPOs = append(groupPOs, *GroupPO)
		}
	}

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, GroupInfoPO)
		if err != nil {
			return err
		}
		for _, GroupPO := range groupPOs {
			_, err = group_dao.UpdateGroupByGroupPO(tx, GroupPO)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 解散群聊
func DissolveGroupInfoByParam(info *param.DissolveGroupInfoParam) error {
	// TODO 加锁
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}

	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}
	userIds := *groupInfoDO.AdminIds
	userIds = append(userIds, *groupInfoDO.UserIds...)
	userIds = append(userIds, groupInfoDO.OwnerID)

	groupInfoPO.IsDeleted = true

	var UserPOList []PO.UserPO
	for _, id := range userIds {
		userPO, err := user_dao.QueryUserInfo(id)
		if err != nil {
			return err
		}
		whiteList, grayList, blackList, err := turnUserGroupList(userPO)
		if err != nil {
			return err
		}

		var white []int64
		for _, ID := range *whiteList {
			if ID != utils.ShiftToNum64(info.GroupID) {
				white = append(white, ID)
			}
		}

		var gray []int64
		for _, ID := range *grayList {
			if ID != utils.ShiftToNum64(info.GroupID) {
				gray = append(gray, ID)
			}
		}

		var black []int64
		for _, ID := range *blackList {
			if ID != utils.ShiftToNum64(info.GroupID) {
				black = append(black, ID)
			}
		}

		str1, str2, str3, err := turnjsonList(white, gray, black)
		if err != nil {
			return err
		}
		userPO.GroupChatWhite = str1
		userPO.GroupChatGray = str2
		userPO.GroupChatBlack = str3
		if len(white) == 0 {
			userPO.GroupChatWhite = nil
		}
		if len(gray) == 0 {
			userPO.GroupChatGray = nil
		}
		if len(black) == 0 {
			userPO.GroupChatBlack = nil
		}
		UserPOList = append(UserPOList, userPO)
	}

	// TODO:Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, groupInfoPO)
		if err != nil {
			return err
		}

		// 修改群聊名单
		for _, userPO := range UserPOList {
			err := user_dao.UpdateUserInfoByPOTx(tx, &userPO)
			if err != nil {
				return err
			}
		}

		// 删除加群关系
		for _, id := range userIds {
			_, err := group_dao.DeleteGroupByUserIDandGroupID(tx, id, utils.ShiftToNum64(info.GroupID))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ApplyJoinGroupByParam(info *param.ApplyJoinGroupParam) error {
	application, err := apply_dao.MGetApplicationByGroupIDandUserID(utils.ShiftToNum64(info.GroupID), utils.ShiftToNum64(info.UserID))
	if err != nil {
		return err
	}
	if application != nil {
		fmt.Println("[ApplyJoinGroupByParam], apply had")
		return nil
	}

	_, err = group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}

	var apply PO.ApplyPO
	apply.Applicant = utils.ShiftToNum64(info.UserID)
	apply.ApplyID = snowflake.GenID()
	apply.CreateTime = utils.GetNowTime()
	apply.Extra = nil
	apply.Reason = info.Reason
	apply.TargetID = utils.ShiftToNum64(info.GroupID)
	apply.Status = 0
	apply.Type = 0
	// TODO:Tx
	err = apply_dao.CreateApplication(&apply)
	if err != nil {
		return err
	}

	// 通知群主和管理员处理
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}
	userIds := *groupInfoDO.AdminIds
	userIds = append(userIds, groupInfoDO.OwnerID)
	for _, id := range userIds {
		msg := VO.MessageVO{
			MsgType:    11,
			ReceiverID: utils.ShiftToStringFromInt64(id),
		}
		HandlePrivateChatMsg(msg) //发送通知
	}

	return nil
}

// 退出群聊
func QuitGroupByParam(info *param.QuitGroupParam) error {
	ret, err := group_dao.IsGroupUser(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err // 查询失败
	}
	if ret {
		userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(info.UserID))
		if err != nil {
			return err
		}
		whiteList, gratList, blackList, err := turnUserGroupList(userInfo)
		if err != nil {
			return err
		}

		var white, gray, black []int64
		for _, id := range *whiteList {
			if id != utils.ShiftToNum64(info.GroupID) {
				white = append(white, id)
			}
		}
		for _, id := range *gratList {
			if id != utils.ShiftToNum64(info.GroupID) {
				gray = append(gray, id)
			}
		}
		for _, id := range *blackList {
			if id != utils.ShiftToNum64(info.GroupID) {
				black = append(black, id)
			}
		}

		str1, str2, str3, err := turnjsonList(white, gray, black)
		if err != nil {
			return err
		}

		if len(white) > 0 {
			userInfo.GroupChatWhite = str1
		} else {
			userInfo.GroupChatWhite = nil
		}
		if len(gray) > 0 {
			userInfo.GroupChatGray = str2
		} else {
			userInfo.GroupChatGray = nil
		}
		if len(black) > 0 {
			userInfo.GroupChatBlack = str3
		} else {
			userInfo.GroupChatBlack = nil
		}

		groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
		if err != nil {
			return err
		}
		groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
		if err != nil {
			return err
		}

		var userIds []int64
		if groupInfoDO.UserIds != nil {
			for _, id := range *groupInfoDO.UserIds {
				if id != utils.ShiftToNum64(info.UserID) {
					userIds = append(userIds, id)
				}
			}
		}
		groupInfoDO.UserIds = &userIds
		var adminIds []int64
		if groupInfoDO.AdminIds != nil {
			for _, id := range *groupInfoDO.AdminIds {
				if id != utils.ShiftToNum64(info.UserID) {
					adminIds = append(adminIds, id)
				}
			}
		}
		groupInfoDO.AdminIds = &adminIds

		GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
		if err != nil {
			return err
		}
		// TODO: Tx
		err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
			// 在该群聊中，执行退群操作
			ret, err := group_dao.DeleteGroupByUserIDandGroupID(tx, utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
			if err != nil || !ret {
				return err // 删除失败
			}

			err = user_dao.UpdateUserInfoByPOTx(tx, &userInfo)
			if err != nil {
				return err
			}

			err = group_dao.UpdateGroupInfo(tx, GroupInfoPO)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil // 退群成功
}

// 查看加群申请
func QueryGroupApplyListByParam(info *param.QueryGroupApplyListParam) (*[]response.GroupJoinApply, error) {
	list, err := apply_dao.MGetApplicationListByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return nil, err
	}

	var result []response.GroupJoinApply
	for _, group := range *list {
		result = append(result, response.GroupJoinApply{
			ApplyID:   group.ApplyID,
			Applicant: group.Applicant,
			Reason:    group.Reason,
			TargetID:  group.TargetID,
		})
	}
	return &result, nil
}

// 同意加入
func AgreeGroupApplyByParam(info *param.AgreeGroupApplyParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(info.Applicant))
	if err != nil {
		return err
	}

	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	// 加入group_info普通成员列表
	var userIds []int64
	if groupInfoPO.UserIds != nil {
		userIds = append(userIds, *groupInfoDO.UserIds...)
	}
	userIds = append(*groupInfoDO.UserIds, utils.ShiftToNum64(info.Applicant))
	groupInfoDO.UserIds = &userIds
	groupInfo, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	// 加入用户白名单
	var groupWhiteList []int64
	if userInfo.GroupChatWhite != nil {
		err := json.Unmarshal([]byte(*userInfo.GroupChatWhite), &groupWhiteList)
		if err != nil {
			fmt.Println("[CreateGroupInfoByParam], Unmarshal GroupChatWhite err is ", err.Error())
			return nil
		}
	}
	groupWhiteList = append(groupWhiteList, utils.ShiftToNum64(info.GroupID))

	var whiteList string
	data, err := json.Marshal(groupWhiteList)
	if err != nil {
		return err
	}
	whiteList = string(data)

	userInfo.GroupChatWhite = &whiteList

	var groupDO DO.GroupDO
	groupDO.CreateTime = utils.GetNowTime()
	groupDO.Extra = new(DO.GroupExtra)
	groupDO.GroupID = utils.ShiftToNum64(info.GroupID)
	groupDO.GroupName = groupInfoPO.GroupName
	groupDO.Type = 0
	groupDO.UserID = utils.ShiftToNum64(info.Applicant)
	var extra DO.GroupExtra
	extra.ReadTime = utils.GetNowTime()
	extra.IsRemark = false
	extra.MyName = userInfo.UserName
	groupDO.Extra = &extra

	groupPO, err := DO.TurnGroupPOfromDO(groupDO)
	if err != nil {
		return err
	}

	// TODO:Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, groupInfo)
		if err != nil {
			return err
		}

		_, err = group_dao.CreateGroupByGroupPO(tx, *groupPO)
		if err != nil {
			return err
		}

		err = user_dao.UpdateUserInfoByPOTx(tx, &userInfo)
		if err != nil {
			return err
		}

		err = apply_dao.DeleteApplicationByApplyID(tx, utils.ShiftToNum64(info.ApplyID))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 不同意加入
func DisAgreeGroupApplyByParam(info *param.DisAgreeGroupApplyParam) error {
	err := mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err := apply_dao.DeleteApplicationByApplyID(tx, utils.ShiftToNum64(info.ApplyID))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 禁言
func SilenceByParam(info *param.SilenceParam) error {
	// TODO: 加锁
	group, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}

	groupInfo, err := DO.MGetGroupInfofromPO(*group)
	if err != nil {
		return err
	}
	var silenceList []int64
	if groupInfo.SilenceList != nil {
		silenceList = append(silenceList, *groupInfo.SilenceList...)
	}
	silenceList = append(silenceList, utils.ShiftToNum64(info.TargetID))
	data, err := json.Marshal(silenceList)
	if err != nil {
		fmt.Println("[SilenceByParam], Marshal err is ", err.Error())
		return err
	}
	silence := string(data)
	group.SilenceList = &silence

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, group)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 解除禁言
func UnSilenceByParam(info *param.UnSilenceParam) error {
	// TODO: 加锁
	group, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}

	groupInfo, err := DO.MGetGroupInfofromPO(*group)
	if err != nil {
		return err
	}

	var silenceList []int64
	for _, id := range *groupInfo.SilenceList {
		if id != utils.ShiftToNum64(info.TargetID) {
			silenceList = append(silenceList, id)
		}
	}
	if len(silenceList) > 0 {
		data, err := json.Marshal(silenceList)
		if err != nil {
			fmt.Println("[UnSilenceByParam], Marshal err is ", err.Error())
			return err
		}
		silence := string(data)
		group.SilenceList = &silence
	} else {
		group.SilenceList = nil
	}

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, group)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func TransferGroupByParam(info *param.TransferGroupParam) error {
	groupInfoRecord, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	OwnerRecord, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	TargetRecord, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.TargetID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}

	GroupInfo, err := DO.MGetGroupInfofromPO(*groupInfoRecord)
	if err != nil {
		return err
	}

	// targte是普通成员
	var userList []int64
	for _, id := range *GroupInfo.UserIds {
		if TargetRecord.Type != 0 || (TargetRecord.Type == 0 && id != utils.ShiftToNum64(info.TargetID)) {
			userList = append(userList, id)
		}
	}
	userList = append(userList, utils.ShiftToNum64(info.UserID)) // 群主变为普通成员
	data, err := json.Marshal(userList)
	if err != nil {
		fmt.Println("[TransferGroupByParam], Marshal err is ", err.Error())
		return err
	}
	users := string(data)
	groupInfoRecord.UserIds = &users

	// target是管理员
	if TargetRecord.Type == 1 {
		var adminList []int64
		for _, id := range *GroupInfo.AdminIds {
			if id != utils.ShiftToNum64(info.TargetID) {
				adminList = append(adminList, id)
			}
		}
		data, err := json.Marshal(adminList)
		if err != nil {
			fmt.Println("[TransferGroupByParam], Marshal err is ", err.Error())
			return err
		}
		admin := string(data)
		if len(adminList) > 0 {
			groupInfoRecord.AdminIds = &admin
		} else {
			groupInfoRecord.AdminIds = nil
		}
	}
	groupInfoRecord.OwnerID = utils.ShiftToNum64(info.TargetID)

	OwnerRecord.Type = 0  // 成为普通成员
	TargetRecord.Type = 2 // 成为群主

	// TODO:Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, groupInfoRecord)
		if err != nil {
			return err
		}
		_, err = group_dao.UpdateGroupByGroupPO(tx, *OwnerRecord)
		if err != nil {
			return err
		}
		_, err = group_dao.UpdateGroupByGroupPO(tx, *TargetRecord)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func SetBlackListByParam(info *param.SetBlackListParam) error {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(info.UserID))
	if err != nil {
		return err
	}

	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var whiteList, grayList []int64
	for _, id := range *groupWhiteList {
		if id != utils.ShiftToNum64(info.GroupID) {
			whiteList = append(whiteList, id)
		}
	}

	for _, id := range *groupGrayList {
		if id != utils.ShiftToNum64(info.GroupID) {
			grayList = append(grayList, id)
		}
	}

	blackList := append(*groupBlackList, utils.ShiftToNum64(info.GroupID))

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	if len(*white) > 0 {
		userInfo.GroupChatWhite = white
	} else {
		userInfo.GroupChatWhite = nil
	}
	if len(*gray) > 0 {
		userInfo.GroupChatGray = gray
	} else {
		userInfo.GroupChatGray = nil
	}
	if len(*black) > 0 {
		userInfo.GroupChatBlack = black
	} else {
		userInfo.GroupChatBlack = nil
	}
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPOTx(tx, &userInfo)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func SetGrayListByParam(info *param.SetGrayListParam) error {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(info.UserID))
	if err != nil {
		return err
	}
	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var whiteList, blackList []int64
	for _, id := range *groupWhiteList {
		if id != utils.ShiftToNum64(info.GroupID) {
			whiteList = append(whiteList, id)
		}
	}

	for _, id := range *groupBlackList {
		if id != utils.ShiftToNum64(info.GroupID) {
			blackList = append(blackList, id)
		}
	}

	grayList := append(*groupGrayList, utils.ShiftToNum64(info.GroupID))

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	if len(*white) > 0 {
		userInfo.GroupChatWhite = white
	} else {
		userInfo.GroupChatWhite = nil
	}
	if len(*gray) > 0 {
		userInfo.GroupChatGray = gray
	} else {
		userInfo.GroupChatGray = nil
	}
	if len(*black) > 0 {
		userInfo.GroupChatBlack = black
	} else {
		userInfo.GroupChatBlack = nil
	}
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPOTx(tx, &userInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func SetWhiteListByParam(info *param.SetWhiteListParam) error {
	userInfo, err := user_dao.QueryUserInfo(utils.ShiftToNum64(info.UserID))
	if err != nil {
		return err
	}

	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var grayList, blackList []int64
	for _, id := range *groupGrayList {
		if id != utils.ShiftToNum64(info.GroupID) {
			grayList = append(grayList, id)
		}
	}

	for _, id := range *groupBlackList {
		if id != utils.ShiftToNum64(info.GroupID) {
			blackList = append(blackList, id)
		}
	}

	whiteList := append(*groupWhiteList, utils.ShiftToNum64(info.GroupID))

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	if len(*white) > 0 {
		userInfo.GroupChatWhite = white
	} else {
		userInfo.GroupChatWhite = nil
	}
	if len(*gray) > 0 {
		userInfo.GroupChatGray = gray
	} else {
		userInfo.GroupChatGray = nil
	}
	if len(*black) > 0 {
		userInfo.GroupChatBlack = black
	} else {
		userInfo.GroupChatBlack = nil
	}
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = user_dao.UpdateUserInfoByPOTx(tx, &userInfo)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func turnUserGroupList(userInfo PO.UserPO) (*[]int64, *[]int64, *[]int64, error) {
	var groupWhiteList []int64
	if userInfo.GroupChatWhite != nil {
		err := json.Unmarshal([]byte(*userInfo.GroupChatWhite), &groupWhiteList)
		if err != nil {
			fmt.Println("[turnUserGroupList], Unmarshal err is ", err.Error())
			return nil, nil, nil, err
		}
	}

	var groupGrayList []int64
	if userInfo.GroupChatGray != nil {
		err := json.Unmarshal([]byte(*userInfo.GroupChatGray), &groupGrayList)
		if err != nil {
			fmt.Println("[turnUserGroupList], Unmarshal err is ", err.Error())
			return nil, nil, nil, err
		}
	}

	var groupBlackList []int64
	if userInfo.GroupChatBlack != nil {
		err := json.Unmarshal([]byte(*userInfo.GroupChatBlack), &groupBlackList)
		if err != nil {
			fmt.Println("[turnUserGroupList], Unmarshal err is ", err.Error())
			return nil, nil, nil, err
		}
	}

	return &groupWhiteList, &groupGrayList, &groupBlackList, nil
}

func turnjsonList(whiteList, grayList, blackList []int64) (*string, *string, *string, error) {
	var white, gray, black string
	data, err := json.Marshal(whiteList)
	if err != nil {
		fmt.Println("[turnjsonList], Marshal err is ", err.Error())
		return nil, nil, nil, err
	}
	white = string(data)

	data, err = json.Marshal(grayList)
	if err != nil {
		fmt.Println("[turnjsonList], Marshal err is ", err.Error())
		return nil, nil, nil, err
	}
	gray = string(data)

	data, err = json.Marshal(blackList)
	if err != nil {
		fmt.Println("[turnjsonList], Marshal err is ", err.Error())
		return nil, nil, nil, err
	}
	black = string(data)

	return &white, &gray, &black, nil
}

func SetGroupAdminByParam(info *param.SetGroupAdminParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	group, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.TargetID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	var userIDs []int64
	for _, id := range *groupInfoDO.UserIds {
		if id != utils.ShiftToNum64(info.TargetID) {
			userIDs = append(userIDs, id)
		}
	}
	groupInfoDO.UserIds = &userIDs

	adminIDs := append(*groupInfoDO.AdminIds, utils.ShiftToNum64(info.TargetID))
	groupInfoDO.AdminIds = &adminIDs

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	group.Type = 1

	// TODO:Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, GroupInfoPO)
		if err != nil {
			return err
		}
		_, err = group_dao.UpdateGroupByGroupPO(tx, *group)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func SetGroupUserByParam(info *param.SetGroupUserParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	group, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.TargetID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	var AdminIds []int64
	for _, id := range *groupInfoDO.AdminIds {
		if id != utils.ShiftToNum64(info.TargetID) {
			AdminIds = append(AdminIds, id)
		}
	}
	groupInfoDO.AdminIds = &AdminIds

	userIds := append(*groupInfoDO.UserIds, utils.ShiftToNum64(info.TargetID))
	groupInfoDO.UserIds = &userIds

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	group.Type = 0

	// TODO:Tx
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, GroupInfoPO)
		if err != nil {
			return err
		}
		_, err = group_dao.UpdateGroupByGroupPO(tx, *group)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 邀请加入群聊
func InviteJoinGroupByParam(info *param.InviteJoinGroupParam) error {
	applyDO := DO.ApplyDO{
		ApplyID:    snowflake.GenID(),
		Applicant:  utils.ShiftToNum64(info.TargetID),
		TargetID:   utils.ShiftToNum64(info.GroupID),
		Type:       2,
		Status:     0,
		Reason:     info.Reason,
		CreateTime: utils.GetNowTime(),
		Extra: &DO.ApplyExtra{
			InvitedID: utils.ShiftToNum64(info.UserID),
		},
	}
	applyPO, err := DO.MGetApplyPOFromDO(&applyDO)
	if err != nil {
		return err
	}

	err = apply_dao.CreateApplication(applyPO)
	if err != nil {
		return err
	}

	msg := VO.MessageVO{
		MsgType:    11,
		ReceiverID: utils.ShiftToStringFromInt64(utils.ShiftToNum64(info.TargetID)),
	}
	HandlePrivateChatMsg(msg) //发送通知

	return nil
}

func QueryInviteGroupByParam(info *param.QueryInviteGroupParam) (*[]response.InviteGroupInfo, error) {
	var result []response.InviteGroupInfo

	list, err := apply_dao.MGetApplicationListByUserID(utils.ShiftToNum64(info.UserID))
	if err != nil {
		return nil, err
	}
	if list != nil && len(*list) > 0 {
		for _, apply := range *list {
			applyDO, err := DO.MGetApplyDOFromPO(&apply)
			if err != nil {
				return nil, err
			}
			resp := response.InviteGroupInfo{
				ApplyID:   applyDO.ApplyID,
				InvitedID: applyDO.Extra.InvitedID,
				Applicant: applyDO.Applicant,
				TargetID:  applyDO.TargetID,
				Reason:    applyDO.Reason,
			}
			result = append(result, resp)
		}
	}

	return &result, nil
}

func AgreeInviteGroupByParam(info *param.AgreeInviteGroupParam) error {
	Info := param.AgreeGroupApplyParam{
		ApplyID:   info.ApplyID,
		Applicant: info.Applicant,
		GroupID:   info.GroupID,
		UserID:    info.UserID,
	}
	err := AgreeGroupApplyByParam(&Info)
	if err != nil {
		return err
	}

	return nil
}

func DisAgreeInviteGroupByParam(info *param.DisAgreeInviteGroupParam) error {
	err := mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err := apply_dao.DeleteApplicationByApplyID(tx, utils.ShiftToNum64(info.ApplyID))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 设置群备注
func SetGroupNameByParam(info *param.SetGroupNameParam) error {
	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return err
	}

	groupDO.GroupName = info.GroupName
	groupDO.Extra.IsRemark = true

	GroupPO, err := DO.TurnGroupPOfromDO(*groupDO)
	if err != nil {
		return err
	}
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		_, err = group_dao.UpdateGroupByGroupPO(tx, *GroupPO)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 设置在本群的昵称
func SetMyNamebyParam(info *param.SetMyNameParam) error {
	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return err
	}
	groupDO.Extra.MyName = info.MyName
	GroupPO, err := DO.TurnGroupPOfromDO(*groupDO)
	if err != nil {
		return err
	}
	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		_, err = group_dao.UpdateGroupByGroupPO(tx, *GroupPO)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// 设置已读时间
func SetGroupReadTimebyParam(info *param.SetGroupReadTimeParam) error {
	groupPO, err := group_dao.MGetGroupByUserIDandGroupID(utils.ShiftToNum64(info.UserID), utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupDO, err := DO.MGetGroupDOfromPO(*groupPO)
	if err != nil {
		return err
	}

	groupDO.Extra.ReadTime = utils.GetNowTime()

	GroupPO, err := DO.TurnGroupPOfromDO(*groupDO)
	if err != nil {
		return err
	}

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		_, err = group_dao.UpdateGroupByGroupPO(tx, *GroupPO)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func SetAIGPTbyParam(info *param.SetAIGPTParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(info.GroupID))
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}
	groupInfoDO.Extra.AIGPT = !groupInfoDO.Extra.AIGPT

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	err = mysql.Tx(mysql.DB, func(tx *sql.Tx) error {
		err = group_dao.UpdateGroupInfo(tx, GroupInfoPO)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func GetPageOldMsgByParam(info *param.GetPageOldMsgParam) (*[]VO.MessageVO, error) {
	return QueryGroupOldMsgList(info.UserID, info.GroupID, info.PageNum, info.Num)
}

// 登录获取历史消息--获取15条消息
func GetGroupOldMsgLoginbyParam(info *param.GetGroupOldMsgLoginParam) (*[]VO.MessageVO, error) {
	return QueryGroupOldMsgLogin(info.UserID, info.GroupID)
}

// 向上滑动刷新消息
func GetGroupOldMsgUpbyParam(info *param.GetGroupOldMsgUpParam) (*[]VO.MessageVO, error) {
	return QueryGroupOldMsgUp(info.UserID, info.GroupID, info.TimeTag)
}

// 根据时间（天）获取消息
func GetGroupOldMsgDaybyParam(info *param.GetGroupOldMsgDayParam) (*[]VO.MessageVO, error) {
	return QueryGroupOldMsgDay(info.UserID, info.GroupID, info.StartTime, info.EndTime)
}

func GetAllUserIDsbyGroupID(GroupID string) (*[]int64, error) {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(utils.ShiftToNum64(GroupID))
	if err != nil {
		return nil, err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return nil, err
	}
	var result []int64
	result = append(result, groupInfoDO.OwnerID)
	result = append(result, *groupInfoDO.AdminIds...)
	result = append(result, *groupInfoDO.UserIds...)

	return &result, nil
}
