package service

import (
	"dou_yin/dao/mysql/apply_dao"
	"dou_yin/dao/mysql/group_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/model/VO/response"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

func MGetGroupInfoByParam(info *param.QueryGroupInfoParam) (*response.GroupInfo, error) {
	group, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return nil, err
	}

	groupInfo, err := DO.MGetGroupInfofromPO(*group)
	if err != nil {
		return nil, err
	}

	return &response.GroupInfo{
		GroupID:     groupInfo.GroupID,
		OwnerID:     groupInfo.OwnerID,
		AdminIds:    *groupInfo.AdminIds,
		SilenceList: *groupInfo.SilenceList,
		UserIds:     *groupInfo.UserIds,
		GroupName:   groupInfo.GroupName,
		Description: *groupInfo.Description,
		CreateTime:  groupInfo.CreateTime,
		IsDeleted:   groupInfo.IsDeleted,
		Extra:       *groupInfo.Extra,
	}, nil
}

func MGetGroupListByParam(info *param.QueryGroupListParam) (*[]response.GroupJoin, error) {
	list, err := group_dao.MGetGroupListByUserID(info.UserID)
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

// 创建群聊
func CreateGroupInfoByParam(info *param.CreateGroupInfoParam) error {
	var groupInfo DO.GroupInfoDO
	groupInfo.GroupID = snowflake.GenID()
	groupInfo.OwnerID = info.OwnerID
	groupInfo.AdminIds = new([]int64)
	groupInfo.Description = info.Description
	groupInfo.CreateTime = utils.GetNowTime()
	groupInfo.Extra = nil
	groupInfo.GroupName = info.GroupName
	groupInfo.IsDeleted = false
	groupInfo.SilenceList = new([]int64)
	// 将好友拉入群聊
	groupInfo.UserIds = info.UserIDs

	groupInfoPO, err := DO.TurnGroupInfoPOfromDO(groupInfo)
	if err != nil {
		return err
	}

	var users []PO.UserPO
	var UserIDs []int64
	if info.UserIDs != nil {
		UserIDs = append(UserIDs, *info.UserIDs...)
	}
	UserIDs = append(UserIDs, info.OwnerID)
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
	}

	var groups []group_dao.GroupPO
	for _, id := range UserIDs {
		var groupDO DO.GroupDO
		groupDO.CreateTime = utils.GetNowTime()
		groupDO.Extra = new(DO.GroupExtra)
		groupDO.GroupID = groupInfo.GroupID
		groupDO.GroupName = groupInfoPO.GroupName
		groupDO.Type = 0
		groupDO.UserID = id

		groupPO, err := DO.TurnGroupPOfromDO(groupDO)
		if err != nil {
			return err
		}

		groups = append(groups, *groupPO)
	}

	// TODO: Tx
	// 创建群
	err = group_dao.CreateGroupInfo(groupInfoPO)
	if err != nil {
		return err
	}
	// 修改用户白名单
	for _, user := range users {
		err = user_dao.UpdateUserInfoByPO(&user)
		if err != nil {
			return err
		}
	}
	// 添加群关系信息
	for _, group := range groups {
		_, err = group_dao.CreateGroupByGroupPO(group)
		if err != nil {
			return err
		}
	}

	return nil
}

