package service

import (
	"dou_yin/dao/mysql/friend_dao"
	"dou_yin/dao/mysql/user_dao"
	"dou_yin/model/VO/param"
	"dou_yin/service/DO"
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
		return friendInfo, nil
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
