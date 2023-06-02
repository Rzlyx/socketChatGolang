package service

import (
	"database/sql"
	"dou_yin/dao/mysql/apply_dao"
	"dou_yin/dao/mysql/friend_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/logger"
	"dou_yin/model/PO"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"encoding/json"
	"errors"
)

func QueryFriendList(param param.QueryFriendListParam) (friendList DO.FriendList, err error) {
	friends, err := friend_dao.QueryFriendshipList(param.UserID)
	if err != nil {
		return friendList, err
	}

	for _, friend := range friends {
		var friendDO DO.Friend
		friendDO.FriendshipID = friend.FriendshipID
		if friend.FirstID != param.UserID {
			friendDO.FriendID = friend.FirstID
			friendDO.Name = friend.SecondRemarkFirst
		} else {
			friendDO.FriendID = friend.SecondID
			friendDO.Name = friend.FirstRemarkSecond
		}

		friendList.Friends = append(friendList.Friends, friendDO)
	}

	return friendList, nil
}

func QueryFriendInfo(param param.QueryFriendInfoParam) (friendInfo DO.FriendInfo, err error) {
	friend, err := user_dao.QueryUserInfo(param.FriendID)
	if err != nil {
		return friendInfo, err
	}

	friendInfo.UserID = friend.UserID
	friendInfo.UserName = friend.UserName
	friendInfo.Sex = friend.Sex
	friendInfo.PhoneNumber = friend.PhoneNumber
	friendInfo.Email = friend.Email
	if friend.Signature != nil {
		friendInfo.Signature = *friend.Signature
	}
	friendInfo.Birthday = friend.Birthday
	friendInfo.Status = friend.Status

	friendship, err := friend_dao.QueryFriendship(param.FriendshipID)
	if err != nil {
		return friendInfo, err
	}

	if friendship.FirstID == param.FriendID {
		if friendship.IsSecondRemarkFirst {
			friendInfo.IsRemark = true
			friendInfo.Remark = friendship.SecondRemarkFirst
		} else {
			friendInfo.IsRemark = false
		}
	} else {
		if friendship.IsFirstRemarkSecond {
			friendInfo.IsRemark = true
			friendInfo.Remark = friendship.FirstRemarkSecond
		} else {
			friendInfo.IsRemark = false
		}
	}

	userInfo, err := user_dao.QueryUserInfo(param.UserID)
	var privateChatBlackList PO.PrivateChatBlack
	if userInfo.PrivateChatBlack != nil {
		err = json.Unmarshal([]byte(*userInfo.PrivateChatBlack), &privateChatBlackList)
		if err != nil {
			logger.Log.Error(err.Error())
			return friendInfo, err
		}
		friendInfo.IsPrivateChatBlack = IsContains(privateChatBlackList.BlackList, param.FriendID)
	} else {
		friendInfo.IsPrivateChatBlack = false
	}

	var friendCircleBlackList PO.FriendCircleBlack
	if userInfo.FriendCircleBlack != nil {
		err = json.Unmarshal([]byte(*userInfo.FriendCircleBlack), &friendCircleBlackList)
		if err != nil {
			logger.Log.Error(err.Error())
			return friendInfo, err
		}
		friendInfo.IsFriendCircleBlack = IsContains(friendCircleBlackList.BlackList, param.FriendID)
	} else {
		friendInfo.IsFriendCircleBlack = false
	}

	return friendInfo, nil
}

func IsContains(list []int64, target int64) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}

	return false
}