// 解散群聊
func DissolveGroupInfoByParam(info *param.DissolveGroupInfoParam) error {
	// TODO 加锁
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
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
			if ID != info.GroupID {
				white = append(white, ID)
			}
		}

		var gray []int64
		for _, ID := range *grayList {
			if ID != info.GroupID {
				gray = append(gray, ID)
			}
		}

		var black []int64
		for _, ID := range *blackList {
			if ID != info.GroupID {
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
	err = group_dao.UpdateGroupInfo(groupInfoPO)
	if err != nil {
		return err
	}

	// 修改群聊名单
	for _, userPO := range UserPOList {
		err := user_dao.UpdateUserInfoByPO(&userPO)
		if err != nil {
			return err
		}
	}

	// 删除加群关系
	for _, id := range userIds {
		_, err := group_dao.DeleteGroupByUserIDandGroupID(id, info.GroupID)
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyJoinGroupByParam(info *param.ApplyJoinGroupParam) error {
	application, err := apply_dao.MGetApplicationByGroupIDandUserID(info.GroupID, info.UserID)
	if err != nil {
		return err
	}
	if application != nil {
		fmt.Println("[ApplyJoinGroupByParam], apply had")
		return nil
	}

	_, err = group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}

	var apply PO.ApplyPO
	apply.Applicant = info.UserID
	apply.ApplyID = snowflake.GenID()
	apply.CreateTime = utils.GetNowTime()
	apply.Extra = nil
	apply.Reason = info.Reason
	apply.TargetID = info.GroupID
	apply.Status = 0
	apply.Type = 0
	// TODO:Tx
	err = apply_dao.CreateApplication(&apply)
	if err != nil {
		return err
	}

	return nil
}

// 退出群聊
func QuitGroupByParam(info *param.QuitGroupParam) error {
	ret, err := group_dao.IsGroupUser(info.UserID, info.GroupID)
	if err != nil {
		return err // 查询失败
	}
	if ret {
		userInfo, err := user_dao.QueryUserInfo(info.UserID)
		if err != nil {
			return err
		}
		whiteList, gratList, blackList, err := turnUserGroupList(userInfo)
		if err != nil {
			return err
		}

		var white, gray, black []int64
		for _, id := range *whiteList {
			if id != info.UserID {
				white = append(white, id)
			}
		}
		for _, id := range *gratList {
			if id != info.UserID {
				gray = append(gray, id)
			}
		}
		for _, id := range *blackList {
			if id != info.UserID {
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

		groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
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
				if id != info.UserID {
					userIds = append(userIds, id)
				}
			}
		}
		groupInfoDO.UserIds = &userIds
		var adminIds []int64
		if groupInfoDO.AdminIds != nil {
			for _, id := range *groupInfoDO.AdminIds {
				if id != info.UserID {
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
		// 在该群聊中，执行退群操作
		ret, err := group_dao.DeleteGroupByUserIDandGroupID(info.UserID, info.GroupID)
		if err != nil || !ret {
			return err // 删除失败
		}

		err = user_dao.UpdateUserInfoByPO(&userInfo)
		if err != nil {
			return err
		}

		err = group_dao.UpdateGroupInfo(GroupInfoPO)
		if err != nil {
			return err
		}
	}
	return nil // 退群成功
}

// 查看加群申请
func QueryGroupApplyListByParam(info *param.QueryGroupApplyListParam) (*[]response.GroupJoinApply, error) {
	list, err := apply_dao.MGetApplicationListByGroupID(info.GroupID)
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
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}
	userInfo, err := user_dao.QueryUserInfo(info.Applicant)
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
	userIds = append(*groupInfoDO.UserIds, info.Applicant)
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
	groupWhiteList = append(groupWhiteList, info.GroupID)

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
	groupDO.GroupID = info.GroupID
	groupDO.GroupName = groupInfoPO.GroupName
	groupDO.Type = 0
	groupDO.UserID = info.Applicant

	groupPO, err := DO.TurnGroupPOfromDO(groupDO)
	if err != nil {
		return err
	}

	// TODO:Tx
	err = group_dao.UpdateGroupInfo(groupInfo)
	if err != nil {
		return err
	}

	_, err = group_dao.CreateGroupByGroupPO(*groupPO)
	if err != nil {
		return err
	}

	err = user_dao.UpdateUserInfoByPO(&userInfo)
	if err != nil {
		return err
	}

	err = apply_dao.DeleteApplicationByApplyID(info.ApplyID)
	if err != nil {
		return err
	}

	return nil
}

// 不同意加入
func DisAgreeGroupApplyByParam(info *param.DisAgreeGroupApplyParam) error {

	err := apply_dao.DeleteApplicationByApplyID(info.ApplyID)
	if err != nil {
		return err
	}

	return nil
}

// 禁言
func SilenceByParam(info *param.SilenceParam) error {
	// TODO: 加锁
	group, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}

	groupInfo, err := DO.MGetGroupInfofromPO(*group)
	if err != nil {
		return err
	}

	silenceList := append(*groupInfo.SilenceList, info.TargetID)
	data, err := json.Marshal(silenceList)
	if err != nil {
		fmt.Println("[SilenceByParam], Marshal err is ", err.Error())
		return err
	}
	silence := string(data)
	group.SilenceList = &silence

	err = group_dao.UpdateGroupInfo(group)
	if err != nil {
		return err
	}
	return nil
}

// 解除禁言
func UnSilenceByParam(info *param.UnSilenceParam) error {
	// TODO: 加锁
	group, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}

	groupInfo, err := DO.MGetGroupInfofromPO(*group)
	if err != nil {
		return err
	}

	var silenceList []int64
	for _, id := range *groupInfo.SilenceList {
		if id != info.TargetID {
			silenceList = append(silenceList, id)
		}
	}
	data, err := json.Marshal(silenceList)
	if err != nil {
		fmt.Println("[UnSilenceByParam], Marshal err is ", err.Error())
		return err
	}
	silence := string(data)
	group.SilenceList = &silence

	err = group_dao.UpdateGroupInfo(group)
	if err != nil {
		return err
	}
	return nil
}

func TransferGroupByParam(info *param.TransferGroupParam) error {
	groupInfoRecord, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}
	OwnerRecord, err := group_dao.MGetGroupByUserIDandGroupID(info.UserID, info.GroupID)
	if err != nil {
		return err
	}
	TargetRecord, err := group_dao.MGetGroupByUserIDandGroupID(info.TargetID, info.GroupID)
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
		if TargetRecord.Type != 0 || (TargetRecord.Type == 0 && id != info.TargetID) {
			userList = append(userList, id)
		}
	}
	userList = append(userList, info.UserID) // 群主变为普通成员
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
			if id != info.TargetID {
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
	groupInfoRecord.OwnerID = info.TargetID

	OwnerRecord.Type = 0  // 成为普通成员
	TargetRecord.Type = 2 // 成为群主

	// TODO:Tx

	err = group_dao.UpdateGroupInfo(groupInfoRecord)
	if err != nil {
		return err
	}
	_, err = group_dao.UpdateGroupByGroupPO(*OwnerRecord)
	if err != nil {
		return err
	}
	_, err = group_dao.UpdateGroupByGroupPO(*TargetRecord)
	if err != nil {
		return err
	}

	return nil
}

func SetBlackListByParam(info *param.SetBlackListParam) error {
	userInfo, err := user_dao.QueryUserInfo(info.UserID)
	if err != nil {
		return err
	}

	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var whiteList, grayList []int64
	for _, id := range *groupWhiteList {
		if id != info.UserID {
			whiteList = append(whiteList, id)
		}
	}

	for _, id := range *groupGrayList {
		if id != info.UserID {
			grayList = append(grayList, id)
		}
	}

	blackList := append(*groupBlackList, info.UserID)

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	userInfo.GroupChatWhite = white
	userInfo.GroupChatGray = gray
	userInfo.GroupChatBlack = black

	err = user_dao.UpdateUserInfoByPO(&userInfo)
	if err != nil {
		return err
	}

	return nil
}

func SetGrayListByParam(info *param.SetGrayListParam) error {
	userInfo, err := user_dao.QueryUserInfo(info.UserID)
	if err != nil {
		return err
	}
	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var whiteList, blackList []int64
	for _, id := range *groupWhiteList {
		if id != info.UserID {
			whiteList = append(whiteList, id)
		}
	}

	for _, id := range *groupBlackList {
		if id != info.UserID {
			blackList = append(blackList, id)
		}
	}

	grayList := append(*groupGrayList, info.UserID)

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	userInfo.GroupChatWhite = white
	userInfo.GroupChatGray = gray
	userInfo.GroupChatBlack = black

	err = user_dao.UpdateUserInfoByPO(&userInfo)
	if err != nil {
		return err
	}

	return nil
}

func SetWhiteListByParam(info *param.SetWhiteListParam) error {
	userInfo, err := user_dao.QueryUserInfo(info.UserID)
	if err != nil {
		return err
	}

	groupWhiteList, groupGrayList, groupBlackList, err := turnUserGroupList(userInfo)
	if err != nil {
		return err
	}
	var grayList, blackList []int64
	for _, id := range *groupGrayList {
		if id != info.UserID {
			grayList = append(grayList, id)
		}
	}

	for _, id := range *groupBlackList {
		if id != info.UserID {
			blackList = append(blackList, id)
		}
	}

	whiteList := append(*groupWhiteList, info.UserID)

	white, gray, black, err := turnjsonList(whiteList, grayList, blackList)
	if err != nil {
		return err
	}
	userInfo.GroupChatWhite = white
	userInfo.GroupChatGray = gray
	userInfo.GroupChatBlack = black

	err = user_dao.UpdateUserInfoByPO(&userInfo)
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
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}
	group, err := group_dao.MGetGroupByUserIDandGroupID(info.TargetID, info.GroupID)
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	var userIDs []int64
	for _, id := range *groupInfoDO.UserIds {
		if id != info.TargetID {
			userIDs = append(userIDs, id)
		}
	}
	groupInfoDO.UserIds = &userIDs

	adminIDs := append(*groupInfoDO.AdminIds, info.TargetID)
	groupInfoDO.AdminIds = &adminIDs

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	group.Type = 1

	// TODO:Tx
	err = group_dao.UpdateGroupInfo(GroupInfoPO)
	if err != nil {
		return err
	}
	_, err = group_dao.UpdateGroupByGroupPO(*group)
	if err != nil {
		return err
	}
	return nil
}

