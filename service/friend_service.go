package service

import (
	"database/sql"
	"dou_yin/dao/mysql/apply_dao"
	"dou_yin/dao/mysql/friend_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/VO/param"
	"dou_yin/pkg/snowflake"
	"dou_yin/pkg/utils"
	"dou_yin/service/DO"
	"errors"
)

func QueryFriendList(param param.QueryFriendListParam) (friendlist DO.FriendList, err error) {
	friends, err := friend_dao.QueryFriendshipList(param.UserID)
	if err != nil {
		return friendlist, err
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

	return friendlist, nil
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
	isFriend, err := IsFriend(param.UserID, param.FriendID)
	if err != nil {
		return application, err
	}
	if isFriend {
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

	// todo: tx
	err = apply_dao.Insert(application)
	if err != nil {
		return application, err
	}

	err = user_dao.AddPrivateChatWhite(param.UserID, param.FriendID)
	if err != nil {
		return application, err
	}

	return application, nil
}

func DeleteFriend(param param.DeleteFriendParam) (err error) {
	err = friend_dao.DeleteFriend(param.FriendshipID)
	if err != nil {
		return err
	}

	return nil
}

func SetPrivateChatBlack(param param.SetPrivateChatBlackParam) (err error) {
	err = user_dao.AddPrivateChatBlack(param.UserID, param.FriendID)
	if err != nil {
		return err
	}

	return nil
}

func UnBlockPrivateChat(param param.UnBlockPrivateChatParam) (err error) {
	err = user_dao.AddPrivateChatWhiteFromBlack(param.UserID, param.FriendID)
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
	//err = apply_dao.UpdateStatus(param.ApplyID, 1)
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

	return nil
}

func DisagreeFriendApply(param param.DisagreeFriendApplyParam) (err error) {
	//err = apply_dao.UpdateStatus(param.ApplyID, 2)
	err = apply_dao.Delete(param.ApplyID)
	if err != nil {
		return err
	}

	return nil
}