func IsFriend(firstID int64, secondID int64) (bool, error) {
	_, err := friend_dao.QueryFriendshipBy2ID(firstID, secondID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func HadApplied(userID int64, friendID int64) (bool, PO.ApplyPO, error) {
	applyPO, err := apply_dao.QueryApplication(userID, friendID)
	if err == sql.ErrNoRows {
		return false, applyPO, nil
	}
	if err != nil {
		return false, applyPO, err
	}

	return true, applyPO, nil
}

func AddFriend(addFriendParam param.AddFriendParam) (application DO.AddFriendApplication, err error) {
	hadApplied, applyPO, err := HadApplied(addFriendParam.FriendID, addFriendParam.UserID)
	if err != nil {
		return application, err
	}
	if hadApplied {
		err := AgreeFriendApply(param.AgreeFriendApplyParam{
			FriendID: addFriendParam.FriendID,
			UserID:   addFriendParam.UserID,
			ApplyID:  applyPO.ApplyID,
		})
		if err != nil {
			return application, err
		}

		application.IsBeFriend = true
		return application, errors.New("be FRIEND by friend's application")
	}
	application.IsBeFriend = false

	isFriend, err := IsFriend(addFriendParam.UserID, addFriendParam.FriendID)
	if err != nil {
		return application, err
	}
	if isFriend {
		logger.Log.Error("is friend")
		return application, errors.New("is friend")
	}

	hadApplied, _, err = HadApplied(addFriendParam.UserID, addFriendParam.FriendID)
	if err != nil {
		return application, err
	}
	if hadApplied {
		return application, errors.New("had applied")
	}

	application.ApplyID = snowflake.GenID()
	application.ApplicantID = addFriendParam.UserID
	application.FriendID = addFriendParam.FriendID
	application.Type = 1
	application.Reason = addFriendParam.Reason
	application.CreateTime = utils.GetNowTime()

	err = apply_dao.Insert(application)
	if err != nil {
		return application, err
	}

	return application, nil
}

func DeleteFriend(param param.DeleteFriendParam) (err error) {
	// todo: tx
	err = friend_dao.DeleteFriend(param.FriendshipID)
	if err != nil {
		return err
	}

	err = RemoveFriendFromWhiteBlackList(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	err = RemoveFriendFromWhiteBlackList(param.FriendID, param.UserID)
	if err != nil {
		return err
	}

	return nil
}

func SetPrivateChatBlack(param param.SetPrivateChatBlackParam) (err error) {
	err = AddPrivateChatBlack(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func UnBlockPrivateChat(param param.UnBlockPrivateChatParam) (err error) {
	err = AddPrivateChatWhiteFromBlack(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func SetFriendCircleBlack(param param.SetFriendCircleBlackParam) (err error) {
	err = AddFriendCircleBlack(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func UnBlockFriendCircle(param param.UnBlockFriendCircleParam) (err error) {
	err = AddFriendCircleWhiteFromBlack(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func QueryFriendApply(param param.QueryFriendApplyParam) (applications DO.FriendApplicationList, err error) {
	applyPOs, err := apply_dao.QueryFriendApply(param.UserID)
	if err != nil {
		return applications, err
	}

	for _, apply := range applyPOs {
		var application DO.FriendApplication
		application.ApplicantID = apply.Applicant
		application.ApplyID = apply.ApplyID
		application.Reason = apply.Reason
		application.Status = apply.Status
		application.UserID = apply.TargetID

		applicantInfo, err := user_dao.QueryUserInfo(application.ApplicantID)
		if err != nil {
			return applications, err
		}
		application.ApplicantName = applicantInfo.UserName

		applications.Applications = append(applications.Applications, application)
	}

	return applications, nil
}

func AgreeFriendApply(param param.AgreeFriendApplyParam) (err error) {
	hadApplied, _, err := HadApplied(param.FriendID, param.UserID)
	if err != nil {
		return err
	}
	if !hadApplied {
		return errors.New("there is no application")
	}
	// todo: tx
	err = apply_dao.Delete(param.ApplyID)
	if err != nil {
		return err
	}

	FirstInfo, err := user_dao.QueryUserInfo(param.FriendID)
	if err != nil {
		return err
	}
	SecondIndo, err := user_dao.QueryUserInfo(param.UserID)
	if err != nil {
		return err
	}

	var friendship DO.Friendship
	friendship.FriendshipID = snowflake.GenID()
	friendship.FirstID = param.FriendID
	friendship.SecondID = param.UserID
	friendship.FirstRemarkSecond = SecondIndo.UserName
	friendship.SecondRemarkFirst = FirstInfo.UserName

	err = friend_dao.Insert(friendship)
	if err != nil {
		return err
	}

	err = AddPrivateChatWhite(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	err = AddFriendCircleWhite(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	err = AddPrivateChatWhite(param.FriendID, param.UserID)
	if err != nil {
		return err
	}

	err = AddFriendCircleWhite(param.FriendID, param.UserID)
	if err != nil {
		return err
	}

	return nil
}

func DisagreeFriendApply(param param.DisagreeFriendApplyParam) (err error) {
	hadApplied, _, err := HadApplied(param.FriendID, param.UserID)
	if err != nil {
		return err
	}
	if !hadApplied {
		return errors.New("there no apply")
	}
	err = apply_dao.Delete(param.ApplyID)
	if err != nil {
		return err
	}

	return nil
}

func AddFriendCircleWhite(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var whiteList PO.FriendCircleWhite
	if userPO.FriendCircleWhite != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &whiteList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	whiteList.WhiteList = append(whiteList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	err = user_dao.UpdateFriendCircleWhite(userID, whiteStr)
	if err != nil {
		return err
	}

	return nil
}

func AddFriendCircleWhiteFromBlack(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var whiteList PO.FriendCircleWhite
	if userPO.FriendCircleWhite != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &whiteList)
		if err != nil {
			return err
		}
	}
	whiteList.WhiteList = append(whiteList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		return err
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	var blackList PO.FriendCircleBlack
	if userPO.FriendCircleBlack != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleBlack), &blackList)
		if err != nil {
			return err
		}
	}
	for index, id := range blackList.BlackList {
		if id == friendID {
			blackList.BlackList = append(blackList.BlackList[:index], blackList.BlackList[index+1:]...)
			break
		}
	}
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		return err
	}
	blackStr := string(blackJson[:])
	if len(blackList.BlackList) == 0 {
		blackStr = ""
	}

	err = user_dao.UpdateFriendCircleBlackWhite(userID, whiteStr, blackStr)
	if err != nil {
		return err
	}

	return nil
}

func AddFriendCircleBlack(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var blackList PO.FriendCircleBlack
	if userPO.FriendCircleBlack != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleBlack), &blackList)
		if err != nil {
			return err
		}
	}
	blackList.BlackList = append(blackList.BlackList, friendID)
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		return err
	}
	blackStr := string(blackJson[:])
	if len(blackList.BlackList) == 0 {
		blackStr = ""
	}

	var whiteList PO.FriendCircleWhite
	err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &whiteList)
	if err != nil {
		return err
	}
	for index, id := range whiteList.WhiteList {
		if id == friendID {
			whiteList.WhiteList = append(whiteList.WhiteList[:index], whiteList.WhiteList[index+1:]...)
			break
		}
	}
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		return err
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	err = user_dao.UpdateFriendCircleBlackWhite(userID, whiteStr, blackStr)
	if err != nil {
		return err
	}

	return nil
}

func AddPrivateChatWhite(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var whiteList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whiteList)
		if err != nil {
			logger.Log.Error(err.Error())
			return nil
		}
	}
	whiteList.WhiteList = append(whiteList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	err = user_dao.UpdatePrivateChatWhite(userID, whiteStr)
	if err != nil {
		return err
	}

	return nil
}

func AddPrivateChatWhiteFromBlack(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var whiteList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whiteList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	whiteList.WhiteList = append(whiteList.WhiteList, friendID)
	whileJson, err := json.Marshal(whiteList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	whiteStr := string(whileJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	var blackList PO.PrivateChatBlack
	err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &blackList)
	if err != nil {
		return err
	}
	for index, id := range blackList.BlackList {
		if id == friendID {
			blackList.BlackList = append(blackList.BlackList[:index], blackList.BlackList[index+1:]...)
			break
		}
	}
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		return err
	}
	blackStr := string(blackJson[:])
	if len(blackList.BlackList) == 0 {
		blackStr = ""
	}

	err = user_dao.UpdatePrivateChatBlackWhite(userID, whiteStr, blackStr)
	if err != nil {
		return err
	}

	return nil
}

func AddPrivateChatBlack(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	var blackList PO.PrivateChatBlack
	if userPO.PrivateChatBlack != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &blackList)
		if err != nil {
			return err
		}
	}
	blackList.BlackList = append(blackList.BlackList, friendID)
	blackJson, err := json.Marshal(blackList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	blackStr := string(blackJson[:])
	if len(blackList.BlackList) == 0 {
		blackStr = ""
	}

	var whiteList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whiteList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	for index, id := range whiteList.WhiteList {
		if id == friendID {
			whiteList.WhiteList = append(whiteList.WhiteList[:index], whiteList.WhiteList[index+1:]...)
			break
		}
	}
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	err = user_dao.UpdatePrivateChatBlackWhite(userID, whiteStr, blackStr)
	if err != nil {
		return err
	}

	return nil
}

func RemoveFriendFromWhiteBlackList(userID int64, friendID int64) (err error) {
	userPO, err := user_dao.QueryUserInfo(userID)
	if err != nil {
		return err
	}

	privateChatInWhite := false
	var whiteList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whiteList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	for index, id := range whiteList.WhiteList {
		if id == friendID {
			privateChatInWhite = true
			whiteList.WhiteList = append(whiteList.WhiteList[:index], whiteList.WhiteList[index+1:]...)
			break
		}
	}
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	whiteStr := string(whiteJson[:])
	if len(whiteList.WhiteList) == 0 {
		whiteStr = ""
	}

	err = user_dao.UpdatePrivateChatWhite(userID, whiteStr)
	if err != nil {
		return err
	}

	if !privateChatInWhite {
		var blackList PO.PrivateChatBlack
		if userPO.PrivateChatBlack != nil {
			err = json.Unmarshal([]byte(*userPO.PrivateChatBlack), &blackList)
			if err != nil {
				logger.Log.Error(err.Error())
				return err
			}
		}
		for index, id := range blackList.BlackList {
			if id == friendID {
				blackList.BlackList = append(blackList.BlackList[:index], blackList.BlackList[index+1:]...)
				break
			}
		}
		blackJson, err := json.Marshal(blackList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
		blackStr := string(blackJson[:])
		if len(blackList.BlackList) == 0 {
			blackStr = ""
		}

		err = user_dao.UpdatePrivateChatBlack(userID, blackStr)
		if err != nil {
			return err
		}
	}

	friendCircleInWhite := false
	var friendCircleWhite PO.FriendCircleWhite
	if userPO.FriendCircleWhite != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &friendCircleWhite)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &friendCircleWhite)
	if err != nil {
		return err
	}
	for index, id := range friendCircleWhite.WhiteList {
		if id == friendID {
			friendCircleWhite.WhiteList = append(friendCircleWhite.WhiteList[:index], friendCircleWhite.WhiteList[index+1:]...)
			friendCircleInWhite = true
			break
		}
	}
	friendCircleWhiteJson, err := json.Marshal(friendCircleWhite)
	if err != nil {
		return err
	}
	friendCircleWhiteStr := string(friendCircleWhiteJson[:])
	if len(friendCircleWhite.WhiteList) == 0 {
		friendCircleWhiteStr = ""
	}

	err = user_dao.UpdateFriendCircleWhite(userID, friendCircleWhiteStr)
	if err != nil {
		return err
	}

	if !friendCircleInWhite {
		var blackList PO.FriendCircleBlack
		err = json.Unmarshal([]byte(*userPO.FriendCircleBlack), &blackList)
		if err != nil {
			return err
		}
		for index, id := range blackList.BlackList {
			if id == friendID {
				blackList.BlackList = append(blackList.BlackList[:index], blackList.BlackList[index+1:]...)
				break
			}
		}
		blackJson, err := json.Marshal(blackList)
		if err != nil {
			return err
		}
		blackStr := string(blackJson[:])
		if len(blackList.BlackList) == 0 {
			blackStr = ""
		}

		err = user_dao.UpdateFriendCircleBlack(userID, blackStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetFriendRemark(param param.SetFriendRemark) (err error) {
	friendship, err := friend_dao.QueryFriendship(param.FriendshipID)
	if err != nil {
		return err
	}

	var realName string
	if param.Remark != nil {
		userInfo, err := user_dao.QueryUserInfo(param.FriendID)
		if err != nil {
			return err
		}

		realName = userInfo.UserName
	}

	if friendship.FirstID == param.UserID {
		err = friend_dao.UpdateFirstRemarkSecond(param.FriendshipID, *param.Remark, realName)
		if err != nil {
			return err
		}
	} else {
		err = friend_dao.UpdateSecondRemarkFirst(param.FriendshipID, *param.Remark, realName)
		if err != nil {
			return nil
		}
	}

	return nil
}
