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
	friendInfo.Signature = *friend.Signature
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

	return friendInfo, nil
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

func HadApplied(userID int64, friendID int64) (bool, error) {
	_, err := apply_dao.QueryApplication(userID, friendID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddFriend(param param.AddFriendParam) (application DO.AddFriendApplication, err error) {
	// todo： 同时申请
	isFriend, err := IsFriend(param.UserID, param.FriendID)
	if err != nil {
		return application, err
	}
	if isFriend {
		logger.Log.Error("is friend")
		return application, errors.New("is friend")
	}

	hadApplied, err := HadApplied(param.UserID, param.FriendID)
	if err != nil {
		return application, err
	}
	if hadApplied {
		return application, errors.New("had applied")
	}

	application.ApplyID = snowflake.GenID()
	application.ApplicantID = param.UserID
	application.FriendID = param.FriendID
	application.Type = 1
	application.Reason = param.Reason
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

	var whileList PO.FriendCircleWhite
	if userPO.FriendCircleWhite != nil {
		err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &whileList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	whileList.WhiteList = append(whileList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whileList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	WhiteStr := string(whiteJson[:])

	err = user_dao.UpdateFriendCircleWhite(userID, WhiteStr)
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
	err = json.Unmarshal([]byte(*userPO.FriendCircleWhite), &whiteList)
	if err != nil {
		return err
	}
	whiteList.WhiteList = append(whiteList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whiteList)
	if err != nil {
		return err
	}
	whiteStr := string(whiteJson[:])

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

	var whileList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whileList)
		if err != nil {
			logger.Log.Error(err.Error())
			return nil
		}
	}
	whileList.WhiteList = append(whileList.WhiteList, friendID)
	whiteJson, err := json.Marshal(whileList)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}
	WhiteStr := string(whiteJson[:])

	err = user_dao.UpdatePrivateChatWhite(userID, WhiteStr)
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

	var whileList PO.PrivateChatWhite
	if userPO.PrivateChatWhite != nil {
		err = json.Unmarshal([]byte(*userPO.PrivateChatWhite), &whileList)
		if err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}
	whileList.WhiteList = append(whileList.WhiteList, friendID)
	whileJson, err := json.Marshal(whileList)
	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	whiteStr := string(whileJson[:])

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

		err = user_dao.UpdateFriendCircleBlack(userID, blackStr)
		if err != nil {
			return err
		}
	}

	return nil
}