func SetGroupUserByParam(info *param.SetGroupUserParam) error {
	groupInfoPO, err := group_dao.MGetGroupInfoByGroupID(info.GroupID)
	if err != nil {
		return err
	}
	group, err := group_dao.MGetGroupByUserIDandGroupID(info.TargetID, info.GroupID)
	if err != nil {
		return err
	}
	groupInfoDO, err := DO.MGetGroupInfofromPO(*groupInfoPO)
	if err != nil {
		return err
	}

	var AdminIds []int64
	for _, id := range *groupInfoDO.AdminIds {
		if id != info.TargetID {
			AdminIds = append(AdminIds, id)
		}
	}
	groupInfoDO.AdminIds = &AdminIds

	userIds := append(*groupInfoDO.UserIds, info.TargetID)
	groupInfoDO.UserIds = &userIds

	GroupInfoPO, err := DO.TurnGroupInfoPOfromDO(*groupInfoDO)
	if err != nil {
		return err
	}

	group.Type = 0

	// TODO:Tx
	err = group_dao.UpdateGroupInfo(GroupInfoPO)
	if err != nil {
		return err
	}
	_, err = group_dao.UpdateGroupByGroupPO(*group)
	if err != nil {
		return err
	}
	return nil
}

// 邀请加入群聊
func InviteJoinGroup(c *gin.Context) {

}

// 设置群备注
func SetGroupName(c *gin.Context) {

}
